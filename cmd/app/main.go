package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	client "orders-center/internal/client"
	historyRepo "orders-center/internal/domain/history/repository"
	historySvc "orders-center/internal/domain/history/service"
	inboxRepo "orders-center/internal/domain/inbox/repository"
	inboxSvc "orders-center/internal/domain/inbox/service"
	orderRepo "orders-center/internal/domain/order/repository"
	orderSvc "orders-center/internal/domain/order/service"
	itemRepo "orders-center/internal/domain/order_item/repository"
	itemSvc "orders-center/internal/domain/order_item/service"
	outboxRepo "orders-center/internal/domain/outbox/repository"
	outboxSvc "orders-center/internal/domain/outbox/service"
	paymentRepo "orders-center/internal/domain/payment/repository"
	paymentSvc "orders-center/internal/domain/payment/service"
	"orders-center/internal/handler"
	"orders-center/internal/service/cron"
	"orders-center/internal/service/order_eno_1c"
	orderFullSvc "orders-center/internal/service/order_full/order_full_service"
	transactional "orders-center/internal/service/transactional"
	"orders-center/internal/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	dsn := "postgres://root:secret@localhost:5432/db?sslmode=disable"

	// Pool config
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(fmt.Errorf("db config parse error: %w", err))
	}

	cfg.MaxConns = 100
	cfg.MinConns = 5
	cfg.MaxConnIdleTime = 10 * time.Minute
	cfg.MaxConnLifetime = 1 * time.Hour

	// Creating pool with config
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
	inboxRepository := inboxRepo.NewInboxRepository(pool)
	inboxService := inboxSvc.NewInboxService(inboxRepository)
	orderService := orderSvc.NewOrderService(orderRepository)
	orderItemService := itemSvc.NewOrderItemService(itemRepository)
	historyService := historySvc.NewHistoryService(historyRepository)
	paymentService := paymentSvc.NewPaymentService(paymentRepository)
	outboxService := outboxSvc.NewOutboxService(outboxRepository)

	//transactional
	txService := transactional.NewTransactionService(pool)

	//orderfull service
	orderFullService := orderFullSvc.NewOrderFullService(orderService, orderItemService, paymentService, historyService, outboxService)

	//HTTP CLIENT
	clientCfg := client.ClientConfig{
		BaseURL: "http://localhost:1234",
		Timeout: 5 * time.Second,
	}

	clientForEno := client.NewClient(clientCfg)

	// 3. Cron + order_eno_1c
	cronScheduler := cron.NewScheduler(5)
	enoService := order_eno_1c.NewOrderEno1c(
		cronScheduler,
		txService,
		orderFullService,
		outboxService,
		inboxService,
		clientForEno,
		3,
	)

	// Usecase
	createOrderUC := usecase.NewCreateOrderUseCase(
		enoService,
		orderService,
		orderItemService,
		paymentService,
		historyService,
		txService,
	)

	//handler
	orderHandler := handler.NewOrderHandler(createOrderUC)

	// ctx для graceful shutdown всех фоновых процессов
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Запуск фонового processing
	go enoService.Run(ctx)

	r := gin.Default()
	r.POST("/orders", orderHandler.CreateOrderFull)

	// Launch HTTP-server с graceful
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
	//go PostOrderFull(ctx)
	<-ctx.Done()
	//cfg := graceful.NewShutDownConfig(5*time.Second, enoService.Reset, enoService.Stop)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Shutting down HTTP server…")
	if err = srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP shutdown error: %v\n", err)
	}

	fmt.Println("Service stopped")
}
