package config

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/fanchann/belajar_testing_mongodb/helpers"
)

func NewMongoConnection(v *viper.Viper) *mongo.Database {
	ctx := context.Background()

	url := v.GetString("database.url")

	creds := options.Credential{
		Username: v.GetString("database.username"),
		Password: v.GetString("database.password"),
	}

	clientOptions := options.Client().ApplyURI(url).SetAuth(creds)
	clientOptions.ApplyURI(url)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		helpers.ErrorLogger(err)

	}

	err = client.Connect(ctx)
	if err != nil {
		helpers.ErrorLogger(err)
	}

	return client.Database(v.GetString("database.name"))
}
