package services

import (
	"context"

	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/items"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/queries"
)

var ItemsService ItemsServiceInterface = &itemsService{}

type ItemsServiceInterface interface {
	Create(context.Context, items.Item) (*items.Item, error)
	Get(context.Context, string) (*items.Item, error)
	Search(context.Context, queries.EsQuery) ([]items.Item, error)
	Delete(context.Context, string) error
	Put(context.Context, items.Item)(*items.Item, error)
	Patch(context.Context, items.PartialUpdateItem, string)(*items.Item, error)
}

type itemsService struct{}

func (s *itemsService) Create(ctx context.Context, item items.Item) (*items.Item, error) {
	if err := item.Save(ctx); err != nil{
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Get(ctx context.Context, id string) (*items.Item, error) {
	item := items.Item{Id: id}

	if err := item.Get(ctx); err != nil {
		return nil, err 
	}
	return &item, nil
}

func (s *itemsService) Search(ctx context.Context, query queries.EsQuery) ([]items.Item, error) {
	dao := items.Item{}
	return dao.Search(ctx, query)
}

func (s *itemsService) Delete(ctx context.Context, id string) error {
	dao := items.Item{}
	return dao.Delete(ctx, id)
}

func (s *itemsService) Put(ctx context.Context, item items.Item) (*items.Item, error) {
	if err := item.Put(ctx); err != nil{
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Patch(ctx context.Context, item items.PartialUpdateItem, id string) (*items.Item, error) {
	if err := item.Patch(ctx, id); err != nil{
		return nil, err
	}

	return s.Get(ctx, id)
}