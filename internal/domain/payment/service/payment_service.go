package service

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"orders-center/internal/domain/payment/entity"
	"orders-center/internal/domain/payment/repository"
	"orders-center/internal/utils"
)

type PaymentService interface {
	GetByID(ctx context.Context, id uuid.UUID) (entity.OrderPayment, error)
	Create(ctx context.Context, payment entity.OrderPayment) (entity.OrderPayment, error)
	GetPaymentsByOrderID(ctx context.Context, id uuid.UUID) ([]entity.OrderPayment, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

/*func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}*/

func NewPaymentService(repo repository.PaymentRepository) PaymentService {

	return &paymentService{repo: repo}
}

func (s *paymentService) GetPaymentsByOrderID(ctx context.Context, id uuid.UUID) ([]entity.OrderPayment, error) {
	return s.repo.GetPaymentsByOrderID(ctx, id)
}

// Получение платежа по ID
func (s *paymentService) GetByID(ctx context.Context, id uuid.UUID) (entity.OrderPayment, error) {
	return s.repo.GetPayment(ctx, id)
}

// Создание нового платежа
func (s *paymentService) Create(ctx context.Context, payment entity.OrderPayment) (entity.OrderPayment, error) {
	// Сериализуем CreditData и CardPaymentData в json.RawMessage
	creditDataRaw, err := json.Marshal(payment.CreditData)
	if err != nil {
		return entity.OrderPayment{}, err
	}

	cardPaymentDataRaw, err := json.Marshal(payment.CardPaymentData)
	if err != nil {
		return entity.OrderPayment{}, err
	}

	// Преобразуем поля с учетом типов данных
	arg := repository.CreatePaymentParams{
		OrderID:        pgtype.UUID{Bytes: [16]byte(payment.OrderID.Bytes()), Valid: true}, // Преобразуем OrderID в pgtype.UUID
		Type:           string(payment.Type),                                               // Преобразуем в строку
		Sum:            utils.ToNumeric(payment.Sum),                                       // Преобразуем Sum в pgtype.Numeric
		Payed:          utils.ToBool(payment.Payed),                                        // Преобразуем в pgtype.Bool
		Info:           utils.ToText(payment.Info),                                         // Преобразуем Info в pgtype.Text
		ContractNumber: utils.ToText(payment.ContractNumber),                               // Преобразуем ContractNumber в pgtype.Text
		ExternalID:     utils.ToText(payment.ExternalID),                                   // Преобразуем ExternalID в pgtype.Text
		CreditData:     creditDataRaw,                                                      // Сериализованные данные CreditData
		CardData:       cardPaymentDataRaw,                                                 // Сериализованные данные CardPaymentData
	}

	return s.repo.CreatePayment(ctx, arg)
}
