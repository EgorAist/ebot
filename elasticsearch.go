package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"sync"
)

type Item struct {
	Author  string `json:"author"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	PubDate string `json:"pubDate"`
	Text    string `json:"text"`
}

func ES(items []Item) error {
	var wg sync.WaitGroup

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}

	for _, item := range items {
		wg.Add(1)

		go func(item Item) {
			defer wg.Done()

			b, err := json.Marshal(item)
			if err != nil {
				log.Fatalf("Error marshal item: %s", err)
			}

			itemID := hashMD5(item)

			req := esapi.IndexRequest{
				Index:      "items",
				DocumentID: itemID,
				Body:       bytes.NewReader(b),
				Refresh:    "true",
			}

			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%s", res.Status(), itemID)
			} else {
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(item)
	}

	wg.Wait()

	return nil
}

func itemToBytes(item Item) []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(item.Title),
			[]byte(item.Author),
			[]byte(item.PubDate),
			[]byte(item.Link),
			[]byte(item.Text),
		},
		[]byte{},
	)

	return data
}

func hashMD5(item Item) string {
	hash := md5.Sum(itemToBytes(item))
	return hex.EncodeToString(hash[:])
}
