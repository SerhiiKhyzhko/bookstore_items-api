package elasticsearch

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/optype"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/result"
)

// ---------------------------------------------------------
// 1. КОНСТАНТИ ТА МАПІНГ (СХЕМА)
// ---------------------------------------------------------

const (
	IndexItems = "items"
)

var(
	ErrorNotFound = errors.New("Document not found") 
)

// налаштування типів та шардів.
const itemMapping = `{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    },
    "mappings": {
        "properties": {
            "title": {
                "type": "text",
                "analyzer": "standard"
            },
            "description": {
                "properties": {
                    "plain_text": {
                        "type": "text",
                        "analyzer": "standard"
                    },
                    "html": {
                        "type": "text",
                        "analyzer": "standard" 
                    }
                }
            },
            "pictures": {
                "properties": {
                    "id": { "type": "long" },
                    "url": { "type": "keyword" }
                }
            },
            "seller": { 
                "type": "long" 
            },
            "video": {
                "type": "keyword"
            },
            "price": {
                "type": "float"
            },
            "available_quantity": {
                "type": "integer"
            },
            "sold_quantity": {
                "type": "integer"
            },
            "status": {
                "type": "keyword"
            },
            "date_created": {
                "type": "date",
                "format": "strict_date_optional_time"
            }
        }
    }
}`

type EsClientInterface interface {
	Init(string)
	Index(string, string, any) error
	Get(string, string) (*get.Response, error)
	Search(string, *types.Query, *int, *int) (*search.Response, error)
	Delete(string, string) error
	Update(string, string, any) error
}

type esClient struct {
	client *elasticsearch.TypedClient
}

var (
	Client EsClientInterface = &esClient{}
)

// ---------------------------------------------------------
// 2. ІНІЦІАЛІЗАЦІЯ (INIT & ENSURE INDEX)
// ---------------------------------------------------------

func (c *esClient) Init(addreses string) {
	cfg := elasticsearch.Config{
		Addresses: strings.Split(addreses, ";"),
		Transport: &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // Для локальної розробки
			ResponseHeaderTimeout: 10 * time.Second,
		},
		// Якщо потрібно залогінитися під іншим користувачем
		// Username: "elastic",
		// Password: "changeme",
	}

	var err error
	c.client, err = elasticsearch.NewTypedClient(cfg)
	if err != nil {
		logger.Error("Error creating elasticsearch typed client", err)
		panic(err)
	}

	res, err := c.client.Info().Do(context.Background())
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintf("Elasticsearch client connected. Claster: %s, Version %s", res.ClusterName, res.Version)
	logger.Info(msg)

	c.ensureIndexCreated()
}

func (c *esClient) ensureIndexCreated() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := c.client.Indices.Exists(IndexItems).Do(ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("error when check the index %s", IndexItems), err)
		return
	}

	if exists {
		logger.Info(fmt.Sprintf("Indedx %s already exists", IndexItems))
		return
	}

	req, err := create.NewRequest().FromJSON(itemMapping)

	if err != nil {
		logger.Error("error parsing itemMapping JSON", err)
		panic(err)
	}

	res, err := c.client.Indices.Create(IndexItems).
		Request(req).
		Do(ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("failed to create index %s", IndexItems), err)
		panic(err)
	}

	if res.Acknowledged {
		logger.Info(fmt.Sprintf("index %s created successfully with settings and mappings", IndexItems))
	}
}

// ---------------------------------------------------------
// 3. МЕТОДИ РОБОТИ З ДАНИМИ (INDEX)
// ---------------------------------------------------------

func (c *esClient) Index(index string, id string, doc any) error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := c.client.Index(index).Id(id).OpType(optype.Create).Document(doc).Do(ctx) //OpType - захист від перезапису

	if err != nil {
		logger.Error("error connecting to elasticsearch", err)
		return err
	}

	logger.Info(fmt.Sprintf("document indexed: %s, result: %s", res.Id_, res.Result))
	return nil
}

func (c *esClient) Get(index string, Id string) (*get.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := c.client.Get(index, Id).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to get id %s", Id), err)
		return nil, err
	}

	return res, nil
}

func (c *esClient) Search(index string, query *types.Query, from *int, size *int) (*search.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.client.Search().Index(index).Request(
		&search.Request{
			Query: query,
			From:  from,
			Size:  size,
		}).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when trying to search documents in index %s", index), err)
		return nil, err
	}
	return result, nil
}

func (c *esClient) Delete(index string, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := c.client.Delete(index, id).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to delete document with id %s from index %s", id, index), err)
		return err
	}
	if res.Result == result.Notfound {
		return ErrorNotFound
	}
	return nil
}

func (c * esClient) Update(index string, id string, doc any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := c.client.Update(index, id).Doc(doc).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to update document with id %s from index %s", id, index), err)
		return err
	}

	return nil
}
