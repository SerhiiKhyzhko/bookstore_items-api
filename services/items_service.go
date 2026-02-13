package services

import (
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/items"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/queries"
)

var ItemsService ItemsServiceInterface = &itemsService{}

type ItemsServiceInterface interface {
	Create(items.Item) (*items.Item, rest_errors.RestErr)
	Get(string) (*items.Item, rest_errors.RestErr)
	Search(queries.EsQuery) ([]items.Item, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
	Put(items.Item)(*items.Item, rest_errors.RestErr)
	Patch(items.PartialUpdateItem, string)(*items.Item, rest_errors.RestErr)
}

type itemsService struct{}

func (s *itemsService) Create(item items.Item) (*items.Item, rest_errors.RestErr) {
	if err := item.Save(); err != nil{
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Get(id string) (*items.Item, rest_errors.RestErr) {
	item := items.Item{Id: id}

	if err := item.Get(); err != nil {
		return nil, err 
	}
	return &item, nil
}

func (s *itemsService) Search(query queries.EsQuery) ([]items.Item, rest_errors.RestErr) {
	dao := items.Item{}
	return dao.Search(query)
}

func (s *itemsService) Delete(id string) rest_errors.RestErr {
	dao := items.Item{}
	return dao.Delete(id)
}

func (s *itemsService) Put(item items.Item) (*items.Item, rest_errors.RestErr) {
	if err := item.Put(); err != nil{
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Patch(item items.PartialUpdateItem, id string) (*items.Item, rest_errors.RestErr) {
	if err := item.Patch(id); err != nil{
		return nil, err
	}

	return s.Get(id)
}