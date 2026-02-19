package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/optype"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/result"
)

type EsClientInterface interface {
	Index(context.Context, string, string, any) error
	Get(context.Context, string, string) (*get.Response, error)
	Search(context.Context, string, *types.Query, *int, *int) (*search.Response, error)
	Delete(context.Context, string, string) (bool, error)
	Update(context.Context, string, string, any) (bool, error)
}

type esClient struct {
	client *elasticsearch.TypedClient
}

func NewEsClient(client *elasticsearch.TypedClient) *esClient {
	return &esClient{client: client}
}

func (c *esClient) Index(ctx context.Context, index string, id string, doc any) error {

	esCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	res, err := c.client.Index(index).Id(id).OpType(optype.Create).Document(doc).Do(esCtx) //OpType - захист від перезапису

	if err != nil {
		logger.Error("error connecting to elasticsearch", err)
		return err
	}

	logger.Info(fmt.Sprintf("document indexed: %s, result: %s", res.Id_, res.Result))
	return nil
}

func (c *esClient) Get(ctx context.Context, index string, Id string) (*get.Response, error) {
	esCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	res, err := c.client.Get(index, Id).Do(esCtx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to get id %s", Id), err)
		return nil, err
	}

	return res, nil
}

func (c *esClient) Search(ctx context.Context, index string, query *types.Query, from *int, size *int) (*search.Response, error) {
	esCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := c.client.Search().Index(index).Request(
		&search.Request{
			Query: query,
			From:  from,
			Size:  size,
		}).Do(esCtx)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when trying to search documents in index %s", index), err)
		return nil, err
	}
	return result, nil
}

func (c *esClient) Delete(ctx context.Context, index string, id string) (bool, error) {
	esCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	res, err := c.client.Delete(index, id).Do(esCtx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to delete document with id %s from index %s", id, index), err)
		return true, err
	}
	
	if res.Result == result.Notfound {
		return false, nil
	}
	
	return true, nil
}

func (c * esClient) Update(ctx context.Context, index string, id string, doc any) (bool, error) {
	esCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := c.client.Update(index, id).Doc(doc).Do(esCtx)
	if err != nil {
		var e *types.ElasticsearchError
		if errors.As(err, &e) && e.Status == 404 {
			return false, nil
		}
		logger.Error(fmt.Sprintf("error when trying to update document with id %s from index %s", id, index), err)
		return true, err
	}

	return true, nil
}
