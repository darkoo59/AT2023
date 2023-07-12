package main

import (
	"fmt"
	notifMessage "github.com/AT-SmFoYcSNaQ/AT2023/Go/notification/messages"
	orderMessage "github.com/AT-SmFoYcSNaQ/AT2023/Go/order/messages"
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"time"
)

type PaymentReq struct {
	Quantity       int32
	PricePerItem   float32
	OrderId        string
	AccountNumber  string
	AccountBalance float32
	UserId         string
}

type PaymentActor struct {
	remoting *remote.Remote
	context  *actor.RootContext
}

func (actor *PaymentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *orderMessage.PaymentReq:
		paymentReq := PaymentReq{
			Quantity:       msg.Quantity,
			PricePerItem:   msg.PricePerItem,
			OrderId:        msg.OrderId,
			AccountBalance: msg.AccountBalance,
			AccountNumber:  msg.AccountNumber,
			UserId:         msg.UserId,
		}
		actor.handlePaymentRequest(paymentReq, context.Self())
	}
}

func (actor *PaymentActor) handlePaymentRequest(paymentReq PaymentReq, self *actor.PID) {
	fmt.Println("Received payment request")
	totalPrice := paymentReq.PricePerItem * float32(paymentReq.Quantity)
	paymentSuccessful := totalPrice <= paymentReq.AccountBalance

	if !paymentSuccessful {
		actor.sendPaymentInfo(paymentReq, paymentSuccessful)
		actor.sendPaymentInfoNotification(paymentReq, paymentSuccessful, self)
		return
	}
	paymentReq.AccountBalance -= totalPrice

	actor.sendPaymentInfo(paymentReq, paymentSuccessful)
	actor.sendPaymentInfoNotification(paymentReq, paymentSuccessful, self)
}

func (actor *PaymentActor) sendPaymentInfo(paymentReq PaymentReq, isSuccessful bool) {
	fmt.Println("Sending payment info to order actor")
	message := &orderMessage.OrderPaymentInfo{
		OrderId:        paymentReq.OrderId,
		IsSuccessful:   isSuccessful,
		AccountNumber:  paymentReq.AccountNumber,
		AccountBalance: paymentReq.AccountBalance,
	}

	spawnResponse, err := actor.remoting.SpawnNamed("127.0.0.1:8090", "order-actor", "order-actor", time.Second)

	if err != nil {
		panic(err)
	}

	actor.context.Send(spawnResponse.Pid, message)
	fmt.Println("Sent payment message to order actor!")
}

func (actor *PaymentActor) sendPaymentInfoNotification(paymentReq PaymentReq, isSuccessful bool, self *actor.PID) {
	fmt.Println("Sending payment info to notification actor")

	paymentMessage := ""

	if isSuccessful {
		paymentMessage = fmt.Sprintf("Payment for orderId %s was successful! New account balance is %.2f", paymentReq.OrderId, paymentReq.AccountBalance)
	} else {
		paymentMessage = fmt.Sprintf("Payment for orderId %s was not successful,account balance did not change.", paymentReq.OrderId)
	}

	spawnResponse, err := actor.remoting.SpawnNamed("127.0.0.1:8098", "notification-actor", "notification-actor", time.Second)

	messageContent := &notifMessage.Message{
		Content: paymentMessage,
		Action:  "",
		OrderId: "",
	}

	message := &notifMessage.Notification{
		Message:    messageContent,
		ReceiverId: paymentReq.UserId,
	}

	if err != nil {
		panic(err)
	}

	actor.context.Send(spawnResponse.Pid, message)
	fmt.Println("Sent payment message to notification actor!")
}

func main() {

	system := actor.NewActorSystem()

	// Configure and start remote communication
	remoteConfig := remote.Configure("192.168.1.34", 8093)
	remoting := remote.NewRemote(system, remoteConfig)

	remoting.Start()

	// Get the root context of the actor system
	context := system.Root

	// Create the payment actor and register it with the remote system
	paymentActorProps := actor.PropsFromProducer(func() actor.Actor { return &PaymentActor{remoting: remoting, context: context} })

	remoting.Register("payment-actor", paymentActorProps)

	spawnResponse, err := remoting.SpawnNamed("192.168.1.48:8092", "notification-actor", "notification-actor", time.Second)

	if err != nil {
		panic(err)
		return
	}

	paymentMessage := "NOTIFICATIONNOTIFICATION"

	messageContent := &notifMessage.Message{
		Content: paymentMessage,
		Action:  "",
		OrderId: "",
	}

	message := &notifMessage.Notification{
		Message:    messageContent,
		ReceiverId: "5a543ba3-9ee2-48f9-b3db-d85c443a1512",
	}

	spawnedPID := spawnResponse.Pid
	context.Send(spawnedPID, message)
	console.ReadLine()
}
