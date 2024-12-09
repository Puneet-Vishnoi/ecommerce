package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cart struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Checkout  bool               `json:"checkout,omitempty" bson:"checkout"`
}

type CartClient struct {
	UserId    string `json:"user_id" bson:"user_id"`
	ProductID string `json:"product_id" bson:"product_id"`
}
