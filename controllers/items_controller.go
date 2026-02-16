package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SerhiiKhyzhko/bookstore-oauth-go/oauth"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/items"
	"github.com/SerhiiKhyzhko/bookstore_items-api/domain/queries"
	"github.com/SerhiiKhyzhko/bookstore_items-api/item_errors"
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

func requestError(reqErr error) rest_errors.RestErr {
	switch {
	case errors.Is(reqErr, item_errors.RequestTimeoutErr):
		return rest_errors.NewRestError("request timeout", http.StatusRequestTimeout, "database error", nil)
	case errors.Is(reqErr, item_errors.NotFoundErr):
		return rest_errors.NewNotFoundError("item not found with given id")
	case errors.Is(reqErr, item_errors.ParseErr):
		return rest_errors.NewInternalServerError("error when trying to parse response", errors.New("database error"))
	default :
		return rest_errors.NewInternalServerError("internal server error", errors.New("database error"))
	}

}

func (i *itemsController) Create(c *gin.Context) {
	ctx := c.Request.Context()
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
	result, err := services.ItemsService.Create(ctx, itemRequest)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr.Message())
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (i *itemsController) Get(c *gin.Context) {
	ctx := c.Request.Context()
	itemId := strings.TrimSpace(c.Param("id"))
	item, err := services.ItemsService.Get(ctx, itemId)
	if err != nil {
		restErr := requestError(err)
		if errors.Is(err, item_errors.NotFoundErr) {
			restErr = rest_errors.NewNotFoundError(fmt.Sprintf("item not found with given id %s", itemId))
		}
		c.JSON(restErr.Status(), restErr.Message())
		return
	}
	c.JSON(http.StatusOK, item)
}

func (i *itemsController) Search(c *gin.Context) {
	ctx := c.Request.Context()
	var query queries.EsQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid query json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	items, searchErr := services.ItemsService.Search(ctx, query)
	if searchErr != nil {
		restErr := requestError(searchErr)
		c.JSON(restErr.Status(), restErr.Message())
		return
	}
	c.JSON(http.StatusOK, items)
}

func (i *itemsController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	itemId := strings.TrimSpace(c.Param("id"))
	if deleteErr := services.ItemsService.Delete(ctx, itemId); deleteErr != nil {
		restErr := requestError(deleteErr)
		if errors.Is(deleteErr, item_errors.NotFoundErr) {
			restErr = rest_errors.NewNotFoundError(fmt.Sprintf("item not found with given id %s", itemId))
		}
		c.JSON(restErr.Status(), deleteErr)
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (i *itemsController) Put(c *gin.Context) {
	ctx := c.Request.Context()
	itemId := strings.TrimSpace(c.Param("id"))

	var itemRequest items.Item
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid update entire item json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	itemRequest.Id = itemId

	result, err := services.ItemsService.Put(ctx, itemRequest)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr.Message())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (i *itemsController) Patch(c *gin.Context) {
	ctx := c.Request.Context()
	itemId := strings.TrimSpace(c.Param("id"))

	var itemRequest items.PartialUpdateItem
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid update item json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, err := services.ItemsService.Patch(ctx, itemRequest, itemId)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr.Message())
		return
	}

	c.JSON(http.StatusOK, result)
}
