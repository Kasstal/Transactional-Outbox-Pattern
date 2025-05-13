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

	"time"
)

func randomString(prefix string) string {
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
		Type:        randomString("order_type_"),
		Status:      randomString("status_"),
		City:        randomString("city_"),
		Subdivision: randomString("subdivision_"),
		Price:       randomFloat(1000.0, 10000.0),
		Platform:    randomString("platform_"),
		GeneralID:   randomUUID(),
		OrderNumber: randomString("ORD-"),
		Executor:    randomString("executor_"),
		CreatedAt:   randomDate(),
		UpdatedAt:   randomDate(),
	}

	// Generate random order items
	numItems := rand.Intn(5) + 1
	var items []orderItem.OrderItem

	for i := 0; i < numItems; i++ {
		item := orderItem.OrderItem{
			ID:            int32(i),
			ProductID:     randomString("prod_"),
			ExternalID:    randomString("ext_"),
			Status:        randomString("status_"),
			BasePrice:     randomFloat(100.0, 1000.0),
			Price:         randomFloat(100.0, 1000.0),
			EarnedBonuses: randomFloat(0.0, 500.0),
			SpentBonuses:  randomFloat(0.0, 500.0),
			Gift:          randomBool(),
			OwnerID:       randomString("user_"),
			DeliveryID:    randomString("delivery_"),
			ShopAssistant: randomString("assistant_"),
			Warehouse:     randomString("warehouse_"),
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
			Info:            randomString("payment_info_"),
			CreditData:      nil,
			ContractNumber:  randomString("contract_"),
			CardPaymentData: nil,
			ExternalID:      randomString("ext_"),
		}
		payments = append(payments, payment)
	}

	// Generate random history
	numHistory := rand.Intn(5) + 1
	var history []ordHistory.History
	for i := 0; i < numHistory; i++ {
		hist := ordHistory.History{
			Type:     randomString("type_"),
			TypeId:   rand.Int31(),
			OldValue: []byte(fmt.Sprintf(`{"old_key":"%s"}`, randomString("old_value_"))),
			Value:    []byte(fmt.Sprintf(`{"new_key":"%s"}`, randomString("new_value_"))),
			Date:     randomDate(),
			UserID:   randomString("user_"),
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
