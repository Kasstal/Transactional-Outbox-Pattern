package entity

import "github.com/gofrs/uuid"

type OrderPayment struct {
	ID              uuid.UUID
	OrderID         uuid.UUID
	Type            PaymentType
	Sum             float64
	Payed           bool
	Info            string
	CreditData      *CreditData
	ContractNumber  string
	CardPaymentData *CardPaymentData
	ExternalID      string
}

type PaymentType string

const (
	PaymentTypeCashAtShop    PaymentType = "cash_at_shop"
	PaymentTypeCashToCourier PaymentType = "cash_to_courier"
	PaymentTypeCard          PaymentType = "card"
	PaymentTypeCardOnline    PaymentType = "card_online"
	PaymentTypeCredit        PaymentType = "credit"
	PaymentTypeBonuses       PaymentType = "bonuses"
	PaymentTypeCashless      PaymentType = "cashless"
	PaymentTypePrepayment    PaymentType = "prepayment"
)

type CreditData struct {
	Bank           string
	Type           string
	NumberOfMonths int16
	PaySumPerMonth float64
	BrokerID       int32
	IIN            string
}

type CardPaymentData struct {
	Provider      string
	TransactionId string
}
