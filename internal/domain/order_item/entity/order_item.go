package entity

import "github.com/google/uuid"

type OrderItem struct {
	ProductID     string
	ExternalID    string
	Status        string
	BasePrice     float64
	Price         float64
	EarnedBonuses float64
	SpentBonuses  float64
	Gift          bool
	OwnerID       string
	DeliveryID    string
	ShopAssistant string
	Warehouse     string
	OrderId       uuid.UUID
}
