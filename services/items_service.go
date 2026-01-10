package services

import (
	"fmt"

	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/items"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
)

var ItemsService ItemsServiceInterface = &itemsService{}

type ItemsServiceInterface interface {
	Create(items.Item) (*items.Item, rest_errors.RestErr)
	Get(items.Item) (*items.Item, rest_errors.RestErr)
}

type itemsService struct{}

func (s *itemsService) Create(item items.Item) (*items.Item, rest_errors.RestErr) {
	fmt.Println("item is:------", item)
	if err := item.Save(); err != nil{
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Get(items.Item) (*items.Item, rest_errors.RestErr) {
	return nil, nil
}
