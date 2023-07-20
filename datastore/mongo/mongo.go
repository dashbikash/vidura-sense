package mongo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dashbikash/vidura-sense/provider"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	log    = provider.GetLogger()
	config = provider.GetConfig()
)

func getConnection() *mongo.Client {
	log.Debugf("Connecting to Mongodb %s", config.Data.MongoUrl)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Data.MongoUrl))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	return client
}

func QueryData() {
	log.Debugf("Connecting to Mongodb %s", config.Data.MongoUrl)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Data.MongoUrl))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("sample").Collection("employees")

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"lname", "Carter"}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found")
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}
