package app

import (
	"github.com/SerhiiKhyzhko/bookstore_items-api/controllers"
	"github.com/gin-gonic/gin"
)

func StartApp(itemsCtrl *controllers.ItemsController) {
	router := gin.Default()
	mapUrls(router, itemsCtrl)
	router.Run(":8000")
}