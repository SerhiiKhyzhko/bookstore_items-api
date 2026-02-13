package controllers

import (
	"net/http"
	"strings"

	"github.com/SerhiiKhyzhko/bookstore-oauth-go/oauth"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/items"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/queries"
	"github.com/SerhiiKhyzhko/bookstore_items-api/services"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

var ItemsController itemsControllerInterface = &itemsController{}

type itemsControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Search(c *gin.Context)
	Delete(c *gin.Context)
	Put(c *gin.Context)	
	Patch(c *gin.Context)
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
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (i *itemsController) Get(c *gin.Context) {
	itemId := strings.TrimSpace(c.Param("id"))
	item, err := services.ItemsService.Get(itemId)
	if err != nil {
		c.JSON(err.Status(), err.Message())
		return
	}
	c.JSON(http.StatusOK, item)
}

func (i *itemsController) Search(c *gin.Context) {
	var query queries.EsQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid query json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	items, searchErr := services.ItemsService.Search(query)
	if searchErr != nil {
		c.JSON(searchErr.Status(), searchErr.Message())
		return
	}
	c.JSON(http.StatusOK, items)
}

func (i *itemsController) Delete(c *gin.Context) {
	itemId := strings.TrimSpace(c.Param("id"))
	if deleteErr := services.ItemsService.Delete(itemId); deleteErr != nil {
		c.JSON(deleteErr.Status(), deleteErr)
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (i *itemsController) Put(c *gin.Context) {
	itemId := strings.TrimSpace(c.Param("id"))
	
	var itemRequest items.Item
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid item json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	itemRequest.Id = itemId

	result, err := services.ItemsService.Put(itemRequest)
	if err != nil {
		c.JSON(err.Status(), err.Message())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (i *itemsController) Patch(c *gin.Context) {
	itemId := strings.TrimSpace(c.Param("id"))

	var itemRequest items.PartialUpdateItem
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid update item json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, err := services.ItemsService.Patch(itemRequest, itemId)
	if err != nil {
		c.JSON(err.Status(), err.Message())
		return
	}

	c.JSON(http.StatusOK, result)
}