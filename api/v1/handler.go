package v1

import (
	"github.com/ToDoApp/api/models"
	"github.com/ToDoApp/config"
	"github.com/ToDoApp/storage"
)

type handlerV1 struct {
	cfg      *config.Config
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
}

type HandlerV1Options struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		cfg:      options.Cfg,
		storage:  options.Storage,
		inMemory: options.InMemory,
	}
}

func errorResponse(err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Error: err.Error(),
	}
}
