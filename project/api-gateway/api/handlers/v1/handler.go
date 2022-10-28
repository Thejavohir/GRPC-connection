package v1

import (
	"github.com/project/api-gateway/config"
	"github.com/project/api-gateway/pkg/logger"
	"github.com/project/api-gateway/services"
	"github.com/project/api-gateway/storage/repo"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redisStorage   repo.InMemoryStorageI
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.InMemoryStorageI
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redisStorage:          c.Redis,
	}
}
