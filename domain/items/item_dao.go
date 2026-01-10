package items

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/SerhiiKhyzhko/bookstore_items-api/clients/elasticsearch"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
)

const (
	indexItems = "items"
)

func (i *Item) Save() rest_errors.RestErr {
	sellerId := strconv.Itoa(int(i.Seller))
	err := elasticsearch.Client.Index(indexItems, i.Id, sellerId, i)
	if err != nil {
		fmt.Println("error is :------", err)
		return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))
	}
	return nil
}
