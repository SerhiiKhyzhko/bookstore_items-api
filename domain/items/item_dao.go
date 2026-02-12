package items

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SerhiiKhyzhko/bookstore_items-api/clients/elasticsearch"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/queries"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
)

const (
	indexItems = "items"
)

func (i *Item) Save() rest_errors.RestErr {
	err := elasticsearch.Client.Index(indexItems, i.Id, i)
	if err != nil {
		return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))
	}
	return nil
}

func (i *Item) Get() rest_errors.RestErr {
	result, err := elasticsearch.Client.Get(indexItems, i.Id)
	if err != nil {
		return rest_errors.NewInternalServerError(fmt.Sprintf("error when trying to get id %s", i.Id), errors.New("database error"))
	}

	if !result.Found {
		return rest_errors.NewNotFoundError(fmt.Sprintf("no item found with id %s", i.Id))
	}

	if err := json.Unmarshal(result.Source_, &i); err != nil {
		return rest_errors.NewInternalServerError("error parsing elasticsearch response", err)
	}

	i.Id = result.Id_

	return nil
}

func (i *Item) Search(query queries.EsQuery) ([]Item, rest_errors.RestErr) {
	searchRequest, err := elasticsearch.Client.Search(indexItems, query.Build(), query.From, query.Size)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to search document", errors.New("database error"))
	}

	result := make([]Item, len(searchRequest.Hits.Hits))
	for index, hit := range searchRequest.Hits.Hits {
		var item Item
		if err := json.Unmarshal(hit.Source_, &item); err != nil {
			return nil, rest_errors.NewInternalServerError("error when trying to parse response", errors.New("database error"))
		}
		item.Id = *hit.Id_
		result[index] = item
	}

	return result, nil
}

func (i *Item) Delete(id string) rest_errors.RestErr {
	if err := elasticsearch.Client.Delete(indexItems, id); err != nil {
		if errors.Is(err, elasticsearch.ErrorNotFound) {
			return rest_errors.NewNotFoundError(fmt.Sprintf("document not found with such id %s in given index %s", id, indexItems))
		}
		return rest_errors.NewInternalServerError("error when trying to delete document", errors.New("database error"))
	}
	return nil
}
