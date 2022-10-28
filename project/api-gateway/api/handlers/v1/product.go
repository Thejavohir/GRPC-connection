package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/project/api-gateway/api/handlers/models"
	pbp "github.com/project/api-gateway/genproto/product_service"
	l "github.com/project/api-gateway/pkg/logger"
	"github.com/project/api-gateway/pkg/utils"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateProduct creates product
// @Summary creates product api
// @Description new product creation
// @Tags Product
// @Accept json
// @Produce json
// @Param request body models.Product true "product"
// @Success 200 "success"
// @Router /v1/product [post]
func (h *handlerV1) CreateProduct(c *gin.Context) {
	var (
		productBody models.Product
	)

	err := c.ShouldBindJSON(&productBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error while binding": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	product, err := h.serviceManager.ProductService().CreateProduct(ctx, &pbp.Product{
		Name:    productBody.Name,
		Model:   productBody.Model,
		OwnerId: productBody.OwnerID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error create product": err.Error(),
		})
		h.log.Error("error creating product", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct gets product by id
// @Summary getting product by ID
// @Description product by ID
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 "success"
// @Router /v1/product/{id} [get]
func (h *handlerV1) GetProduct(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	id, err := strconv.ParseInt(guid, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error parsing": err.Error(),
		})
		h.log.Error("failed to parse string to int", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	product, err := h.serviceManager.ProductService().GetProduct(ctx, &pbp.GetProductRequest{
		Id: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed getting product", l.Error(err))
	}
	c.JSON(http.StatusOK, product)
}

// GetUser gets user products by id
// @Summary getting user products by ID
// @Description user products by ID
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 "success"
// @Router /v1/products/{id} [get]
func (h *handlerV1) GetUserProducts(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true
	guid := c.Param("id")

	id, err := strconv.ParseInt(guid, 10, 64)
	fmt.Println(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error parsing to int", l.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	userProducts, err := h.serviceManager.ProductService().GetUserProducts(ctx, &pbp.GetUserProductsRequest{
		OwnerId: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error getting getting user products", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, userProducts)
}

// ListProducts gets all products
// @Summary getting all products
// @Description all products
// @Tags Product
// @Accept json
// @Produce json
// @Param page query string false "query params"
// @Param limit query string false "query params"
// @Success 200 {object} models.ListProducts
// @Failure 400 {object} models.Error
// @Router /v1/products/all [get]
func (h *handlerV1) ListProducts(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	params, err := utils.ParseQueryParams(queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":err[0],
		})
		h.log.Error("failed to parse query params to JSON"+ err[0])
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, errS := h.serviceManager.ProductService().ListProducts(ctx, &pbp.LPreq{Page: params.Page, Limit: params.Limit})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errS.Error(),
		})
		h.log.Error("failed to list users", l.Error(errS))
		return
	}
	c.JSON(http.StatusOK, response)
}

