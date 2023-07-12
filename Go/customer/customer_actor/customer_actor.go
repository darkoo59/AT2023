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
	PID         *actor.PID
	RootContext *actor.RootContext
	Remoting    *remote.Remote
	Logger      *zap.Logger
}

func CreateCustomerActor(logger *zap.Logger) *CustomerActor {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatal(err.Error())
		return nil
	}

	system := actor.NewActorSystem()
	remoteConfig := remote.Configure(loadConfig.ActorCustomerAddress, loadConfig.ActorCustomerPort)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()

	actorContext := system.Root
	customerActor := &CustomerActor{Logger: logger, RootContext: actorContext, Remoting: remoting}
	customerActorProps := actor.PropsFromProducer(func() actor.Actor {
		return customerActor
	})
	pid := actorContext.Spawn(customerActorProps)
	customerActor.PID = pid
	remoting.Register("customer-actor", customerActorProps)
	logger.Info("Customer actor registered")

	return customerActor
}

type Order struct {
	UserId   string
	ItemId   string
	Quantity int
}

func (customerActor *CustomerActor) Send(pid *actor.PID, message interface{}) {
	customerActor.RootContext.Send(pid, message)
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
		loadConfig.ActorOrderAddress+":"+fmt.Sprint(loadConfig.ActorOrderPort),
		"order-actor",
		"order-actor",
		5*time.Second)
	if err != nil {
		customerActor.Logger.Error(err.Error())
		panic(err)
	}

	customerActor.Send(spawnResponse.Pid, order)
	customerActor.Logger.Info("Message sent to order-actor")
}
