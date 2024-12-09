package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Verification struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email  string             `json:"email" bson:"email"`
	Otp    int64              `json:"otp" bson:"otp"`
	Status bool               `json:"status" bson:"status"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
}
