package main

import (
	"github.com/SerhiiKhyzhko/bookstore_items-api/app"
	"github.com/SerhiiKhyzhko/bookstore_items-api/clients/elasticsearch"
)

func main() {
	elasticsearch.Client.Init()
	app.StartApp()
}
