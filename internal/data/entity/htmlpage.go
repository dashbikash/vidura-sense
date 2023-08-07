package entity

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cespare/xxhash"
	mongostore "github.com/dashbikash/vidura-sense/internal/datastorage/mongo-store"
	"github.com/dashbikash/vidura-sense/internal/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HtmlPage struct {
	ID     string `json:"_id" bson:"_id"`
	URL    string `json:"url"`
	Hash   string `json:"hash"`
	Scheme string `json:"scheme" `
	Host   string `json:"host"`
	Title  string `json:"title" `

	Meta struct {
		Charset     string `json:"charset"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Language    string `json:"language"`
		Viewport    string `json:"viewport"`
	} `json:"meta"`
	Body      string    `json:"body"`
	Links     []string  `json:"links"`
	UpdatedOn time.Time `json:"updated_on" bson:"updated_on"`
	UpdatedBy struct {
		AppID  string `json:"app_id"`
		Proxy  string `json:"proxy"`
		NodeIP string `json:"node_ip" bson:"node_ip"`
	} `json:"updated_by" bson:"updated_by"`
}

func (htmlPage *HtmlPage) StoreUpdated() {
	htmlPage.ID = strconv.FormatUint(xxhash.Sum64String(htmlPage.URL), 10)
	htmlPage.Hash = strconv.FormatUint(xxhash.Sum64String(htmlPage.URL+"\n"+htmlPage.Body), 10)

	if ds := mongostore.DefaultClient(); ds != nil {
		result := ds.CreateOrReplaceOne(system.Config.Data.Mongo.Collections.Htmlpages, htmlPage, bson.D{{Key: "_id", Value: htmlPage.ID}})
		system.Log.Info(fmt.Sprintf("%d document updated.", result))
	}
}
func HtmlPageCreateBlankEntries(pages *[]interface{}) {

	if ds := mongostore.DefaultClient(); ds != nil {

		result := ds.InsertOrIgnore(system.Config.Data.Mongo.Collections.Htmlpages, *pages)
		system.Log.Info(fmt.Sprintf("%d blank document created.", result))
	}
}

type UrlMeta struct {
	URL         string `bson:"_id"`
	LastUpdated time.Time
}

func GetUrlMeta(host string, urlMeta *UrlMeta) {

	if ds := mongostore.DefaultClient(); ds != nil {

		cursor, err := ds.Database().Collection(system.Config.Data.Mongo.Collections.Htmlpages).Aggregate(context.TODO(),
			mongo.Pipeline{
				{{Key: "$match", Value: bson.D{{Key: "host", Value: host}}}},
				{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$host"}, {Key: "LastUpdated", Value: bson.D{{Key: "$max", Value: "$updated_on"}}}}}},
			},
		)
		if err != nil {
			panic(err)
		}
		if cursor.Next(context.TODO()) {
			if err = cursor.Decode(&urlMeta); err != nil {
				system.Log.Fatal(err.Error())
			}
		}

	}
}

func BlankHtmlPage(targetUrl string) interface{} {

	return struct {
		ID        string `json:"_id" bson:"_id"`
		URL       string
		UpdatedOn time.Time `json:"updated_on"`
		UpdatedBy struct {
			AppID string `json:"app_id" bson:"app_id"`
		} `json:"updated_by" bson:"updated_by"`
	}{
		ID:        strconv.FormatUint(xxhash.Sum64String(targetUrl), 10),
		URL:       targetUrl,
		UpdatedOn: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedBy: struct {
			AppID string `json:"app_id" bson:"app_id"`
		}{
			AppID: system.Config.Application.ID,
		},
	}
}
