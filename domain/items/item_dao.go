package items

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SerhiiKhyzhko/bookstore_items-api/clients/elasticsearch"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/queries"
	"github.com/SerhiiKhyzhko/bookstore_items-api/item_errors"
)

const (
	indexItems = "items"
)

type ItemDaoInterface interface {
	Save(context.Context, Item) error
	Get(context.Context, string) (*Item, error)
	Search(context.Context, queries.EsQuery) ([]Item, error)
	Delete(context.Context, string) error
	Put(context.Context, Item) error
	Patch(context.Context, PartialUpdateItem, string) error
}

type itemDaoStruct struct {
	client elasticsearch.EsClientInterface
}

func NewItemDao(esClient elasticsearch.EsClientInterface) *itemDaoStruct {
	return &itemDaoStruct{client: esClient}
}

func (d *itemDaoStruct) Save(ctx context.Context, item Item) error {
	err := d.client.Index(ctx, indexItems, item.Id, item)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return item_errors.RequestTimeoutErr
		}
		return fmt.Errorf("save failed %w", err)
	}
	return nil
}

func (d *itemDaoStruct) Get(ctx context.Context, id string) (*Item, error) {
	var item Item
	result, err := d.client.Get(ctx, indexItems, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, item_errors.RequestTimeoutErr
		}
		return nil, fmt.Errorf("get failed %w", err)
	}

	if !result.Found {
		return nil, item_errors.NotFoundErr
	}

	if err := json.Unmarshal(result.Source_, &item); err != nil {
		return nil, item_errors.ParseErr
	}

	item.Id = result.Id_

	return &item, nil
}

func (d *itemDaoStruct) Search(ctx context.Context, query queries.EsQuery) ([]Item, error) {
	searchRequest, err := d.client.Search(ctx, indexItems, query.Build(), query.From, query.Size)
	if err != nil {
		return nil, fmt.Errorf("search failed %w", err)
	}
	result := make([]Item, len(searchRequest.Hits.Hits))
	for index, hit := range searchRequest.Hits.Hits {
		var item Item
		if err := json.Unmarshal(hit.Source_, &item); err != nil {
			return nil, item_errors.ParseErr
		}
		item.Id = *hit.Id_
		result[index] = item
	}

	return result, nil
}

func (d *itemDaoStruct) Delete(ctx context.Context, id string) error {
	itemFound, err := d.client.Delete(ctx, indexItems, id)
	if err != nil {
		return fmt.Errorf("delete failed %w", err)
	}
	if !itemFound {
		return item_errors.NotFoundErr
	}
	
	return nil
}

func (d *itemDaoStruct) Put(ctx context.Context, item Item) error {
	itemFound, err := d.client.Update(ctx, indexItems, item.Id, item)
	if err != nil {
		return fmt.Errorf("update of entire item failed %w", err)
	}
	if !itemFound {
		return item_errors.NotFoundErr
	}

	return nil
}

func (d *itemDaoStruct) Patch(ctx context.Context, partialUpdateItem PartialUpdateItem, id string) error {
	itemFound, err := d.client.Update(ctx, indexItems, id, partialUpdateItem)
	if err != nil {
		
		return fmt.Errorf("update item`s field(s) failed %w", err)
	}
	if !itemFound {
		return item_errors.NotFoundErr
	}

	return nil
}