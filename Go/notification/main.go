package main

import (
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/notification/messages"
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type NotificationActor struct{}

func (*NotificationActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Notification:
		println(context.Message().(*messages.Notification).Message)
		context.Send(msg.Sender, &messages.Response{
			SomeValue: "response",
		})
	}
}

func main() {
	system := actor.NewActorSystem()
	remoteConfig := remote.Configure("127.0.0.1", 8092)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()

	remoting.Register("notification-actor", actor.PropsFromProducer(func() actor.Actor { return &NotificationActor{} }))
	console.ReadLine()
}
