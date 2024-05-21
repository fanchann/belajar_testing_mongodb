package products

import (
	"context"
	"fmt"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type newProductImpl struct {
	collection *mongo.Collection
}

func NewProductImpl(c *mongo.Collection) IProduct {
	return &newProductImpl{
		collection: c,
	}
}

func (n *newProductImpl) DeleteProduct(ctx context.Context, idProduct string) error {
	productObjectID, _ := primitive.ObjectIDFromHex(idProduct)

	_, err := n.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: productObjectID}})
	return err
}

func (n *newProductImpl) UpdateProduct(ctx context.Context, idProduct string, product *ProductModel) (*ProductModel, error) {
	productObjectId, _ := primitive.ObjectIDFromHex(idProduct)
	var productModel ProductModel

	filter := bson.D{{Key: "_id", Value: productObjectId}}
	product.ID = productObjectId
	updateProductData := bson.D{{Key: "$set", Value: product}}
	err := n.collection.FindOneAndUpdate(ctx, filter, updateProductData).Decode(&productModel)

	return product, err
}

func (n *newProductImpl) FindProductWithID(ctx context.Context, idProduct string) (*ProductModel, error) {
	var product ProductModel

	// Convert idProduct to ObjectID if it is in string format
	objID, err := primitive.ObjectIDFromHex(idProduct)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	err = n.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (n *newProductImpl) FindProductsWithPagination(ctx context.Context, page int64, pageSize int64) ([]ProductModel, error) {

	// Calculate the number of documents to skip
	skip := (page - 1) * pageSize

	// Set the options for pagination
	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	findOptions.SetSkip(skip)

	cursor, err := n.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error finding products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []ProductModel
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("error decoding products: %w", err)
	}

	return products, nil
}

func (n *newProductImpl) FindAllProduct(ctx context.Context) []ProductModel {
	cursor, err := n.collection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)
	if err != nil {
		return nil
	}

	var products []ProductModel
	if err := cursor.All(ctx, &products); err != nil {
		return nil
	}
	return products
}

func (n *newProductImpl) InsertProduct(ctx context.Context, product *ProductModel) (*ProductModel, error) {
	_, err := n.collection.InsertOne(ctx, product)
	return product, err
}
