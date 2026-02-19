package main

import (
	"github.com/SerhiiKhyzhko/bookstore-oauth-go/oauth"
	"github.com/SerhiiKhyzhko/bookstore_items-api/app"
	"github.com/SerhiiKhyzhko/bookstore_items-api/clients/elasticsearch"
	"github.com/SerhiiKhyzhko/bookstore_items-api/config"
	"github.com/SerhiiKhyzhko/bookstore_items-api/controllers"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/items"
	"github.com/SerhiiKhyzhko/bookstore_items-api/internal/elsticsearch_client"
	"github.com/SerhiiKhyzhko/bookstore_items-api/services"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load();err != nil {
		logger.Info("Error loading .env file")
	}

	config.Init()
	oauth.Init(config.RestyBaseUrl)
	esClient, err := elsticsearch_client.NewElasticClient(config.EsHosts)
	if err != nil{
		logger.Fatal("CRITICAL: Failed to connect to Elasticsearch: ", err)
	}
	if err = elsticsearch_client.EnsureIndexCreated(esClient); err != nil {
		logger.Fatal("CRITICAL: Failed to check/create index: ", err)
	}
	elasticsearch := elasticsearch.NewEsClient(esClient)
	dao := items.NewItemDao(elasticsearch)
	service := services.NewItemsService(dao)
	controller := controllers.NewItemsController(service)
	app.StartApp(controller)
}
