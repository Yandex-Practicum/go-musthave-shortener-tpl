package workers

import (
	"context"
	"fmt"

	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
)

// Worker - интерфейс воркера.
//
//go:generate mockgen -source=worker.go -destination=mock_worker.go -package=workers
type Worker interface {
	SendDeletionRequestToWorker(req DeletionRequest) error
}

// DeletionRequest - запрос на удаление URL из хранилища.
type DeletionRequest struct {
	User string
	URLs []string
}

// WorkerDeleted - воркер для удаления URL из хранилища.
var deleteQueue = make(chan DeletionRequest, 10)

// WorkerDeleted - воркер для удаления URL из хранилища.
type WorkerDeleted struct {
	storage      *service.Service
	errorChannel chan error
}

// NewWorkerDeleted - конструктор воркера.
func NewWorkerDeleted(storage *service.Service) *WorkerDeleted {
	return &WorkerDeleted{
		storage: storage,
	}
}

// StartWorkerDeletion стартует воркер для удаления URL из хранилища.
func (w *WorkerDeleted) StartWorkerDeletion(ctx context.Context) {
	for {
		select {
		case req := <-deleteQueue:
			go w.processDeletion(ctx, req)
		case <-ctx.Done():
			return
		}
	}
}

// processDeletion обрабатывает удаление URL из хранилища.
func (w *WorkerDeleted) processDeletion(ctx context.Context, req DeletionRequest) {
	if err := w.storage.DeletedURLs(req.URLs, req.User); err != nil {
		select {
		case w.errorChannel <- err:
		case <-ctx.Done():
			fmt.Println("Operation canceled, skipping error reporting.")
		}
	}
}

// StartErrorListener обрабатывает ошибки воркера.
func (w *WorkerDeleted) StartErrorListener(ctx context.Context) {
	for {
		select {
		case err := <-w.errorChannel:
			fmt.Printf("Error processing deletion request: %v\n", err)
		case <-ctx.Done():
			fmt.Println("Error listener shutting down due to context cancellation.")
			return
		}
	}
}

// SendDeletionRequestToWorker отправляет запрос на удаление URL из хранилища.
func (w *WorkerDeleted) SendDeletionRequestToWorker(req DeletionRequest) error {
	select {
	case deleteQueue <- req:
		return nil
	default:
		return fmt.Errorf("the deletion request queue is currently full, please try again later")
	}
}
