package service

import (
	"appstore/backend"
	"appstore/constants"
	"appstore/model"
	"reflect"

	"github.com/olivere/elastic/v7"
)

func SearchApps(title string, description string) ([]model.App, error) {
	if title == "" {
		return SearchAppsByDescription(description)
	}
	if description == "" {
		return SearchAppsByTitle(title)
	}
	query1 := elastic.NewMatchQuery("title", title)
	query2 := elastic.NewMatchQuery("description", description)
	query := elastic.NewBoolQuery().Must(query1, query2)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
	if err != nil {
		return nil, err
	}

	return getAppFromSearchResult(searchResult), nil
}

func SearchAppsByTitle(title string) ([]model.App, error) {
	query := elastic.NewMatchQuery("title", title)
	query.Operator("AND")
	if title == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
	if err != nil {
		return nil, err
	}

	return getAppFromSearchResult(searchResult), nil
}

func SearchAppsByDescription(description string) ([]model.App, error) {
	query := elastic.NewMatchQuery("description", description)
	query.Operator("AND")
	if description == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
	if err != nil {
		return nil, err
	}

	return getAppFromSearchResult(searchResult), nil
}

func SearchAppsByID(appID string) (*model.App, error) {
	query := elastic.NewMatchQuery("id", appID)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
	if err != nil {
		return nil, err
	}
	result := getAppFromSearchResult(searchResult)
	if len(result) == 1 {
		return &result[0], nil
	}
	return nil, nil
}

func getAppFromSearchResult(searchResult *elastic.SearchResult) []model.App {
	var ptyp model.App
	var apps []model.App
	for _, item := range searchResult.Each(reflect.TypeOf(ptyp)) {
		p := item.(model.App)
		apps = append(apps, p)
	}
	return apps
}