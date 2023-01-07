package backend

import (
	"context"
	"fmt"

	"appstore/constants"

	"github.com/olivere/elastic/v7"
)

var (
	ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct {
	client *elastic.Client
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}

func InitElasticsearchBackend() {
	// new ElasticsearchBackend object to be used for ES operation
	client, err := elastic.NewClient(
		elastic.SetURL(constants.ES_URL),
		elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_PASSWORD))
	if err != nil {
		panic(err)
	}

	exists, err := client.IndexExists(constants.APP_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		mapping := `{
				"mappings": {
					"properties": {
						"id":       { "type": "keyword" },
						"user":     { "type": "keyword" },
						"title":      { "type": "text"},
						"description":  { "type": "text" },
						"price":      { "type": "keyword", "index": false },
						"url":     { "type": "keyword", "index": false }
					}
				}
			}`
		_, err := client.CreateIndex(constants.APP_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		mapping := `{
						 "mappings": {
							 "properties": {
								"username": {"type": "keyword"},
								"password": {"type": "keyword"},
								"age": {"type": "long", "index": false},
								"gender": {"type": "keyword", "index": false}
							 }
						}
					}`
		_, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Indexes are created.")

	ESBackend = &ElasticsearchBackend{client: client}
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
	_, err := backend.client.Index().Index(index).Id(id).BodyJson(i).Do(context.Background())
	return err
}

func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
	_, err := backend.client.DeleteByQuery().Index(index).Query(query).Pretty(true).Do(context.Background())
	return err;
}
