package main

import (
	"github.com/SerhiiKhyzhko/bookstore-oauth-go/oauth"
	"github.com/SerhiiKhyzhko/bookstore_items-api/app"
	"github.com/SerhiiKhyzhko/bookstore_items-api/clients/elasticsearch"
	"github.com/SerhiiKhyzhko/bookstore_items-api/config"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load();err != nil {
		logger.Info("Error loading .env file")
	}

	config.Init()
	oauth.Init(config.RestyBaseUrl)
	elasticsearch.Client.Init(config.EsHosts)
	app.StartApp()
}
