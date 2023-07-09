package main

import (
	"encoding/json"
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/notification/messages"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/notification/socket"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type NotificationActor struct {
	hub *socket.Hub
}

type NotificationResponse struct {
	Message string
	Type    int32
}

func (h *NotificationActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Notification:
		fmt.Println("Received notification: " + context.Message().(*messages.Notification).Message)
		client := (*h.hub).GetClient(msg.ReceiverId)
		if client == nil {
			return
		}
		data, _ := json.Marshal(&NotificationResponse{
			Message: msg.Message,
			Type:    msg.Type,
		})
		client.Send <- data
		context.Send(msg.Sender, &messages.Response{
			Message: "Notification sent.",
		})
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.New(corsConfig()))
	hub := socket.NewHub()
	go hub.Run()
	r.GET("/:userId", func(c *gin.Context) {
		socket.ServeWs(hub, c.Writer, c.Request, c.Param("userId"))
	})

	system := actor.NewActorSystem()
	remoteConfig := remote.Configure("127.0.0.1", 8092)
	remoting := remote.NewRemote(system, remoteConfig)
	remoting.Start()
	remoting.Register("notification-actor", actor.PropsFromProducer(func() actor.Actor {
		return &NotificationActor{
			hub: hub,
		}
	}))

	r.Run(":10000")
}

func corsConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	return corsConfig
}
