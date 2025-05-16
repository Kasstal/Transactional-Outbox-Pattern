package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	orderFull "orders-center/internal/service/order_full/entity"
	"orders-center/internal/usecase"
)

type OrderHandler struct {
	usecase usecase.UseCase
}

func NewOrderHandler(usecase usecase.UseCase) *OrderHandler {
	return &OrderHandler{usecase: usecase}
}

func (h *OrderHandler) CreateOrderFull(c *gin.Context) {

	var orderFull orderFull.OrderFull
	if err := c.ShouldBindJSON(&orderFull); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	//validation

	if err := h.usecase.Create(c.Request.Context(), orderFull); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Successfully accepted full order ")
	c.Status(http.StatusCreated)
}
