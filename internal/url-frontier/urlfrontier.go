package urlfrontier

import (
	"context"

	mongostore "github.com/dashbikash/vidura-sense/internal/datastorage/mongo-store"
	"github.com/dashbikash/vidura-sense/internal/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetNewUrls(limit int) []string {

	var result []struct{ URL string }
	cursor, err := mongostore.DefaultClient().Database().Collection(system.Config.Data.Mongo.Collections.Htmlpages).Aggregate(context.TODO(),
		mongo.Pipeline{
			{{Key: "$match", Value: bson.D{{Key: "body", Value: nil}}}},
			{{Key: "$sort", Value: bson.D{{Key: "updated_on", Value: 1}}}},
			{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "url", Value: 1}}}},
			{{Key: "$limit", Value: limit}},
		})
	if err != nil {
		system.Log.Error(err.Error())
		return nil
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		system.Log.Error(err.Error())
		return nil
	}
	urls := make([]string, 0)
	for _, v := range result {
		urls = append(urls, v.URL)
	}
	return urls

}
