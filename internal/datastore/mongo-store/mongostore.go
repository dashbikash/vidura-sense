package mongostore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dashbikash/vidura-sense/internal/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	log    = common.GetLogger()
	config = common.GetConfig()
	ctx    = context.Background()
)

type MongoStore struct {
	mongoClient *mongo.Client
}

func DefaultClient() *MongoStore {
	log.Debug("Connecting to Mongodb " + config.Data.Mongo.MongoUrl)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Data.Mongo.MongoUrl))
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return &MongoStore{mongoClient: client}
}

func (ds *MongoStore) CreateOrReplaceOne(collection string, data interface{}, filter interface{}) int64 {
	defer ds.dispose()
	opts := options.Replace().SetUpsert(true)
	result, err := ds.mongoClient.Database(config.Data.Mongo.Database).Collection(collection).ReplaceOne(ctx, filter, data, opts)
	if err != nil {
		log.Error(err.Error())
		return 0
	}

	return result.MatchedCount
}
func (ds *MongoStore) dispose() {

	if err := ds.mongoClient.Disconnect(ctx); err != nil {
		log.Error(err.Error())
	}

}
func QueryData() {
	log.Debug("Connecting to Mongodb " + config.Data.Mongo.MongoUrl)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Data.Mongo.MongoUrl))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("sample").Collection("employees")

	var result bson.M
	err = coll.FindOne(ctx, bson.D{{"lname", "Carter"}}).Decode(&result)
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
