package entity

import "github.com/gofrs/uuid"

type OrderItem struct {
	ID            int32     `json:"id"`
	ProductID     string    `json:"product_id"`
	ExternalID    string    `json:"external_id"`
	Status        string    `json:"status"`
	BasePrice     float64   `json:"base_price"`
	Price         float64   `json:"price"`
	EarnedBonuses float64   `json:"earned_bonuses"`
	SpentBonuses  float64   `json:"spent_bonuses"`
	Gift          bool      `json:"gift"`
	OwnerID       string    `json:"owner_id"`
	DeliveryID    string    `json:"delivery_id"`
	ShopAssistant string    `json:"shop_assistant"`
	Warehouse     string    `json:"warehouse"`
	OrderId       uuid.UUID `json:"order_id"`
}
