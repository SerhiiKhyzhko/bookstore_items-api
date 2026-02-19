package elsticsearch_client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/elastic/go-elasticsearch/v9"
)

const (
	IndexItems = "items"
	// shard settings
	itemMapping = `{
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
)

func NewElasticClient(addreses string) (*elasticsearch.TypedClient, error) {
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
	var client *elasticsearch.TypedClient
	client, err = elasticsearch.NewTypedClient(cfg)
	if err != nil {
		logger.Error("Error creating elasticsearch typed client", err)
		return nil, err
	}

	res, err := client.Info().Do(context.Background())
	if err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("Elasticsearch client connected. Claster: %s, Version %s", res.ClusterName, res.Version)
	logger.Info(msg)

	return client, nil
}

func EnsureIndexCreated(client *elasticsearch.TypedClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.Indices.Exists(IndexItems).Do(ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("error when check the index %s", IndexItems), err)
		return err
	}

	if exists {
		logger.Info(fmt.Sprintf("Indedx %s already exists", IndexItems))
		return nil
	}

	resp, err := client.Indices.
		Create(IndexItems).
		Raw(strings.NewReader(itemMapping)).
		Do(ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("failed to create index %s", IndexItems), err)
		return err
	}

	if resp.Acknowledged {
		logger.Info(fmt.Sprintf("index %s created successfully with settings and mappings", IndexItems))
	}
	return nil
}
