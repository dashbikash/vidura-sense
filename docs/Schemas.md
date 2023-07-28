# Web Crawler Schemas

## 1 HtmlPage

```json
{
    "id": "hashurl",
    "url": "google.com/search/about",
    "scheme": "https",
    "host":"google.com",
    "title":"Site Title",
    "meta": {
        "charset": "utf-8",
        "author":"bikash",
        "description":"description",
        "language":"en_IN",
        "viewport":"width=device-width, initial-scale=1"
    },
    "body":"body content",
    "links": ["link1","link2"
    ],
    "updated_on": "2023-07-17T10:00:00",
    "updated_by": {
        "proxy": "",
        "node_ip":""
    }
}
```

## 2 FeedItem

``` json
{
    "url":"url",
    "source_url":"url",
    "title":"title",
    "description":"description",
    "published_on":"2023-07-17T10:00:00",
    "updated_on": "2023-07-17T10:00:00",
    "updated_by": {
        "agent": "bot1",
        "proxy": "",
        "node_ip":""
    }
}
```
