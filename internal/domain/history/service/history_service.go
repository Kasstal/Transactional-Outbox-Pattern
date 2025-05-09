package service

import (
	"context"
	"github.com/gofrs/uuid"
	"orders-center/internal/domain/history/entity"
	"orders-center/internal/domain/history/repository"
)

type HistoryService interface {
	GetByID(ctx context.Context, id int32) (entity.History, error)
	Create(ctx context.Context, history entity.History) (entity.History, error)
	GetHistoriesByOrderId(ctx context.Context, orderId uuid.UUID) ([]entity.History, error)
}

type historyService struct {
	repo repository.HistoryRepository
}

func (s *historyService) GetHistoriesByOrderId(ctx context.Context, orderId uuid.UUID) ([]entity.History, error) {
	return s.repo.GetHistoriesByOrderID(ctx, orderId)
}

func NewHistoryService(repo repository.HistoryRepository) HistoryService {
	return &historyService{repo: repo}
}

// Получение истории по ID
func (s *historyService) GetByID(ctx context.Context, id int32) (entity.History, error) {
	return s.repo.GetHistory(ctx, id)
}

// Создание новой записи в истории
func (s *historyService) Create(ctx context.Context, history entity.History) (entity.History, error) {
	// Преобразуем поля с учетом типов данных
	arg := repository.CreateHistoryParams{
		Type:     history.Type,
		TypeID:   history.TypeId,
		OldValue: history.OldValue, // предполагается, что OldValue и Value уже в формате json.RawMessage
		Value:    history.Value,    // предполагается, что OldValue и Value уже в формате json.RawMessage
		UserID:   history.UserID,
		OrderID:  history.OrderID, // Преобразуем UUID в pgtype.UUID
	}

	return s.repo.CreateHistory(ctx, arg)
}
