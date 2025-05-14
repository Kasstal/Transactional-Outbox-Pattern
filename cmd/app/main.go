package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"orders-center/cmd/handler"
	"orders-center/cmd/usecase"
	historyRepo "orders-center/internal/domain/history/repository"
	historySvc "orders-center/internal/domain/history/service"
	orderRepo "orders-center/internal/domain/order/repository"
	orderSvc "orders-center/internal/domain/order/service"
	itemRepo "orders-center/internal/domain/order_item/repository"
	itemSvc "orders-center/internal/domain/order_item/service"
	outboxRepo "orders-center/internal/domain/outbox/repository"
	outboxSvc "orders-center/internal/domain/outbox/service"
	paymentRepo "orders-center/internal/domain/payment/repository"
	paymentSvc "orders-center/internal/domain/payment/service"
	"orders-center/internal/service/cron"
	"orders-center/internal/service/order_eno_1c"
	orderFullSvc "orders-center/internal/service/order_full/order_full_service"
	transactional "orders-center/internal/service/transactional"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	dsn := "postgres://root:secret@localhost:5432/db?sslmode=disable"

	// Создаем конфигурацию пула соединений
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(fmt.Errorf("db config parse error: %w", err))
	}

	cfg.MaxConns = 100
	cfg.MinConns = 5
	cfg.MaxConnIdleTime = 10 * time.Minute
	cfg.MaxConnLifetime = 1 * time.Hour

	// Создание пула соединений с настроенной конфигурацией
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(fmt.Errorf("db connect: %w", err))
	}
	defer pool.Close()

	orderRepository := orderRepo.NewOrderRepository(pool)
	historyRepository := historyRepo.NewHistoryRepository(pool)
	itemRepository := itemRepo.NewOrderItemRepository(pool)
	outboxRepository := outboxRepo.NewOutboxRepository(pool)
	paymentRepository := paymentRepo.NewPaymentRepository(pool)
	orderService := orderSvc.NewOrderService(orderRepository)
	orderItemService := itemSvc.NewOrderItemService(itemRepository)
	historyService := historySvc.NewHistoryService(historyRepository)
	paymentService := paymentSvc.NewPaymentService(paymentRepository)
	outboxService := outboxSvc.NewOutboxService(outboxRepository)

	//transactional
	txService := transactional.NewTransactionService(pool)

	// Usecase
	createOrderUC := usecase.NewCreateOrderUseCase(
		orderService,
		orderItemService,
		paymentService,
		historyService,
		outboxService,
		txService,
	)

	//orderfull service
	orderFullService := orderFullSvc.NewOrderFullService(orderService, orderItemService, paymentService, historyService, outboxService)

	//handler
	orderHandler := handler.NewOrderHandler(createOrderUC)

	// 3. Cron + order_eno_1c
	cronScheduler := cron.NewScheduler(15)
	enoService := order_eno_1c.NewOrderEno1c(
		cronScheduler,
		txService,
		orderFullService,
		outboxService,
	)

	// ctx для graceful shutdown всех фоновых процессов
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 3. Запуск фонового processing
	go enoService.Run(ctx)

	r := gin.Default()
	r.POST("/orders", orderHandler.CreateOrderFull)

	// 4. Запуск HTTP-сервера с graceful
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
	go PostOrderFull(ctx)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Shutting down HTTP server…")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP shutdown error: %v\n", err)
	}
	enoService.Reset()

	fmt.Println("Service stopped")
}
