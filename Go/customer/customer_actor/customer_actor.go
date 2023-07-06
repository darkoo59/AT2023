package customer_actor

import (
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/config"
	messages "github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/message"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"go.uber.org/zap"
	"time"
)

type CustomerActor struct {
	Remoting *remote.Remote
	Context  *actor.RootContext
	Logger   *zap.Logger
}

type Order struct {
	UserId   string
	ItemId   string
	Quantity int
}

func (customerActor *CustomerActor) Spawn() *actor.PID {
	customerActorProps := actor.PropsFromProducer(func() actor.Actor {
		return customerActor
	})
	return customerActor.Context.Spawn(customerActorProps)
}

func (customerActor *CustomerActor) Send(pid *actor.PID, message interface{}) {
	customerActor.Context.Send(pid, message)
}

func (customerActor *CustomerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.ReceiveOrder_Request:
		customerActor.sendOrderRequest(msg)
	}
}

func (customerActor *CustomerActor) sendOrderRequest(order *messages.ReceiveOrder_Request) {
	customerActor.Logger.Info("Message received from http server")

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		customerActor.Logger.Error(err.Error())
	}

	spawnResponse, err := customerActor.Remoting.SpawnNamed(
		loadConfig.ActorHostAddress+":"+fmt.Sprint(loadConfig.ActorOrderPort),
		"order-actor",
		"order-actor",
		time.Second)
	if err != nil {
		customerActor.Logger.Error(err.Error())
		panic(err)
	}
	customerActor.Send(spawnResponse.Pid, order)
	customerActor.Logger.Info("Message sent to order-actor")
}
