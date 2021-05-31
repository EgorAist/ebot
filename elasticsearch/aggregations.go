package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/EgorAist/ebot/models"
	"github.com/olivere/elastic/v7"
	"log"
)

func (es *ESClient) TermAggregationByField(field string) ([]models.Aggregation, error) {
	aggregationQuery := elastic.NewTermsAggregation().
		Field(addKeyWord(field)).
		Size(30).
		OrderByCountDesc()

	result, err := es.client.Search().
		Index(indexName).
		Aggregation(indexName, aggregationQuery).
		Do(context.Background())
	if err != nil {
		e, ok := err.(*elastic.Error)
		if ok {
			log.Printf("Got error from elastic: %s", e.Details)
		}
		return []models.Aggregation{}, err
	}

	rawMsg := result.Aggregations[indexName]

	ar := elastic.AggregationBucketKeyItems{}

	err = json.Unmarshal(rawMsg, &ar)
	if err != nil {
		return nil, err
	}

	var termsAggregations []models.Aggregation

	for _, item := range ar.Buckets {
		termsAggregations = append(termsAggregations, models.Aggregation{
			Key:   item.Key,
			Count: item.DocCount,
		})
	}
	return termsAggregations, nil
}

func (es *ESClient) CardinalityAggregationByField(field string) (float64, error) {
	aggregationQuery := elastic.NewCardinalityAggregation().Field(addKeyWord(field))

	result, err := es.client.Search().
		Index(indexName).
		Aggregation(indexName, aggregationQuery).
		Do(context.Background())
	if err != nil {
		e, ok := err.(*elastic.Error)
		if ok {
			log.Printf("Got error from elastic %s", e.Details)
		}
		return 0, err
	}

	rawMsg := result.Aggregations[indexName]

	var ar elastic.AggregationValueMetric

	err = json.Unmarshal(rawMsg, &ar)
	if err != nil {
		return 0, err
	}

	return *ar.Value, nil
}

func (es *ESClient) DateHistogramAggregation() ([]models.Aggregation, error) {
	dailyAggregation := elastic.NewDateHistogramAggregation().
		Field("pubDate").
		CalendarInterval(dateInterval).
		Format(dateFormat)

	result, err := es.client.Search().
		Index(indexName).
		Aggregation(indexName, dailyAggregation).
		Do(context.Background())
	if err != nil {
		e, ok := err.(*elastic.Error)
		if ok {
			log.Printf("Got error from elastic %s", e.Details)
		}
		return []models.Aggregation{}, err
	}

	hist, found := result.Aggregations.Histogram(indexName)
	if !found {
		return []models.Aggregation{}, errors.New("Not found ")
	}

	var dateHistogramAggregations []models.Aggregation

	for _, bucket := range hist.Buckets {
		dateHistogramAggregations = append(dateHistogramAggregations, models.Aggregation{
			Key:   *bucket.KeyAsString,
			Count: bucket.DocCount,
		})
	}

	return dateHistogramAggregations, nil
}

func addKeyWord(field string) string {
	if field == "title" || field == "description" {
		return field + ".keyword"
	}

	return field
}
