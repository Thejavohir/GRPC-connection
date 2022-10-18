package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/project/api-gateway/api/handlers/models"
	pbp "github.com/project/api-gateway/genproto/product_service"
	l "github.com/project/api-gateway/pkg/logger"
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
		Name: productBody.Name,
		Model: productBody.Model,
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
