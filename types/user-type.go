package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserClient struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Password string `json:"password" bson:"password"`
}

type UserUpdateClient struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Password  string             `json:"password" bson:"password"`
	UserType  string             `json:"user_type" bson:"user_type"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at" bson:"updated_at"`
}

type Address struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Address1 string             `json:"address_1" bson:"address_1"`
	UserId   primitive.ObjectID `json:"user_id" bson:"user_id"`
	City     string             `json:"city" bson:"city"`
	Country  string             `json:"country" bson:"country"`
}

type AddressClient struct {
	Address1 string `json:"address_1" bson:"address_1"`
	UserId   string `json:"user_id" bson:"user_id"`
	City     string `json:"city" bson:"city"`
	Country  string `json:"country" bson:"country"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
