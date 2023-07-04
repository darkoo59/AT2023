package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type PaymentActor struct {
	remoting *remote.Remote
	context  *actor.RootContext
}

func main() {

}
