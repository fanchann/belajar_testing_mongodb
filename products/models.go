package products

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductModel struct {
	ID       primitive.ObjectID `bson:"_id,omiempty"`
	Name     string             `bson:"name"`
	Quantity uint               `bson:"quantity"`
	Price    uint               `bson:"price"`
}
