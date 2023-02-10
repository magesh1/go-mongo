package db

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection

var ctx = context.TODO()

func DbInit() {

	address := viper.GetString("mongo.address")

	clientOption := options.Client().ApplyURI(address)
	client, err := mongo.Connect(ctx, clientOption)

	if err != nil {
		log.Error().Msg("Error connecting to MongoDB database: " + err.Error())
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Error().Msg("Error pinging to MongoDB database: " + err.Error())
		return
	}

	Collection = client.Database("mydb").Collection("urlshort")

	log.Info().Msg("Connected to MongoDB database")

}
