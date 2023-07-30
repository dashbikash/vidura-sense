package entity

import (
	"fmt"
	"time"

	mongostore "github.com/dashbikash/vidura-sense/internal/datastorage/mongo-store"
	"github.com/dashbikash/vidura-sense/internal/system"
	"go.mongodb.org/mongo-driver/bson"
)

type HtmlPage struct {
	URL    string `json:"url"`
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Title  string `json:"title"`
	Meta   struct {
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
		Proxy  string `json:"proxy"`
		NodeIP string `json:"node_ip" bson:"node_ip"`
	} `json:"updated_by" bson:"updated_by"`
}

type FeedItem struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	SourceURL   string    `json:"source_url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedOn time.Time `json:"published_on"`
	UpdatedOn   time.Time `json:"updated_on"`
	UpdatedBy   struct {
		Agent  string `json:"agent"`
		Proxy  string `json:"proxy"`
		NodeIP string `json:"node_ip"`
	} `json:"updated_by"`
}

func (htmlPage *HtmlPage) StoreHtml() {

	if ds := mongostore.DefaultClient(); ds != nil {
		result := ds.CreateOrReplaceOne(system.Config.Data.Mongo.Collections.Htmlpages, htmlPage, bson.D{{Key: "url", Value: htmlPage.URL}})
		system.Log.Info(fmt.Sprintf("%d document updated.", result))
	}
}
