package router

import (
	"github.com/gin-gonic/gin"
	"orders-center/internal/handler"
)

type Router struct {
	handler *handler.OrderHandler
}

func NewRouter(orderHandler *handler.OrderHandler) *Router {
	return &Router{
		handler: orderHandler,
	}
}

func (r *Router) InitRouter(endpoint string) *gin.Engine {
	router := gin.New()
	//middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	//routes
	router.POST(endpoint, r.handler.CreateOrderFull)

	return router
}
