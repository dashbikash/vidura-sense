package datatype

import "time"

type HtmlDocument struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Protocol string `json:"protocol"`
	Title    string `json:"title"`
	Meta     struct {
		Charset     string `json:"charset"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Language    string `json:"language"`
		Viewport    string `json:"viewport"`
	} `json:"meta"`
	TextBlocks []struct {
		Heading string `json:"heading"`
		Content string `json:"content"`
	} `json:"text_blocks"`
	Links []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"links"`
	UpdatedOn time.Time `json:"updated_on"`
	UpdatedBy struct {
		Agent  string `json:"agent"`
		Proxy  string `json:"proxy"`
		NodeIP string `json:"node_ip"`
	} `json:"updated_by"`
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
