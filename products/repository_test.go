package products_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"github.com/fanchann/belajar_testing_mongodb/products"
)

const (
	dbName   = "db"
	collName = "products"
)

var (
	ctx = context.Background()

	productParam = &products.ProductModel{
		Name:     "PRODUCT TEST",
		Quantity: 100,
		Price:    500000,
	}

	productModel = &products.ProductModel{
		ID:       primitive.NewObjectID(),
		Name:     productParam.Name,
		Quantity: productParam.Quantity,
		Price:    productParam.Price,
	}
)

func TestProductsRepository(t *testing.T) {
	mongoTest := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mongoTest.Run("success insert user", func(mt *mtest.T) {

		produtImpl := products.NewProductImpl(mt.DB.Collection(collName))
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(),
		)

		response, err := produtImpl.InsertProduct(ctx, productModel)
		assert.Nil(t, err)
		assert.Equal(t, response, productModel)

	})

	mongoTest.Run("success find product by id", func(mt *mtest.T) {
		produtImpl := products.NewProductImpl(mt.DB.Collection(collName))
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				dbName+"."+collName,
				mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: productModel.ID},
					{Key: "name", Value: productModel.Name},
					{Key: "quantity", Value: productModel.Quantity},
					{Key: "price", Value: productModel.Price},
				},
			),
		)
		response, err := produtImpl.FindProductWithID(ctx, productModel.ID.Hex())
		assert.Nil(t, err)
		assert.Equal(t, response, productModel)
	})

	mongoTest.Run("product not found", func(mt *mtest.T) {
		productImpl := products.NewProductImpl(mt.Coll)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(
			mtest.WriteError{
				Index:   1,
				Code:    1100,
				Message: "object not found",
			},
		),
		)

		response, err := productImpl.FindProductWithID(ctx, "xxxx")
		assert.Nil(t, response)
		assert.Error(t, err)
	})

	mongoTest.Run("delete product", func(mt *mtest.T) {
		productImpl := products.NewProductImpl(mt.Coll)
		mt.AddMockResponses(
			bson.D{
				{"ok", 1},
				{"acknowledged", true},
			},
		)

		err := productImpl.DeleteProduct(ctx, productModel.ID.Hex())
		assert.Nil(t, err)
	})

	mongoTest.Run("delete product but product not found", func(mt *mtest.T) {
		productImpl := products.NewProductImpl(mt.Coll)
		mt.AddMockResponses(
			mtest.CreateWriteErrorsResponse(
				mtest.WriteError{
					Index:   -1,
					Code:    -1,
					Message: "object not found",
				},
			),
		)

		err := productImpl.DeleteProduct(ctx, "xxxx")
		assert.Error(t, err)
	})

	mongoTest.Run("update product", func(mt *mtest.T) {
		productImpl := products.NewProductImpl(mt.Coll)
		mt.AddMockResponses(
			bson.D{
				{"ok", 1},
				{"value", bson.D{
					{Key: "_id", Value: productModel.ID},
					{Key: "name", Value: productModel.Name},
					{Key: "quantity", Value: productModel.Quantity},
					{Key: "price", Value: productModel.Price},
				},
				},
			},
		)

		response, err := productImpl.UpdateProduct(ctx, productModel.ID.Hex(), productParam)
		assert.Nil(t, err)
		assert.Equal(t, productModel, response)
	})

	mongoTest.Run("update product but product not found", func(mt *mtest.T) {
		productImpl := products.NewProductImpl(mt.Coll)
		mt.AddMockResponses(
			mtest.CreateWriteErrorsResponse(
				mtest.WriteError{
					Index:   -1,
					Code:    -1,
					Message: "object not found",
				},
			),
		)

		response, err := productImpl.UpdateProduct(ctx, "xxx", &products.ProductModel{})
		assert.Error(t, err)
		assert.Equal(t, &products.ProductModel{}, response)
	})

}
