package entity

import "html/template"

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"guid"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Post struct {
	ID      int    `json:"ID"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
	Link    string `json:"Link"`
	PubTime int64  `json:"PubTime"`
}

type PostForPublic struct {
	Title       string `json:"title"`
	ContentHTML template.HTML
	Link        template.URL
	PubTime     string `json:"PubTime"`
}
