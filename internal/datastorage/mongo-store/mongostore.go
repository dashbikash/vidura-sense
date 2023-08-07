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
	database    *mongo.Database
}

func DefaultClient() *MongoStore {
	system.Log.Debug("Connecting to Mongodb " + system.Config.Data.Mongo.MongoUrl)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(system.Config.Data.Mongo.MongoUrl))
	if err != nil {
		system.Log.Error(err.Error())
		return nil
	}

	return &MongoStore{mongoClient: client, database: client.Database(system.Config.Data.Mongo.Database)}
}
func DefaultClientWithDatabase(db string) *MongoStore {
	system.Log.Debug("Connecting to Mongodb " + system.Config.Data.Mongo.MongoUrl)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(system.Config.Data.Mongo.MongoUrl))
	if err != nil {
		system.Log.Error(err.Error())
		return nil
	}

	return &MongoStore{mongoClient: client, database: client.Database(db)}
}

func (ds *MongoStore) Database() *mongo.Database {
	return ds.database
}

func (ds *MongoStore) CreateOrReplaceOne(collection string, data interface{}, filter interface{}) int64 {
	defer ds.dispose()
	opts := options.Replace().SetUpsert(true)
	result, err := ds.database.Collection(collection).ReplaceOne(ctx, filter, data, opts)
	if err != nil {
		system.Log.Error(err.Error())
		return 0
	}

	return result.UpsertedCount
}

func (ds *MongoStore) InsertOrIgnore(collection string, data []interface{}) int {
	defer ds.dispose()
	opts := &options.InsertManyOptions{}
	opts.SetOrdered(false)

	result, _ := ds.database.Collection(collection).InsertMany(ctx, data, opts)

	return len(result.InsertedIDs)
}

func (ds *MongoStore) Aggregate(collection string, pipeline interface{}, dtype *[]interface{}) {
	defer ds.dispose()
	cursor, err := ds.database.Collection(collection).Aggregate(ctx, pipeline)
	if err != nil {
		system.Log.Error(err.Error())
		return
	}

	if err = cursor.All(context.TODO(), dtype); err != nil {
		system.Log.Error(err.Error())
		return
	}
}
func (ds *MongoStore) dispose() {

	if err := ds.mongoClient.Disconnect(ctx); err != nil {
		system.Log.Error(err.Error())
	}

}
