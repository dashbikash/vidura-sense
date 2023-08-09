package entity

import (
	"time"
)

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
