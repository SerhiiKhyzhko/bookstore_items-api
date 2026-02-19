package app

import (
	"github.com/SerhiiKhyzhko/bookstore_items-api/controllers"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine, itemsCtrl *controllers.ItemsController) {
	router.POST("/items", itemsCtrl.Create)
	router.GET("/items/:id", itemsCtrl.Get)
	router.POST("/items/search", itemsCtrl.Search)
	router.DELETE("/items/:id", itemsCtrl.Delete)
	router.PATCH("/items/:id", itemsCtrl.Patch)
	router.PUT("/items/:id", itemsCtrl.Put)
}