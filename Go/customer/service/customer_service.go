package service

import (
	"context"
	"fmt"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/config"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/customer_actor"
	messages "github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/message"
	"github.com/AT-SmFoYcSNaQ/AT2023/Go/customer/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type CustomerService struct {
	logger        *zap.Logger
	db            *gorm.DB
	customerActor *customer_actor.CustomerActor
}

func CreateCustomerService(customerActor *customer_actor.CustomerActor, zapLogger *zap.Logger) *CustomerService {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		zapLogger.Fatal(err.Error())
		return nil
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
	databaseURL := loadConfig.DatabaseURL
	db, err := gorm.Open(postgres.Open(databaseURL+"&application_name=$ docs_simplecrud_gorm"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.AutoMigrate(&model.Customer{})
	if err != nil {
		zapLogger.Fatal(err.Error())
		return nil
	}

	return &CustomerService{zapLogger, db, customerActor}
}

type OrderBody struct {
	ItemId   string  `json:"itemId"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Item struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Name     string             `json:"name"`
	Quantity uint32             `json:"quantity"`
	Price    float64            `json:"price"`
}

func (service *CustomerService) GetItemsFromMongoDatabase() *[]Item {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb://localhost:27017").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		service.logger.Error(err.Error())
		return nil
	}

	pingErr := client.Ping(context.TODO(), readpref.Primary())
	if pingErr != nil {
		service.logger.Error(pingErr.Error())
	}

	service.logger.Info("Connected to mongo database")

	collection := client.Database("inventory").Collection("items")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		service.logger.Error(err.Error())
		return nil
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			service.logger.Error(err.Error())
		}
	}(cursor, context.TODO())

	var items []Item
	for cursor.Next(context.TODO()) {
		var item Item
		err := cursor.Decode(&item)
		if err != nil {
			service.logger.Error(err.Error())
			return nil
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil
	}

	return &items
}

func (service *CustomerService) Insert(customer *model.Customer) (id string, err error) {
	createdId := service.db.Create(customer)
	if createdId.Error != nil {
		service.logger.Error(createdId.Error.Error())
		return "", createdId.Error
	}

	return fmt.Sprint(createdId), nil
}

func (service *CustomerService) GetByEmail(email string) (*model.Customer, error) {
	var customer model.Customer
	if err := service.db.First(&customer, "email = ?", email).Error; err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	return &customer, nil
}

func (service *CustomerService) GetById(id string) (*model.Customer, error) {
	var customer model.Customer
	if err := service.db.First(&customer, "id = ?", id).Error; err != nil {
		service.logger.Error(err.Error())
		return nil, err
	}

	return &customer, nil
}

func (service *CustomerService) EmailExists(email string) bool {
	var customer model.Customer
	if err := service.db.First(&customer, "email = ?", email).Error; err != nil {
		return false
	}

	return true
}

func (service *CustomerService) Order(orderBody *OrderBody, customerId string) error {
	balance := service.GetBalanceByCustomerId(customerId)
	order := &messages.ReceiveOrder_Request{
		UserId:         customerId,
		ItemId:         orderBody.ItemId,
		Quantity:       int32(orderBody.Quantity),
		AccountBalance: balance,
		PricePerItem:   orderBody.Price,
	}

	pid := service.customerActor.Spawn()
	service.customerActor.Send(pid, order)
	service.logger.Info("Message sent to customer-actor")
	return nil
}

func (service *CustomerService) GetBalanceByCustomerId(customerId string) float64 {
	var customer model.Customer
	if err := service.db.First(&customer, "id = ?", customerId).Error; err != nil {
		return 0
	}
	return customer.Balance
}
