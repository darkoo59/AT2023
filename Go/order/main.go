package main

import (
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/order/model"
	"time"

	"github.com/AT-SmFoYcSNaQ/AT2023/Go/order/messages"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/order/service"
	paymentMessages "github.com/AT-SmFoYcSNaQ/AT2023/Go/payment/messages/Go/messages"
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type OrderActor struct {
	remoting *remote.Remote
	context  *actor.RootContext
	service  *service.OrderService
}

func (actor *OrderActor) Receive(context actor.Context) {
	// Handle incoming messages
	switch msg := context.Message().(type) {
	case *messages.ReceiveOrder_Request:
		// Received order from customer
		order := model.Order{
			UserId:         msg.UserId,
			ItemId:         msg.ItemId,
			Quantity:       msg.Quantity,
			AccountBalance: msg.AccountBalance,
			PricePerItem:   msg.PricePerItem,
			OrderStatus:    "Pending",
		}
		actor.handleOrderReceived(&order, context.Self()) // Pass the order and self reference
	case messages.CheckAvailability_Response:
		// Availability response from inventory actor
		actor.handleAvailabilityChecked(&msg, context.Self()) // Pass availability status and self reference
	case paymentMessages.OrderPaymentInfo:
		// Payment response from payment actor
		actor.handlePaymentInfoReceived(&msg, context.Self()) // Pass payment status and self reference
	}
}

func (actor *OrderActor) handleOrderReceived(order *model.Order, self *actor.PID) {
	fmt.Println("Received message from customer!")

	orderCreated, err := actor.service.Insert(order)
	if err != nil {
		return
	}

	// Create a message to check availability
	message := &messages.CheckAvailability_Request{
		Sender:   self,
		ItemId:   order.ItemId,
		Quantity: order.Quantity,
		OrderId:  orderCreated,
	}

	// Spawn the inventory actor
	spawnResponse, err := actor.remoting.SpawnNamed("127.0.0.1:8098", "inventory-actor", "inventory-actor", time.Second)
	if err != nil {
		panic(err)
	}

	// Send the availability request to the inventory actor
	actor.context.Send(spawnResponse.Pid, message)
	fmt.Println("Sent message to the inventory actor!")
}

func (actor *OrderActor) handleAvailabilityChecked(request *messages.CheckAvailability_Response, self *actor.PID) {
	fmt.Println("Received message from inventory actor!")

	// Spawn the notification actor
	spawnResponse, err := actor.remoting.SpawnNamed("127.0.0.1:8092", "notification-actor", "notification-actor", time.Second)
	if err != nil {
		panic(err)
	}

	if request.IsAvailable {
		// Item is available

		orderUpdated, err := actor.service.GetById(request.OrderId)
		if err != nil {
			panic(err)
		}
		orderUpdated.OrderStatus = "Pending"
		_, err = actor.service.Insert(orderUpdated)
		if err != nil {
			return
		}
		actor.context.Send(spawnResponse.Pid, &messages.OrderUpdated_Request{Sender: self, Status: "Pending"})
		actor.prepareOrder(10 * time.Second)
		orderUpdated, err = actor.service.GetById(request.OrderId)
		if err != nil {
			panic(err)
		}
		orderUpdated.OrderStatus = "Prepared"
		_, err = actor.service.Insert(orderUpdated)
		if err != nil {
			return
		}
		actor.context.Send(spawnResponse.Pid, &messages.OrderUpdated_Request{Sender: self, Status: "Prepared"})
		actor.processPayment(self, request) // Pass self reference for payment actor
	} else {
		// Item is out of stock
		message := &messages.OrderUpdated_Request{
			Sender: self,
			Status: "OutOfStock",
		}

		orderUpdated, err := actor.service.GetById(request.OrderId)
		if err != nil {
			panic(err)
		}
		orderUpdated.OrderStatus = "OutOfStock"
		_, err = actor.service.Insert(orderUpdated)
		if err != nil {
			return
		}

		actor.context.Send(spawnResponse.Pid, message)
	}
}

func (actor *OrderActor) handlePaymentInfoReceived(request *paymentMessages.OrderPaymentInfo, self *actor.PID) {
	fmt.Println("Received message from payment actor!")

	// Spawn the notification actor
	spawnResponse, err := actor.remoting.SpawnNamed("127.0.0.1:8092", "notification-actor", "notification-actor", time.Second)
	if err != nil {
		panic(err)
	}

	status := "PaymentFailed"
	if request.IsSuccessful {
		status = "Payment"
	}

	orderUpdated, err := actor.service.GetById(request.OrderId)
	if err != nil {
		panic(err)
	}
	orderUpdated.OrderStatus = status
	_, err = actor.service.Insert(orderUpdated)
	if err != nil {
		return
	}

	message := &messages.OrderUpdated_Request{
		Sender: self,
		Status: status,
	}

	actor.context.Send(spawnResponse.Pid, message)
}

func (actor *OrderActor) prepareOrder(seconds time.Duration) {
	fmt.Println("Order preparing process in progress!")
	time.Sleep(seconds)
	fmt.Println("Order preparing process done!")
}

func (actor *OrderActor) processPayment(self *actor.PID, request *messages.CheckAvailability_Response) {
	// Spawn the payment actor
	spawnResponse, err := actor.remoting.SpawnNamed("127.0.0.1:8093", "payment-actor", "payment-actor", time.Second)
	if err != nil {
		panic(err)
	}

	order, err := actor.service.GetById(request.OrderId)
	if err != nil {
		panic(err)
	}

	message := &messages.PaymentReq{
		Quantity:       order.Quantity,
		PricePerItem:   request.ItemPrice,
		OrderId:        request.OrderId,
		AccountBalance: float32(order.AccountBalance),
		UserId:         order.UserId,
	}
	actor.context.Send(spawnResponse.Pid, message)
}

/*
Configuration for order-actor:
  - kind: order-actor
  - address: 127.0.0.1:8090

In order to works, required configuration for other actors are:
  - kind: notification-actor, address: 127.0.0.1:8092
  - kind: inventory-actor, address: 127.0.0.1:8098
  - kind: payment-actor, address: 127.0.0.1:8093
*/
func main() {

	system := actor.NewActorSystem()
	orderService := service.CreateOrderService()

	// Configure and start remote communication
	remoteConfig := remote.Configure("127.0.0.1", 8090)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()

	// Get the root context of the actor system
	context := system.Root

	// Create the order actor and register it with the remote system
	orderActorProps := actor.PropsFromProducer(func() actor.Actor { return &OrderActor{remoting: remoting, context: context, service: orderService} })
	remoting.Register("order-actor", orderActorProps)

	console.ReadLine()
}
