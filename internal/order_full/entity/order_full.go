package entity

import (
	history "orders-center/internal/domain/history/entity"
	order "orders-center/internal/domain/order/entity"
	orderItem "orders-center/internal/domain/order_item/entity"
	payment "orders-center/internal/domain/payment/entity"
)

type OrderFull struct {
	Order    order.Order
	Items    []orderItem.OrderItem
	Payments []payment.OrderPayment
	History  history.History
}
