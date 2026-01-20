package items

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SerhiiKhyzhko/bookstore_items-api/clients/elasticsearch"
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
