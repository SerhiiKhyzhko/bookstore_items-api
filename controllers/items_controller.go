package controllers

import (
	"net/http"

	"github.com/SerhiiKhyzhko/bookstore-oauth-go/oauth"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/items"
	"github.com/SerhiiKhyzhko/bookstore_items-api/services"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

var ItemsController itemsControllerInterface = &itemsController{}

type itemsControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
}
type itemsController struct{}

func (i *itemsController) Create(c *gin.Context) {
	if err := oauth.AutenticationRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}

	var itemRequest items.Item
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid item json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	itemRequest.Seller = oauth.GetClientId(c.Request)
	result, err := services.ItemsService.Create(itemRequest)
	if err != nil {
		c.JSON(err.Status(), err.Message())
	}

	c.JSON(http.StatusCreated, result)
}

func (i *itemsController) Get(c *gin.Context) {

}
