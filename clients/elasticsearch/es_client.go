package elasticsearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/joho/godotenv"
)

// ---------------------------------------------------------
// 1. КОНСТАНТИ ТА МАПІНГ (СХЕМА)
// ---------------------------------------------------------

const (
	IndexItems = "items"
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
	Init()
	Index(index string, id string, sellerID string, doc any) rest_errors.RestErr
}

type esClient struct {
	client *elasticsearch.Client
}

var (
	Client EsClientInterface = &esClient{}
)

// ---------------------------------------------------------
// 3. ІНІЦІАЛІЗАЦІЯ (INIT & ENSURE INDEX)
// ---------------------------------------------------------

func getESAddresses() []string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addreses := strings.Split(os.Getenv("ADDRESSES"), ";")
	return addreses
}

func (c *esClient) Init() {
	cfg := elasticsearch.Config{
		Addresses: getESAddresses(),
		Transport: &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // Для локальної розробки
			ResponseHeaderTimeout: 10 * time.Second,
		},
		// Якщо потрібно залогінитися під іншим користувачем
		// Username: "elastic",
		// Password: "changeme",
	}

	var err error
	c.client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	res, err := c.client.Info()
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		panic(fmt.Sprintf("Elasticsearch connection failed: %s", res.Status()))
	}

	logger.Info("Elasticsearch client connected successfully")

	c.ensureIndexCreated()
}

func (c *esClient) ensureIndexCreated() {
	ctx := context.Background()

	exists, err := c.client.Indices.Exists([]string{IndexItems})
	if err != nil {
		panic(err)
	}

	if exists.StatusCode != 404 {
		return
	}

	res, err := c.client.Indices.Create(
		IndexItems,
		c.client.Indices.Create.WithBody(strings.NewReader(itemMapping)),
		c.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		panic(fmt.Sprintf("Error creating index %s: %s", IndexItems, res.Status()))
	}

	logger.Info(fmt.Sprintf("Index '%s' created successfully with custom mapping", IndexItems))
}

// ---------------------------------------------------------
// 4. МЕТОДИ РОБОТИ З ДАНИМИ (INDEX)
// ---------------------------------------------------------

func (c *esClient) Index(index string, id string, sellerID string, doc any) rest_errors.RestErr {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	bytesData, err := json.Marshal(doc)
	if err != nil {
		logger.Error("error marshaling json", err)
		return rest_errors.NewInternalServerError("error trying to marshal document", err)
	}

	res, err := c.client.Index(
		index,
		strings.NewReader(string(bytesData)),
		c.client.Index.WithDocumentID(id),
		c.client.Index.WithRouting(sellerID),
		c.client.Index.WithContext(ctx),
	)

	if err != nil {
		logger.Error("error connecting to elasticsearch", err)
		return rest_errors.NewInternalServerError("database connection error", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error(fmt.Sprintf("elasticsearch returned error: %s", res.Status()), nil)
		return rest_errors.NewInternalServerError("error indexing document", nil)
	}

	return nil
}
