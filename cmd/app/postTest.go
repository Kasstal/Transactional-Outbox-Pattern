package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	orderFullSvc "orders-center/cmd/order_full/order_full_service"
	db "orders-center/db/sqlc"
	"orders-center/internal/utils"
	"time"
)

func PostOrderFull(ctx context.Context, q *db.Queries) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

			order_full := utils.RandomOrderFull()
			data, err := json.Marshal(order_full)
			req, err := http.NewRequest("POST", "http://localhost:8080/orders", bytes.NewBuffer(data))
			if err != nil {
				log.Printf("failed to create request: %v", err)
			}

			// Set the appropriate headers
			req.Header.Set("Content-Type", "application/json")

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("failed to send request: %v", err)
			}

			if resp.StatusCode != http.StatusCreated {
				log.Printf("server returned non-created status: %v", resp.Status)
			}

			log.Println("Order successfully posted to Orders-Center!")
			resp.Body.Close()
			orderFull := orderFullSvc.NewOrderFullService(q)
			orderCheck, err := orderFull.GetOrderFull(ctx, order_full.Order.ID)
			if err != nil {
				log.Printf("failed to get order: %v", err)
			}
			log.Println(orderCheck.Order.ID)
		}
		time.Sleep(2 * time.Second)
	}

}
