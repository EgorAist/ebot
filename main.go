package main

import (
	"fmt"
	"github.com/EgorAist/ebot/elasticsearch"
	"github.com/EgorAist/ebot/models"
	"github.com/EgorAist/ebot/parser"
)

const url = "https://lenta.ru/rss/last24"
const (
	wordSearch                 = "Боевой"
	wordTermAggregation        = "author"
	wordCardinalityAggregation = "author"
)


func main()  {
	client, err := elasticsearch.NewESClient()
	if err != nil{
		fmt.Println(err)
		return
	}

	err = client.CreateIndex()
	if err != nil{
		fmt.Println(err)
		return
	}

	items, err := client.GetLastItems()
	if err != nil{
		fmt.Println(err)
		return
	}

	existsItems := map[string]*models.Item{}
	for _, item := range items {
		existsItems[item.Link] = item
	}

	newItems, err := parser.ParsePage(url)
	if err != nil{
		fmt.Println(err)
		return
	}

	var newItemsForPars []*models.Item

	for _, newItem := range newItems {
		if _, ok := existsItems[newItem.Link]; !ok {
			newItemsForPars = append(newItemsForPars, newItem)
		}
	}

	parsedItems, err := parser.ParsItems(newItemsForPars)
	if err != nil{
		fmt.Println(err)
		return
	}

	err = client.InsertItems(parsedItems)
	if err != nil{
		fmt.Println(err)
		return
	}

	itemsByKey, err := client.FullTextSearch(wordSearch)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("\nFind by word, ", wordSearch)
	for _, item := range itemsByKey {
		fmt.Println("Title: " 	+ item.Title)
		fmt.Println("PubDate: " + item.PubDate.String())
		fmt.Println("Author: " 	+ item.Author)
		fmt.Println("Link: " 	+ item.Link)
		fmt.Println("Text: " 	+ item.Text)
	}

	termAggregation, err := client.TermAggregationByField(wordTermAggregation)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("\nTerm aggregation by %s\n", wordTermAggregation)
	for _, t := range termAggregation {
		fmt.Printf("%s: %v\n", wordTermAggregation, t.Key)
		fmt.Printf("Count: %d\n", t.Count)
	}

	cardinalityAggregation, err := client.CardinalityAggregationByField(wordCardinalityAggregation)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("\nCardinality aggregation by %s: %.0f\n", wordCardinalityAggregation, cardinalityAggregation)

	dataHistogram, err := client.DateHistogramAggregation()
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("\nDate histogram aggregation")
	for _, t := range dataHistogram {
		fmt.Printf("Key: %v\n", t.Key)
		fmt.Printf("Count: %d\n", t.Count)
	}
}
