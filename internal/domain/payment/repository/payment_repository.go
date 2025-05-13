package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/payment/entity"
	transactional "orders-center/internal/service/transactional"
	"orders-center/internal/utils"
)

type paymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) PaymentRepository {
	return &paymentRepository{pool: pool}
}

func (r *paymentRepository) GetPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) ([]entity.OrderPayment, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	payments, err := query.GetPaymentsByOrderID(ctx, utils.ToUUID(orderID))
	if err != nil {
		return nil, err
	}
	paymentsEntity := make([]entity.OrderPayment, len(payments))
	for _, payment := range payments {
		if err != nil {
			return []entity.OrderPayment{}, err
		}
		sum, err := payment.Sum.Float64Value()
		if err != nil {
			return []entity.OrderPayment{}, err
		}
		creditData, err := entity.GetCreditData(payment.CreditData)
		if err != nil {
			return []entity.OrderPayment{}, err
		}

		cardPaymentData, err := entity.GetCardPaymentData(payment.CardData)
		if err != nil {
			return []entity.OrderPayment{}, err
		}
		paymentType, err := entity.GetPaymentType(payment.Type)
		if err != nil {
			return []entity.OrderPayment{}, err
		}
		paymentEntity := entity.OrderPayment{
			ID:              payment.ID.Bytes,
			OrderID:         payment.OrderID.Bytes,
			Type:            paymentType,
			Sum:             sum.Float64,
			Payed:           payment.Payed.Bool,
			Info:            payment.Info.String,
			ContractNumber:  payment.ContractNumber.String,
			ExternalID:      payment.ExternalID.String,
			CreditData:      creditData,
			CardPaymentData: cardPaymentData,
		}
		paymentsEntity = append(paymentsEntity, paymentEntity)
	}

	return paymentsEntity, nil
}

func (r *paymentRepository) CreatePayment(ctx context.Context, arg CreatePaymentParams) (entity.OrderPayment, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	sqlcArg := db.CreatePaymentParams{
		OrderID:        arg.OrderID,
		Type:           arg.Type,
		Sum:            arg.Sum,
		Payed:          arg.Payed,
		Info:           arg.Info,
		ContractNumber: arg.ContractNumber,
		ExternalID:     arg.ExternalID,
		CreditData:     arg.CreditData,
		CardData:       arg.CardData,
	}

	payment, err := query.CreatePayment(ctx, sqlcArg)
	if err != nil {
		return entity.OrderPayment{}, err
	}
	sum, err := payment.Sum.Float64Value()
	if err != nil {
		return entity.OrderPayment{}, err
	}
	creditData, err := entity.GetCreditData(payment.CreditData)
	if err != nil {
		return entity.OrderPayment{}, err
	}

	cardPaymentData, err := entity.GetCardPaymentData(payment.CardData)
	if err != nil {
		return entity.OrderPayment{}, err
	}
	paymentType, err := entity.GetPaymentType(payment.Type)
	if err != nil {
		return entity.OrderPayment{}, err
	}
	paymentEntity := entity.OrderPayment{
		ID:              payment.ID.Bytes,
		OrderID:         payment.OrderID.Bytes,
		Type:            paymentType,
		Sum:             sum.Float64,
		Payed:           payment.Payed.Bool,
		Info:            payment.Info.String,
		ContractNumber:  payment.ContractNumber.String,
		ExternalID:      payment.ExternalID.String,
		CreditData:      creditData,
		CardPaymentData: cardPaymentData,
	}
	return paymentEntity, nil
}
func (r *paymentRepository) GetPayment(ctx context.Context, id uuid.UUID) (entity.OrderPayment, error) {

	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}
	payment, err := query.GetPayment(ctx, pgtype.UUID{Bytes: id})
	if err != nil {
		return entity.OrderPayment{}, err
	}
	sum, err := payment.Sum.Float64Value()
	if err != nil {
		return entity.OrderPayment{}, err
	}
	creditData, err := entity.GetCreditData(payment.CreditData)
	if err != nil {
		return entity.OrderPayment{}, err
	}

	cardPaymentData, err := entity.GetCardPaymentData(payment.CardData)
	if err != nil {
		return entity.OrderPayment{}, err
	}
	paymentType, err := entity.GetPaymentType(payment.Type)
	if err != nil {
		return entity.OrderPayment{}, err
	}
	paymentEntity := entity.OrderPayment{
		ID:              payment.ID.Bytes,
		OrderID:         payment.OrderID.Bytes,
		Type:            paymentType,
		Sum:             sum.Float64,
		Payed:           payment.Payed.Bool,
		Info:            payment.Info.String,
		ContractNumber:  payment.ContractNumber.String,
		ExternalID:      payment.ExternalID.String,
		CreditData:      creditData,
		CardPaymentData: cardPaymentData,
	}
	return paymentEntity, nil
}
func (r *paymentRepository) DeletePayment(ctx context.Context, id uuid.UUID) error {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	return query.DeletePayment(ctx, pgtype.UUID{Bytes: id})
}
