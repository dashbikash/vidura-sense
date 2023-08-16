package entity

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cespare/xxhash"
	mongostore "github.com/dashbikash/vidura-sense/internal/datastorage/mongo-store"
	"github.com/dashbikash/vidura-sense/internal/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HtmlPage struct {
	ID     string `json:"_id" bson:"_id"`
	URL    string `json:"url"`
	Hash   string `json:"hash"`
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Title  string `json:"title" `

	Meta struct {
		Charset     string `json:"charset"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Keywords    string `json:"keywords"`
		Language    string `json:"language"`
		Viewport    string `json:"viewport"`
	} `json:"meta"`
	Body    string `json:"body"`
	Summary struct {
		Content string
		Nav     []struct {
			Href string
			Text string
		}
	} `json:"summary"`
	UpdatedOn primitive.DateTime `json:"updated_on" bson:"updated_on"`
	UpdatedBy struct {
		AppID  string `json:"app_id"`
		Proxy  string `json:"proxy"`
		NodeIP string `json:"node_ip" bson:"node_ip"`
	} `json:"updated_by" bson:"updated_by"`
	LockExpiry primitive.DateTime `json:"lock_expiry,omitempty" bson:"lock_expiry,omitempty"`
}

func (htmlPage *HtmlPage) Store() {
	htmlPage.ID = strconv.FormatUint(xxhash.Sum64String(htmlPage.URL), 10)
	htmlPage.Hash = strconv.FormatUint(xxhash.Sum64String(htmlPage.URL+"\n"+htmlPage.Body), 10)

	if ds := mongostore.DefaultClient(); ds != nil {
		result := ds.Collection(system.Config.Data.Mongo.Collections.Htmlpages).CreateOrReplaceOne(htmlPage, bson.D{{Key: "_id", Value: htmlPage.ID}})
		system.Log.Info(fmt.Sprintf("%d document updated.", result))
	}
}
func HtmlPageLockUpdate(targetUrl string) {

	if ds := mongostore.DefaultClient(); ds != nil {
		ds.Collection(system.Config.Data.Mongo.Collections.Htmlpages).FindOneAndUpdate(bson.M{"$set": bson.M{"lock_expiry": time.Now().Add(24 * time.Hour)}}, bson.D{{Key: "url", Value: targetUrl}})
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

func NewBlankHtmlPage(targetUrl string) interface{} {
	urlParsed, _ := url.Parse(targetUrl)
	httpScheme := urlParsed.Scheme
	urlParsed.Scheme = ""

	return struct {
		ID        string `json:"_id" bson:"_id"`
		URL       string
		Scheme    string
		Host      string
		UpdatedOn primitive.DateTime `json:"updated_on" bson:"updated_on"`
		UpdatedBy struct {
			AppID string `json:"app_id" bson:"app_id"`
		} `json:"updated_by" bson:"updated_by"`
	}{
		ID:        strconv.FormatUint(xxhash.Sum64String(strings.Trim(urlParsed.String(), "/")), 10),
		URL:       strings.Trim(urlParsed.String(), "/"),
		Scheme:    httpScheme,
		Host:      urlParsed.Host,
		UpdatedOn: primitive.DateTime(time.Now().Local().UnixMilli()),
		UpdatedBy: struct {
			AppID string `json:"app_id" bson:"app_id"`
		}{
			AppID: system.Config.Application.ID,
		},
	}
}

func NewBlankHtmlPageV2(targetUrl string) interface{} {

	return HtmlPage{
		ID:        strconv.FormatUint(xxhash.Sum64String(targetUrl), 10),
		URL:       targetUrl,
		UpdatedOn: primitive.DateTime(time.Date(1986, 7, 2, 0, 0, 0, 0, time.UTC).Unix()),

		UpdatedBy: struct {
			AppID  string `json:"app_id"`
			Proxy  string `json:"proxy"`
			NodeIP string `json:"node_ip" bson:"node_ip"`
		}{},
	}
}

func HtmlPageCreateBlankEntries(pages *[]interface{}) {

	if ds := mongostore.DefaultClient(); ds != nil {

		result := ds.Collection(system.Config.Data.Mongo.Collections.Htmlpages).InsertOrIgnore(*pages)
		system.Log.Info(fmt.Sprintf("%d blank document created.", result))
	}
}
