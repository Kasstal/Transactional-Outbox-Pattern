package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"orders-center/cmd/cron"
	"orders-center/cmd/order_eno_1c"
	"orders-center/cmd/order_full/entity"
	transactional "orders-center/cmd/transactional"
	"orders-center/cmd/usecase"
	db "orders-center/db/sqlc"
	outboxSvc "orders-center/internal/domain/outbox/service"

	/*historySvc "orders-center/internal/domain/history/service"
	orderSvc "orders-center/internal/domain/order/service"
	itemSvc "orders-center/internal/domain/order_item/service"
	outboxSvc "orders-center/internal/domain/outbox/service"
	paymentSvc "orders-center/internal/domain/payment/service"*/
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. Настроить пул соединений с параметрами
	dsn := "postgres://root:secret@localhost:5432/db?sslmode=disable"

	// Создаем конфигурацию пула соединений
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(fmt.Errorf("db config parse error: %w", err))
	}
	// Настроим параметры пула
	cfg.MaxConns = 200                     // Увеличиваем максимальное количество соединений в пуле
	cfg.MinConns = 5                       // Минимальное количество соединений
	cfg.MaxConnIdleTime = 10 * time.Minute // Время бездействия соединений
	cfg.MaxConnLifetime = 1 * time.Hour    // Максимальное время жизни соединения

	// Создание пула соединений с настроенной конфигурацией
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(fmt.Errorf("db connect: %w", err))
	}
	defer pool.Close()
	q := db.New(pool)
	/*orderService := orderSvc.NewOrderService(q)
	orderItemService := itemSvc.NewOrderItemService(q)
	historyService := historySvc.NewHistoryService(q)
	paymentService := paymentSvc.NewPaymentService(q)*/
	outboxService := outboxSvc.NewOutboxService(q)

	// 4. Транзакционный сервис
	txService := transactional.NewTransactionService(pool)

	// 5. UseCase для создания заказа
	createOrderUC := usecase.NewCreateOrderUseCase(
		txService,
	)

	// 7. Cron + order_eno_1c
	cronScheduler := cron.NewScheduler(6)
	enoService := order_eno_1c.NewOrderEno1c(
		cronScheduler,
		txService,
		outboxService,
	)

	// ctx для graceful shutdown всех фоновых процессов
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 8. Запуск фонового processing
	go enoService.Run(ctx)

	// 9. Настройка Gin
	r := gin.Default()
	r.POST("/orders", func(c *gin.Context) {
		var of entity.OrderFull
		if err := c.ShouldBindJSON(&of); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := createOrderUC.Create(c.Request.Context(), of); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	// 10. Запуск HTTP-сервера с graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		fmt.Println("Gin server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Gin serve error: %v\n", err)
			stop()
		}
	}()
	go PostOrderFull(ctx, q)
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Shutting down HTTP server…")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP shutdown error: %v\n", err)
	}
	fmt.Println("Service stopped")
}
