package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/EgorAist/ebot/models"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"strings"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		XMLName xml.Name  `xml:"channel"`
		Items   []XMLItem `xml:"item"`
	} `xml:"channel"`
}

type XMLItem struct {
	XMLName     xml.Name `xml:"item"`
	Guid        string   `xml:"guid"`
	Author      string   `xml:"author"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Category    string   `xml:"category"`
}

func ParsePage(url string) ([]*models.Item, error) {
	var rss RSS

	byteValue, err := LoadPage(url)
	if err != nil {
		return nil, errors.Wrap(err, "in ParsePage")
	}

	err = xml.Unmarshal(byteValue, &rss)
	if err != nil {
		return nil, err
	}

	var items []*models.Item

	for _, item := range rss.Channel.Items {
		date, _ := time.Parse(time.RFC1123Z, item.PubDate)
		items = append(items, &models.Item{
			Author:  item.Author,
			Title:   item.Title,
			Link:    item.Link,
			PubDate: date,
		})
	}

	return items, nil
}

func ParsItems(items []*models.Item) ([]*models.Item, error) {

	for _, item := range items {
		text, err := ParseItemPage(item.Link)
		if err != nil {
			return nil, err
		}

		item.Text = text
	}

	return items, nil
}

func PrintItems(items []*models.Item)  {
	for _, item := range items{
		fmt.Println("Title: " 	+ item.Title)
		fmt.Println("Author: " 	+ item.Author)
		fmt.Println("Link: " 	+ item.Link)
		fmt.Println("Text: " 	+ item.Text)
	}
}

func ParseItemPage(url string) (string, error) {
	byteValue, err := LoadPage(url)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(byteValue)

	z := html.NewTokenizer(reader)

	text := ""

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return text, err

		case html.StartTagToken:
			t := z.Token()

			if t.Data != "div" {
				continue
			}
			for _, attr := range t.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "js-topic__text") {
					text = parseTopicText(z)
				}
			}
		}
	}
}

func parseTopicText(z *html.Tokenizer) string {
	text := ""
	countDivs := 1
	open := false
	for countDivs != 0 {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return text

		case html.StartTagToken:
			t := z.Token()

			if t.Data == "div" {
				countDivs += 1
			}

			if t.Data == "p" {
				open = true
				tt = z.Next()
				if tt != html.StartTagToken {
					t = z.Token()
					text += t.Data + " "
				}
			}

			if t.Data == "a" && open == true {
				tt = z.Next()
				t = z.Token()
				text += t.Data
			}

		case html.TextToken:
			if open {
				t := z.Token()
				text += t.Data + " "
			}

		case html.EndTagToken:
			t := z.Token()

			if t.Data == "div" {
				countDivs -= 1
			}

			if t.Data == "p"{
				open = false
			}
		}
	}

	return text
}
