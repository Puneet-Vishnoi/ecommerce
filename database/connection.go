package database

import (
	"context"
	"ecommerce-project/constant"
	"ecommerce-project/types"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type manager struct {
	connection *mongo.Client
	ctx        context.Context
	cancel     context.CancelFunc
}

var Mgr *manager

type Manager interface {
	Insert(interface{}, string) (interface{}, error)
	GetSingleRecordByEmail(string, string) *types.Verification
	UpdateVerification(types.Verification, string) error
	UpdateEmailVerifiedStatus(types.Verification, string) error
	GetSingleRecordByEmailForUser(string, string) types.Verification
	GetListProducts(int, int, int, string)([]types.Product, int64, error)
	SearchProduct(int, int, int, string, string)([]types.Product, int64, error)
	GetSingleProductById(primitive.ObjectID, string)(types.Product, error)
	UpdateProduct(types.Product, string)error
	DeleteProduct(primitive.ObjectID, string)error
	GetSingleAddress(primitive.ObjectID, string)(types.Address, error)
	GetSingleUserByUserId(primitive.ObjectID, string)(types.User, error)
	UpdateUser(types.User, string) error
	GetCartObjectById(primitive.ObjectID, string)(types.Cart, error)
	GetCartObjectListForUser(primitive.ObjectID, string)([]types.Cart, error)
	UpdateCartToCheckout(types.Cart, string)error
}

// ConnectDb connects to the MongoDB database and initializes the global manager.
func ConnectDb() {
	uri := os.Getenv("BD_HOST")
	if uri == "" {
		uri = constant.MDBUri // Fall back to constant URI if not found in environment
	}

	// Create the client options using ApplyURI
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", uri))



	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx,clientOptions)


	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Test the connection by pinging the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Printf("Successfully connected to the database at %s", uri)

	// Initialize the global Mgr variable
	Mgr = &manager{
		connection: client,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Close gracefully closes the MongoDB connection and cleans up resources
func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	// Cancel the context and close the client connection
	defer cancel()

	// Disconnect from MongoDB
	err := client.Disconnect(ctx)
	if err != nil {
		log.Printf("Error while disconnecting MongoDB client: %v", err)
	} else {
		log.Println("MongoDB connection closed successfully.")
	}
}
