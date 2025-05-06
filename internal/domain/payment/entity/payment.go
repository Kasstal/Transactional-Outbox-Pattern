package entity

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
)

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
	Bank           string  `json:"bank"`
	Type           string  `json:"type"`
	NumberOfMonths int16   `json:"number_of_months"`
	PaySumPerMonth float64 `json:"pay_sum_per_month"`
	BrokerID       int32   `json:"broker_id"`
	IIN            string  `json:"iin"`
}

type CardPaymentData struct {
	Provider      string `json:"provider"`
	TransactionId string `json:"transaction_id"`
}

func GetCreditData(raw json.RawMessage) (*CreditData, error) {
	creditData := &CreditData{}
	err := json.Unmarshal(raw, creditData)
	if err != nil {
		return nil, err
	}
	return creditData, nil
}

func GetCardPaymentData(raw json.RawMessage) (*CardPaymentData, error) {
	cardPaymentData := &CardPaymentData{}
	err := json.Unmarshal(raw, cardPaymentData)
	if err != nil {
		return nil, err
	}
	return cardPaymentData, nil
}

func GetPaymentType(str string) (PaymentType, error) {
	switch str {
	case string(PaymentTypeCashAtShop),
		string(PaymentTypeCashToCourier),
		string(PaymentTypeCard),
		string(PaymentTypeCardOnline),
		string(PaymentTypeCredit),
		string(PaymentTypeBonuses),
		string(PaymentTypeCashless),
		string(PaymentTypePrepayment):
		return PaymentType(str), nil
	default:
		return "", fmt.Errorf("invalid payment type: %s", str)
	}
}
