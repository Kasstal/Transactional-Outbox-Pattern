package utils

import (
	"fmt"
	"github.com/gofrs/uuid"
	"math/rand"
	ordHistory "orders-center/internal/domain/history/entity"
	order "orders-center/internal/domain/order/entity"
	orderItem "orders-center/internal/domain/order_item/entity"
	payment "orders-center/internal/domain/payment/entity"
	orderFull "orders-center/internal/service/order_full/entity"
	"strings"

	"time"
)

func randomStringNumeric(prefix string) string {
	return fmt.Sprintf("%s%d", prefix, rand.Intn(1000))
}

func randomUUID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomDate() time.Time {
	return time.Now().Add(time.Duration(rand.Intn(1000)) * time.Hour * 24)
}

func randomBool() bool {
	return rand.Intn(2) == 1
}
func randomPaymentType() payment.PaymentType {
	paymentTypes := []payment.PaymentType{
		payment.PaymentTypeCashAtShop,
		payment.PaymentTypeCashToCourier,
		payment.PaymentTypeCard,
		payment.PaymentTypeCardOnline,
		payment.PaymentTypeCredit,
		payment.PaymentTypeBonuses,
		payment.PaymentTypeCashless,
		payment.PaymentTypePrepayment,
	}

	// Return a random payment type from the slice
	return paymentTypes[rand.Intn(len(paymentTypes))]
}

func RandomOrderFull() orderFull.OrderFull {
	order := order.Order{
		ID:          randomUUID(),
		Type:        randomStringNumeric("order_type_"),
		Status:      randomStringNumeric("status_"),
		City:        randomStringNumeric("city_"),
		Subdivision: randomStringNumeric("subdivision_"),
		Price:       randomFloat(1000.0, 10000.0),
		Platform:    randomStringNumeric("platform_"),
		GeneralID:   randomUUID(),
		OrderNumber: randomStringNumeric("ORD-"),
		Executor:    randomStringNumeric("executor_"),
		CreatedAt:   randomDate(),
		UpdatedAt:   randomDate(),
	}

	// Generate random order items
	numItems := rand.Intn(5) + 1
	var items []orderItem.OrderItem

	for i := 0; i < numItems; i++ {
		item := orderItem.OrderItem{
			ID:            int32(i),
			ProductID:     randomStringNumeric("prod_"),
			ExternalID:    randomStringNumeric("ext_"),
			Status:        randomStringNumeric("status_"),
			BasePrice:     randomFloat(100.0, 1000.0),
			Price:         randomFloat(100.0, 1000.0),
			EarnedBonuses: randomFloat(0.0, 500.0),
			SpentBonuses:  randomFloat(0.0, 500.0),
			Gift:          randomBool(),
			OwnerID:       randomStringNumeric("user_"),
			DeliveryID:    randomStringNumeric("delivery_"),
			ShopAssistant: randomStringNumeric("assistant_"),
			Warehouse:     randomStringNumeric("warehouse_"),
			OrderId:       order.ID,
		}
		items = append(items, item)
	}

	// Generate random payments
	numPayments := rand.Intn(3) + 1
	var payments []payment.OrderPayment
	for i := 0; i < numPayments; i++ {
		payment := payment.OrderPayment{
			ID:              randomUUID(),
			OrderID:         order.ID,
			Type:            randomPaymentType(),
			Sum:             randomFloat(100.0, 1000.0),
			Payed:           randomBool(),
			Info:            randomStringNumeric("payment_info_"),
			CreditData:      randomCreditData(),
			ContractNumber:  randomStringNumeric("contract_"),
			CardPaymentData: randomCardPaymentData(),
			ExternalID:      randomStringNumeric("ext_"),
		}
		payments = append(payments, payment)
	}

	// Generate random history
	numHistory := rand.Intn(5) + 1
	var history []ordHistory.History
	for i := 0; i < numHistory; i++ {
		hist := ordHistory.History{
			Type:     randomStringNumeric("type_"),
			TypeId:   rand.Int31(),
			OldValue: []byte(fmt.Sprintf(`{"old_key":"%s"}`, randomStringNumeric("old_value_"))),
			Value:    []byte(fmt.Sprintf(`{"new_key":"%s"}`, randomStringNumeric("new_value_"))),
			Date:     randomDate(),
			UserID:   randomStringNumeric("user_"),
			OrderID:  order.ID,
		}
		history = append(history, hist)
	}

	return orderFull.OrderFull{
		Order:    order,
		Items:    items,
		Payments: payments,
		History:  history,
	}

}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteRune(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}

func randomDigits(n int) string {
	digits := []rune("0123456789")
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteRune(digits[rand.Intn(len(digits))])
	}
	return sb.String()
}

func randomCreditType() string {
	types := []string{"Mortgage", "Personal", "Auto", "Business"}
	return types[rand.Intn(len(types))]
}

func randomBank() string {
	banks := []string{"Bank of America", "Chase", "Wells Fargo", "Citibank", "HSBC"}
	return banks[rand.Intn(len(banks))]
}

func randomProvider() string {
	providers := []string{"Visa", "Mastercard", "American Express", "Discover"}
	return providers[rand.Intn(len(providers))]
}

func randomTransactionID() string {
	// e.g. 16 alphanumeric chars
	return randomString(8) + randomDigits(8)
}

func randomIIN() string {
	// IIN usually 12 digits (e.g. national ID)
	return randomDigits(12)
}

func randomCreditData() *payment.CreditData {
	return &payment.CreditData{
		Bank:           randomBank(),
		Type:           randomCreditType(),
		NumberOfMonths: int16(rand.Intn(60) + 1),                         // 1 to 60 months
		PaySumPerMonth: float64(int(randomFloat(1000, 50000)*100)) / 100, // Rounded to 2 decimals
		BrokerID:       int32(rand.Intn(1000) + 1),                       // Broker ID 1-1000
		IIN:            randomIIN(),
	}
}

func randomCardPaymentData() *payment.CardPaymentData {
	return &payment.CardPaymentData{
		Provider:      randomProvider(),
		TransactionId: randomTransactionID(),
	}
}
