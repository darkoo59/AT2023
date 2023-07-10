package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/AT-SmFoYcSNaQ/AT2023/Go/inventory/model"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/order/messages"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	dbUri           = "mongodb://localhost:27017"
	itemsCollection *mongo.Collection
)

var (
	collectionItems = []interface{}{
		model.Item{Name: "Samsung Galaxy S23 Ultra", Price: 2000, Quantity: 50},
		model.Item{Name: "Lenovo IP3", Price: 1477.99, Quantity: 20},
		model.Item{Name: "Sony PlayStation 5", Price: 600.99, Quantity: 30},
		model.Item{Name: "AMD Ryzen 5 5600x", Price: 301, Quantity: 40},
		model.Item{Name: "Intel Core i5-11400F", Price: 180.99, Quantity: 100},
		model.Item{Name: "Logitech G915", Price: 340.99, Quantity: 17},
		model.Item{Name: "Razer Viper Mini", Price: 45, Quantity: 4},
		model.Item{Name: "Xwave usb hub 4-port", Price: 10, Quantity: 323},
		model.Item{Name: "Apple macbook pro 16", Price: 5000, Quantity: 2},
	}
)

type InventoryActor struct {
}

func (act *InventoryActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		log.Println("Inventory actor started")
	case messages.CheckAvailability_Request:
		available, item := CheckItemAvailability(msg.ItemId, int(msg.Quantity))
		ctx.Send(msg.Sender, messages.CheckAvailability_Response{
			OrderId:     "????????????????????????????",
			IsAvailable: available,
			Quantity:    msg.Quantity,
			ItemName:    item.Name,
			ItemPrice:   float32(item.Price),
		})
	}
}

func NewInventoryActor() actor.Actor {
	return &InventoryActor{}
}

func CheckItemAvailability(itemId string, quantity int) (bool, *model.Item) {
	itemObjectId, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		panic("Invalid item object id")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var item model.Item
	err = itemsCollection.FindOne(ctx, bson.D{{Key: "_id", Value: itemObjectId}}).Decode(&item)
	if err != nil {
		log.Printf("No item with id: %v", itemId)
		return false, &item
	}

	if uint32(quantity) > item.Quantity {
		log.Printf("There are no enough items, item id: %v", itemId)
		return false, &item
	}

	return true, &item
}

func ConnectToDb() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	pingErr := client.Ping(context.TODO(), readpref.Primary())
	if pingErr != nil {
		log.Fatal("Ping, mongoDB: ", err.Error())
	}

	log.Printf("Connected to mongo database")

	itemsCollection = client.Database("inventory").Collection("items")

}

func SeedItems() {
	itemsCollection.Drop(context.TODO())
	result, err := itemsCollection.InsertMany(context.TODO(), collectionItems)
	if err != nil {
		panic(err)
	}
	log.Printf("\n%v items seeded successfully", len(result.InsertedIDs))
}

func main() {
	port := 20001
	port1 := 8098
	if len(os.Args) > 2 {
		port, _ = strconv.Atoi(os.Args[1])
		port1, _ = strconv.Atoi(os.Args[2])
	}

	system := actor.NewActorSystem()
	remoteConfig := remote.Configure("127.0.0.1", port)

	cp := automanaged.NewWithConfig(1*time.Second, port1, "localhost:8098")
	clusterKind := cluster.NewKind(
		"inventory-actor",
		actor.PropsFromProducer(NewInventoryActor),
	)
	lookup := disthash.New()
	clusterConfig := cluster.Configure("cluster-inventory", cp, lookup, remoteConfig, cluster.WithKinds(clusterKind))
	clust := cluster.New(system, clusterConfig)

	clust.StartMember()
	defer clust.Shutdown(false)

	finishChan := make(chan os.Signal, 1)
	signal.Notify(finishChan, os.Interrupt)

	ConnectToDb()
	SeedItems()

	<-finishChan
}
