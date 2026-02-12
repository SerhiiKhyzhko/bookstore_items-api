package app

import "github.com/SerhiiKhyzhko/bookstore_items-api/controllers"

func mapUrls() {
	router.POST("/items", controllers.ItemsController.Create)
	router.GET("/items/:id", controllers.ItemsController.Get)
	router.POST("/items/search", controllers.ItemsController.Search)
	router.DELETE("/items/:id", controllers.ItemsController.Delete)
}