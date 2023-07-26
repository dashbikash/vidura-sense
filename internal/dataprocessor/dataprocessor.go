package dataprocessor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dashbikash/vidura-sense/internal/common"
	mongostore "github.com/dashbikash/vidura-sense/internal/datastore/mongo-store"
	"github.com/dashbikash/vidura-sense/internal/datatype"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	log    = common.GetLogger()
	config = common.GetConfig()
)

func StoreHtml(htmlPage *datatype.HtmlPage) {
	ds := mongostore.DefaultClient()
	if ds != nil {
		result := ds.CreateOrReplaceOne(config.Data.Mongo.Collections.Htmlpages, htmlPage, bson.D{{"url", htmlPage.URL}})
		log.Info(fmt.Sprintf("%d document updated.", result))
	}
}

func SanitizeText(text string) string {
	var sb strings.Builder
	cr := regexp.MustCompile(`\n+`)
	tab := regexp.MustCompile(`\t+`)

	text = cr.ReplaceAllString(text, "\n")
	text = tab.ReplaceAllString(text, "\t")
	txtBlocks := strings.Split(text, "\n")

	for _, blk := range txtBlocks {

		blk = strings.TrimSpace(blk)
		if len(blk) > 0 {
			sb.WriteString(blk + "\n")
		}
	}
	return sb.String()

}
func SanitizeHtml(htm string) string {

	p := bluemonday.UGCPolicy()

	html := p.Sanitize(htm)

	return html
}
