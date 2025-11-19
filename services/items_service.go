package services

import (
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/items"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
)

var ItemsService ItemsServiceInterface = &itemsService{}

type ItemsServiceInterface interface {
	Create(items.Item) (*items.Item, *rest_errors.RestErr)
	Get(items.Item) (*items.Item, *rest_errors.RestErr)
}

type itemsService struct{}

func (s *itemsService) Create(items.Item) (*items.Item, *rest_errors.RestErr) {
	return nil, nil
}

func (s *itemsService) Get(items.Item) (*items.Item, *rest_errors.RestErr) {
	return nil, nil
}
