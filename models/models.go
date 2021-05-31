package models

import "time"

type Item struct {
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Link    string    `json:"link"`
	PubDate time.Time `json:"pubDate"`
	Text    string    `json:"text"`
}

type Aggregation struct {
	Key   interface{}
	Count int64
}

