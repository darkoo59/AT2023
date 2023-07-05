package main

import (
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type PaymentActor struct {
	remoting *remote.Remote
	context  *actor.RootContext
}

func (actor *PaymentActor) Receive(context actor.Context) {

}

//func (actor *PaymentActor) notifyAboutPayment()

func main() {

	system := actor.NewActorSystem()

	// Configure and start remote communication
	remoteConfig := remote.Configure("127.0.0.1", 8093)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()

	// Get the root context of the actor system
	context := system.Root

	// Create the order actor and register it with the remote system
	orderActorProps := actor.PropsFromProducer(func() actor.Actor { return &PaymentActor{remoting: remoting, context: context} })
	remoting.Register("payment-actor", orderActorProps)

	console.ReadLine()
}
