package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/project/api-gateway/api/handlers/v1"
	"github.com/project/api-gateway/config"
	"github.com/project/api-gateway/pkg/logger"
	"github.com/project/api-gateway/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/project/api-gateway/api/docs"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
}

// New ...
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
	})

	api := router.Group("/v1")
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/users/:id", handlerV1.GetUser)

	api.POST("/product", handlerV1.CreateProduct)

	// url := ginSwagger.URL("swagger/docs.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
