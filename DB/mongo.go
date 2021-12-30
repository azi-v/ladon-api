package DB

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client

func InitMongoConn(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = MongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
}

type CollResources struct {
	ID       string `bson:"id,omitempty"`
	PolicyID string `bson:"policy_id,omitempty"`
	Pesource string `bson:"pesource,omitempty"`
}
type CollSubject struct {
	ID       string `bson:"id,omitempty"`
	PolicyID string `bson:"policy_id,omitempty"`
	Subject string `bson:"subject,omitempty"`
}
