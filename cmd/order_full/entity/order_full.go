package entity

import (
	history "orders-center/internal/domain/history/entity"
	order "orders-center/internal/domain/order/entity"
	orderItem "orders-center/internal/domain/order_item/entity"
	payment "orders-center/internal/domain/payment/entity"
)

type OrderFull struct {
	Order    order.Order            `json:"order"`
	Items    []orderItem.OrderItem  `json:"items"`
	Payments []payment.OrderPayment `json:"payments"`
	History  []history.History      `json:"history"`
}
