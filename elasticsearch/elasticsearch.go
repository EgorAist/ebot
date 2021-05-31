package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/EgorAist/ebot/models"

	"github.com/olivere/elastic/v7"
	"os"

	"log"
)

const esAddress = "http://127.0.0.1:9200"
const indexName = "items"

const (
	dateFormat   = "dd.MM.YYYY"
	dateInterval = "1d"
)

type ESClient struct {
	client *elastic.Client
}

func NewESClient() (ESClient, error) {
	client, err := elastic.NewSimpleClient(elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)))
	if err != nil {
		return ESClient{}, err
	}

	return ESClient{client: client}, nil
}

func (es *ESClient) CreateIndex() error {
	ctx := context.Background()
	exists, err := es.client.
		IndexExists(indexName).
		Do(ctx)
	if err != nil {
		log.Println("Error with index exists request", err)
	}

	if exists {
		return nil
	}

	_, err = es.client.CreateIndex(indexName).Do(ctx)
	if err != nil {
		log.Println("Error with create index", err)
		return err
	}

	err = es.createMapping()
	if err != nil {
		return err
	}

	return nil
}

func (es *ESClient) createMapping() error {
	_, err := es.client.PutMapping().
		Index(indexName).
		BodyString(PutMapping).
		Do(context.Background())
	if err != nil {
		log.Println("Error with put mapping", err)
		return err
	}

	return nil
}

func (es *ESClient) InsertItems(items []*models.Item) error {
	for _, item := range items {
		_, err := es.client.Index().
			Index(indexName).
			Id(item.Link).
			BodyJson(&item).
			Do(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *ESClient) GetLastItems() ([]*models.Item, error) {
	query := elastic.NewMatchAllQuery()
	res, err := es.client.Search().
		Index(indexName).
		Query(query).
		Size(100).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var items []*models.Item

	for _, hit := range res.Hits.Hits {
		var t *models.Item
		err := json.Unmarshal(hit.Source, &t)
		if err != nil {
			return nil, err
		}

		items = append(items, t)
	}

	return items, nil
}

func (es *ESClient) FullTextSearch(key string) ([]models.Item, error) {
	queryString := elastic.NewQueryStringQuery(key)

	res, err := es.client.Search().
		Index(indexName).
		Query(queryString).
		Do(context.Background())
	if err != nil {
		return []models.Item{}, err
	}

	if res.Hits.TotalHits.Value == 0 {
		return []models.Item{}, errors.New("Not found ")
	}

	var items []models.Item

	for _, hit := range res.Hits.Hits {
		var t models.Item
		err := json.Unmarshal(hit.Source, &t)
		if err != nil {
			return nil, err
		}

		items = append(items, t)
	}

	return items, nil
}
