package services

import (
	"context"

	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/items"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/queries"
)

type ItemsServiceInterface interface {
	Create(context.Context, items.Item) (*items.Item, error)
	Get(context.Context, string) (*items.Item, error)
	Search(context.Context, queries.EsQuery) ([]items.Item, error)
	Delete(context.Context, string) error
	Put(context.Context, items.Item)(*items.Item, error)
	Patch(context.Context, items.PartialUpdateItem, string)(*items.Item, error)
}

type itemsService struct{
	itemDao items.ItemDaoInterface
}

func NewItemsService(itemDao items.ItemDaoInterface) *itemsService {
	return &itemsService{itemDao: itemDao}
}

func (s *itemsService) Create(ctx context.Context, item items.Item) (*items.Item, error) {
	if err := s.itemDao.Save(ctx, item); err != nil{
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Get(ctx context.Context, id string) (*items.Item, error) {
	result, err := s.itemDao.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	result.Id = id

	return result, nil
}

func (s *itemsService) Search(ctx context.Context, query queries.EsQuery) ([]items.Item, error) {
	return s.itemDao.Search(ctx, query)
}

func (s *itemsService) Delete(ctx context.Context, id string) error {
	return s.itemDao.Delete(ctx, id)
}

func (s *itemsService) Put(ctx context.Context, item items.Item) (*items.Item, error) {
	if err := s.itemDao.Put(ctx, item); err != nil{
		return nil, err
	}
	return &item, nil
}

func (s *itemsService) Patch(ctx context.Context, item items.PartialUpdateItem, id string) (*items.Item, error) {
	if err := s.itemDao.Patch(ctx, item, id); err != nil{
		return nil, err
	}

	return s.itemDao.Get(ctx, id)
}