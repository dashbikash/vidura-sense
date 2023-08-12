package urlfrontier

import (
	"context"
	"time"

	mongostore "github.com/dashbikash/vidura-sense/internal/datastorage/mongo-store"
	"github.com/dashbikash/vidura-sense/internal/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetNewUrls(limit int) []string {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"$and": bson.A{
			bson.M{"hash": nil},
			bson.M{"$or": bson.A{
				bson.M{"lock_expiry": nil},
				bson.M{"lock_expiry": bson.M{"$lt": time.Now()}},
			}},
		}}}},
		{{Key: "$sort", Value: bson.M{"updated_on": 1}}},
		{{Key: "$project", Value: bson.M{"_id": 0, "url": 1, "scheme": 1}}},
		{{Key: "$limit", Value: limit}},
	}

	var result []struct {
		URL    string
		Scheme string
	}

	cursor, err := mongostore.DefaultClient().Database().Collection(system.Config.Data.Mongo.Collections.Htmlpages).Aggregate(context.TODO(),
		pipeline)
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
		urls = append(urls, v.Scheme+"://"+v.URL)
	}
	return urls

}

func GetExclusiveDomainNewUrls(limit int, filterDomains []string) []string {

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"$and": bson.A{
			bson.M{"hash": nil},
			bson.M{"host": bson.M{"$in": filterDomains}},
			bson.M{"$or": bson.A{
				bson.M{"lock_expiry": nil},
				bson.M{"lock_expiry": bson.M{"$lt": time.Now()}},
			}},
		}}}},
		{{Key: "$sort", Value: bson.M{"updated_on": 1}}},
		{{Key: "$project", Value: bson.M{"_id": 0, "url": 1, "scheme": 1}}},
		{{Key: "$limit", Value: limit}},
	}

	var result []struct {
		URL    string
		Scheme string
	}

	cursor, err := mongostore.DefaultClient().Database().Collection(system.Config.Data.Mongo.Collections.Htmlpages).Aggregate(context.TODO(),
		pipeline)
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
		urls = append(urls, v.Scheme+"://"+v.URL)
	}
	return urls

}
