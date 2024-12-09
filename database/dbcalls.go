package database

import (
	"context"
	"ecommerce-project/constant"
	"ecommerce-project/types"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Insert inserts a new document into the specified MongoDB collection.
// Parameters:
// - data: The document to be inserted (as an interface{}).
// - collectionName: The name of the MongoDB collection where the document will be stored.
// Returns:
// - InsertedID: The ID of the newly inserted document.
// - error: Error if any issue occurs during the operation.
func (mgr *manager) Insert(data interface{}, collectionName string) (interface{}, error) {
	// Retrieve the collection object from the database connection
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	// Insert the provided data into the collection
	result, err := orgCollection.InsertOne(context.TODO(), data)

	// Check for errors in the insertion process
	if err != nil {
		return nil, err
	}

	// Log the result of the operation for debugging
	log.Println(result)

	// Return the ID of the inserted document
	return result.InsertedID, nil
}

// GetSingleRecordByEmail retrieves a single document matching the provided email from a specified collection.
// Parameters:
// - email: The email address to filter by.
// - collectionName: The name of the MongoDB collection to search.
// Returns:
// - *types.Verification: The verification record matching the email.
func (mgr *manager) GetSingleRecordByEmail(email string, collectionName string) *types.Verification {
	// Create an empty Verification object to hold the response
	resp := &types.Verification{}

	// Define the filter criteria for the query
	filter := bson.D{{Key: "email", Value: email}}

	// Retrieve the collection object
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	// Execute the query and decode the result into the response object
	err := orgCollection.FindOne(context.TODO(), filter).Decode(&resp)

	// Log any error that occurs during the query
	fmt.Println(err)

	// Return the verification record (or an empty object if no match is found)
	return resp
}

// UpdateVerification updates the verification details of a user in a specified collection.
// Parameters:
// - data: The updated Verification object.
// - collectionName: The name of the MongoDB collection to update.
// Returns:
// - error: Error if any issue occurs during the update operation.
func (mgr *manager) UpdateVerification(data types.Verification, collectionName string) error {
	// Retrieve the collection object
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	// Define the filter and update criteria
	filter := bson.D{{Key: "email", Value: data.Email}}
	update := bson.D{{Key: "$set", Value: data}}

	// Execute the update operation
	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)

	// Return any error encountered during the operation
	return err
}

// UpdateEmailVerifiedStatus marks a user's email as verified in the database.
// Parameters:
// - req: The Verification object containing the updated status.
// - collectionName: The name of the MongoDB collection to update.
// Returns:
// - error: Error if any issue occurs during the update operation.
func (mgr *manager) UpdateEmailVerifiedStatus(req types.Verification, collectionName string) error {
	// Retrieve the collection object
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	// Define the filter and update criteria
	filter := bson.D{{Key: "email", Value: req.Email}}
	update := bson.D{{Key: "$set", Value: req}}

	// Execute the update operation
	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)

	// Return any error encountered during the operation
	return err
}

// GetSingleRecordByEmailForUser retrieves a user record matching the provided email from a specified collection.
// Parameters:
// - email: The email address to filter by.
// - collectionName: The name of the MongoDB collection to search.
// Returns:
// - *types.User: The user record matching the email.
func (mgr *manager) GetSingleRecordByEmailForUser(email, collectionName string) *types.User {
	// Create an empty User object to hold the response
	resp := &types.User{}

	// Retrieve the collection object
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	// Define the filter criteria for the query
	filter := bson.D{{Key: "email", Value: email}}

	// Execute the query and decode the result into the response object
	_ = orgCollection.FindOne(context.TODO(), filter).Decode(&resp)

	// Return the user record (or an empty object if no match is found)
	return resp
}


func (mgr *manager) GetListProducts(page, limit, offset int, collectionName string) ([]types.Product, int64, error) {
	// Calculate skip value based on page and limit
	skip := (page - 1) * limit
	if offset > 0 {
			skip = offset
	}

	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	// Set find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Query documents
	cur, err := orgCollection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
			return nil, 0, err
	}
	defer cur.Close(context.TODO())

	// Decode documents
	var products []types.Product
	if err := cur.All(context.TODO(), &products); err != nil {
			return nil, 0, err
	}

	// Count total documents
	count, err := orgCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
			return nil, 0, err
	}

	return products, count, nil
}


func (mgr *manager) SearchProduct(page, limit, offset int, search,collectionName string) ([]types.Product, int64, error) {
	// Calculate skip value based on page and limit
	skip := (page - 1) * limit
	if offset > 0 {
			skip = offset
	}

	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	// Set find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	searchFilter := bson.M{}

	if len(search) >= 3 {
		searchFilter["$or"] = []bson.M{
			{"name" : primitive.Regex{Pattern: ".*" + search + ".*", Options : "i"}},
			{"description" : primitive.Regex{Pattern: ".*" + search + ".*", Options : "i"}},
		}
	}


	// Query documents
	cur, err := orgCollection.Find(context.TODO(), searchFilter, findOptions)
	if err != nil {
			return nil, 0, err
	}
	defer cur.Close(context.TODO())

	// Decode documents
	var products []types.Product
	if err := cur.All(context.TODO(), &products); err != nil {
			return nil, 0, err
	}

	// Count total documents
	count, err := orgCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
			return nil, 0, err
	}

	return products, count, nil
}

func (mgr *manager)	GetSingleProductById(id primitive.ObjectID, collectionName string)(types.Product, error){
	filter := bson.D{{Key :"_id", Value: id}}
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	var product types.Product
	err := orgCollection.FindOne(context.TODO(), filter).Decode(&product)

	return product, err
}


func (mgr *manager) UpdateProduct(p types.Product, colllectionName string)error{
	orgCollection := mgr.connection.Database(constant.Database).Collection(colllectionName)
	filter := bson.D{{Key: "_id", Value: p.Id}}
	update := bson.D{{Key: "$set", Value: p}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)

	return err 
}

func (mgr *manager)	DeleteProduct(id primitive.ObjectID, collectionName string)error{
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "_id", Value: id}}

	_, err := orgCollection.DeleteOne(context.TODO(), filter)
	return err
}

func (mgr *manager)	GetSingleAddress(id primitive.ObjectID, collectionName string)(types.Address, error){
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "user_id", Value: id}}
	var address types.Address
	err := orgCollection.FindOne(context.TODO(), filter).Decode(&address)
	return address, err
}

func (mgr *manager)	GetSingleUserByUserId(id primitive.ObjectID, collectionName string)(types.User, error){
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "_id", Value: id}}
	var user types.User
	err := orgCollection.FindOne(context.TODO(), filter).Decode(&user)
	return user, err
}

func (mgr *manager)	GetCartObjectById(id primitive.ObjectID, collectionName string)(types.Cart, error){
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "_id", Value: id}}
	var cart types.Cart
	err := orgCollection.FindOne(context.TODO(), filter).Decode(&cart)
	return cart, err
}

func (mgr *manager)	GetCartObjectListForUser(userID primitive.ObjectID, collectionName string)([]types.Cart, error){
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	// Define the filter for the user's carts
	filter := bson.D{{Key: "userId", Value: userID}}

	// Find multiple documents
	cursor, err := orgCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Decode documents
	var cartItems []types.Cart
	if err := cursor.All(context.TODO(), &cartItems); err != nil {
			return nil, err
	}

	return cartItems, nil
}


func (mgr *manager) UpdateUser(u types.User, collectionName string)error{
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "_id", Value: u.Id}}
	update := bson.D{{Key: "$set", Value: u}}
	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (mgr *manager)	UpdateCartToCheckout(userID primitive.ObjectID, collectionName string)error{
	// orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	// filter := bson.D{{Key: "_id", Value: c.Id}}
	// update := bson.D{{Key: "$set", Value: c}}
	// _, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	// return err
	// Get the collection
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	// Define the filter to match documents with the given userId
	filter := bson.D{{Key: "user_id", Value: userID}}

	// Define the update operation to set "checkout" to true
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "checkout", Value: true}}}}

	// Update all matching documents
	_, err := orgCollection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
