package model

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	// USD currency
	USD = "USD"

	// RUB currency
	RUB = "RUB"

	// EUR currency
	EUR = "EUR"
)

// Product model declaration
type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Price       float32            `bson:"price"`
	Currency    string             `bson:"currency"`
}
