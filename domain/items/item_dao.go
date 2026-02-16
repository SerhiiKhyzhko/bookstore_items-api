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

func (i *Item) Save(ctx context.Context) error {
	err := elasticsearch.Client.Index(ctx, indexItems, i.Id, i)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return item_errors.RequestTimeoutErr
		}
		return fmt.Errorf("save failed %w", err)
		//return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database error"))
	}
	return nil
}

func (i *Item) Get(ctx context.Context) error {
	result, err := elasticsearch.Client.Get(ctx, indexItems, i.Id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return item_errors.RequestTimeoutErr
		}
		return fmt.Errorf("get failed %w", err)
		//return rest_errors.NewInternalServerError(fmt.Sprintf("error when trying to get id %s", i.Id), errors.New("database error"))
	}

	if !result.Found {
		return item_errors.NotFoundErr
		//return rest_errors.NewNotFoundError(fmt.Sprintf("no item found with id %s", i.Id))
	}

	if err := json.Unmarshal(result.Source_, &i); err != nil {
		return item_errors.ParseErr
		//return rest_errors.NewInternalServerError("error parsing elasticsearch response", err)
	}

	i.Id = result.Id_

	return nil
}

func (i *Item) Search(ctx context.Context, query queries.EsQuery) ([]Item, error) {
	searchRequest, err := elasticsearch.Client.Search(ctx, indexItems, query.Build(), query.From, query.Size)
	if err != nil {
		return nil, fmt.Errorf("search failed %w", err)
		//return nil, rest_errors.NewInternalServerError("error when trying to search document", errors.New("database error"))
	}

	result := make([]Item, len(searchRequest.Hits.Hits))
	for index, hit := range searchRequest.Hits.Hits {
		var item Item
		if err := json.Unmarshal(hit.Source_, &item); err != nil {
			return nil, item_errors.ParseErr
			//return nil, rest_errors.NewInternalServerError("error when trying to parse response", errors.New("database error"))
		}
		item.Id = *hit.Id_
		result[index] = item
	}

	return result, nil
}

func (i *Item) Delete(ctx context.Context, id string) error {
	if err := elasticsearch.Client.Delete(ctx, indexItems, id); err != nil {
		if errors.Is(err, elasticsearch.ErrorNotFound) {
			return item_errors.NotFoundErr
			//return rest_errors.NewNotFoundError(fmt.Sprintf("document not found with such id %s in given index %s", id, indexItems))
		}
		return fmt.Errorf("delete failed %w", err)
		//return rest_errors.NewInternalServerError("error when trying to delete document", errors.New("database error"))
	}
	return nil
}

func (i *Item) Put(ctx context.Context) error {
	err := elasticsearch.Client.Update(ctx, indexItems, i.Id, i)
	if err != nil {
		return fmt.Errorf("update of entire item failed %w", err)
		//return rest_errors.NewInternalServerError("error when trying to fully update item", errors.New("database error"))
	}
	return nil
}

func (p *PartialUpdateItem) Patch(ctx context.Context, id string) error {
	err := elasticsearch.Client.Update(ctx, indexItems, id, p)
	if err != nil {
		return fmt.Errorf("update item`s field(s) failed %w", err)
		//return rest_errors.NewInternalServerError("error when trying to update item`s field(s)", errors.New("database error"))
	}
	return nil
}