package main

import (
	"fmt"
	"time"

	"github.com/AT-SmFoYcSNaQ/AT2023/Go/order/messages"
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type ReceivedOrder struct {
	UserID   string
	ItemID   string
	Quantity int32
}

type OrderActor struct {
	remoting *remote.Remote
	context  *actor.RootContext
}

func (actor *OrderActor) Receive(context actor.Context) {
	// Handle incoming messages
	switch msg := context.Message().(type) {
	case *messages.ReceiveOrder_Request:
		// Received order from customer
		order := ReceivedOrder{
			UserID:   msg.UserId,
			ItemID:   msg.ItemId,
			Quantity: msg.Quantity,
		}
		actor.handleOrderReceived(order, context.Self()) // Pass the order and self reference
	case *messages.CheckAvailability_Response:
		// Availability response from inventory actor
		actor.handleAvailabilityChecked(msg.IsAvailable, context.Self()) // Pass availability status and self reference
	case *messages.PaymentInfo_Response:
		// Payment response from payment actor
		actor.handlePaymentInfoReceived(msg.Successful, context.Self()) // Pass payment status and self reference
	}
}

func (actor *OrderActor) handleOrderReceived(order ReceivedOrder, self *actor.PID) {
	fmt.Println("Received message from customer!")

	// Create a message to check availability
	message := &messages.CheckAvailability_Request{
		Sender:   self,
		ItemId:   order.ItemID,
		Quantity: order.Quantity,
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

func (actor *OrderActor) handleAvailabilityChecked(isAvailable bool, self *actor.PID) {
	fmt.Println("Received message from inventory actor!")

	// Spawn the notification actor
	spawnResponse, err := actor.remoting.SpawnNamed("127.0.0.1:8092", "notification-actor", "notification-actor", time.Second)
	if err != nil {
		panic(err)
	}

	if isAvailable {
		// Item is available
		actor.context.Send(spawnResponse.Pid, &messages.OrderUpdated_Request{Sender: self, Status: "Pending"})
		actor.prepareOrder(10 * time.Second)
		actor.context.Send(spawnResponse.Pid, &messages.OrderUpdated_Request{Sender: self, Status: "Prepared"})
		actor.processPayment(self) // Pass self reference for payment actor
	} else {
		// Item is out of stock
		message := &messages.OrderUpdated_Request{
			Sender: self,
			Status: "OutOfStock",
		}
		actor.context.Send(spawnResponse.Pid, message)
	}
}

func (actor *OrderActor) handlePaymentInfoReceived(successful bool, self *actor.PID) {
	fmt.Println("Received message from payment actor!")

	// Spawn the notification actor
	spawnResponse, err := actor.remoting.SpawnNamed("127.0.0.1:8092", "notification-actor", "notification-actor", time.Second)
	if err != nil {
		panic(err)
	}

	status := "PaymentFailed"
	if successful {
		status = "Payment"
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

func (actor *OrderActor) processPayment(self *actor.PID) {
	// Spawn the payment actor
	spawnResponse, err := actor.remoting.SpawnNamed("127.0.0.1:8093", "payment-actor", "payment-actor", time.Second)
	if err != nil {
		panic(err)
	}

	message := &messages.EmptyMessage{Sender: self, Message: ""}
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

	// Configure and start remote communication
	remoteConfig := remote.Configure("127.0.0.1", 8090)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()

	// Get the root context of the actor system
	context := system.Root

	// Create the order actor and register it with the remote system
	orderActorProps := actor.PropsFromProducer(func() actor.Actor { return &OrderActor{remoting: remoting, context: context} })
	remoting.Register("order-actor", orderActorProps)

	console.ReadLine()
}
