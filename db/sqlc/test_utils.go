package db

import (
	"github.com/gofrs/uuid"
	"math/rand"
	"strconv"
)

func randomString(prefix string) string {
	return prefix + strconv.Itoa(rand.Intn(1000))
}

// Вспомогательные функции для генерации случайных данных
func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomOrderType() string {
	types := []string{"online", "offline", "wholesale"}
	return types[rand.Intn(len(types))]
}

func randomOrderStatus() string {
	statuses := []string{"created", "processing", "completed", "cancelled"}
	return statuses[rand.Intn(len(statuses))]
}

func randomCity() string {
	cities := []string{"Moscow", "Saint Petersburg", "Novosibirsk", "Yekaterinburg"}
	return cities[rand.Intn(len(cities))]
}

func randomSubdivision() string {
	subdivisions := []string{"north", "south", "east", "west", "center"}
	return subdivisions[rand.Intn(len(subdivisions))]
}

func randomPlatform() string {
	platforms := []string{"web", "mobile", "terminal", "api"}
	return platforms[rand.Intn(len(platforms))]
}

func randomUUID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func randomOrderNumber() string {
	return "ORD-" + strconv.Itoa(rand.Intn(1000))
}

func randomUserID() string {
	return "user-" + strconv.Itoa(rand.Intn(1000))
}

func randomProductID() string {
	return "prod-" + strconv.Itoa(rand.Intn(1000))
}

func randomExternalID() string {
	return "ext-" + strconv.Itoa(rand.Intn(1000))
}

func randomItemStatus() string {
	statuses := []string{"reserved", "shipped", "delivered", "returned"}
	return statuses[rand.Intn(len(statuses))]
}

func randomDeliveryID() string {
	return "delivery-" + strconv.Itoa(rand.Intn(1000))
}

func randomWarehouse() string {
	return "warehouse-" + strconv.Itoa(rand.Intn(1000))
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomPaymentType() string {
	types := []string{"credit", "card", "cash", "bonuses"}
	return types[rand.Intn(len(types))]
}

func randomPaymentInfo() string {
	return "payment-" + strconv.Itoa(rand.Intn(1000))
}

func randomContractNumber() string {
	return "CNT-" + strconv.Itoa(rand.Intn(1000))
}

func randomBank() string {
	banks := []string{"sberbank", "tinkoff", "alfa", "vtb"}
	return banks[rand.Intn(len(banks))]
}

func randomTerm() string {
	return strconv.Itoa(rand.Intn(1000))
}

func randomLast4() string {
	return strconv.Itoa(rand.Intn(1000))
}

func randomCardSystem() string {
	systems := []string{"visa", "mastercard", "mir", "unionpay"}
	return systems[rand.Intn(len(systems))]
}

func randomHistoryType() string {
	types := []string{"order", "item", "payment", "user"}
	return types[rand.Intn(len(types))]
}
