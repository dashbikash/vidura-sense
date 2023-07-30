package mongostore

import (
	"context"

	"github.com/dashbikash/vidura-sense/internal/system"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx = context.Background()
)

type MongoStore struct {
	mongoClient *mongo.Client
}

func DefaultClient() *MongoStore {
	system.Log.Debug("Connecting to Mongodb " + system.Config.Data.Mongo.MongoUrl)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(system.Config.Data.Mongo.MongoUrl))
	if err != nil {
		system.Log.Error(err.Error())
		return nil
	}

	return &MongoStore{mongoClient: client}
}

func (ds *MongoStore) CreateOrReplaceOne(collection string, data interface{}, filter interface{}) int64 {
	defer ds.dispose()
	opts := options.Replace().SetUpsert(true)
	result, err := ds.mongoClient.Database(system.Config.Data.Mongo.Database).Collection(collection).ReplaceOne(ctx, filter, data, opts)
	if err != nil {
		system.Log.Error(err.Error())
		return 0
	}

	return result.UpsertedCount
}
func (ds *MongoStore) dispose() {

	if err := ds.mongoClient.Disconnect(ctx); err != nil {
		system.Log.Error(err.Error())
	}

}
