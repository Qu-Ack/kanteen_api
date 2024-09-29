package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/Qu-Ack/kanteen_api/internal/database"
	"github.com/google/uuid"
)

type Item struct {
	Name             string `json:"name"`
	Price            int    `json:"price"`
	ID               int    `json:"id"`
	EatInQuantity    int    `json:"eatInQuantity"`
	TakeAwayQuantity int    `json:"takeAwayQuantity"`
	ServiceType      string `json:"serviceType"`
}

func (apiconfig apiConfig) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	order_id := r.PathValue("ORDERID")
	int_order_id, err := strconv.Atoi(order_id)
	if err != nil {
		log.Println("Error In HandleGetOrder while converting id to int", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}
	order, err := apiconfig.DB.GetOrder(r.Context(), int32(int_order_id))

	if err != nil {
		log.Println("Error In HandleGetOrder while getting order from db", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	order_items, err := apiconfig.DB.GetOrderItemsForOrder(r.Context(), int32(int_order_id))

	if err != nil {
		log.Println("Error In HandleGetOrder while getting order items", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 200, map[string]any{
		"order_id":     order.OrderID,
		"user_name":    order.UserName,
		"user_mobile":  order.UserMobile,
		"total":        order.Total,
		"order_status": order.OrderStatus,
		"items":        order_items,
	})

}

func (apiconfig apiConfig) HandlePostOrder(w http.ResponseWriter, r *http.Request) {
	type body struct {
		UserId string          `json:"user_id"`
		Items  map[string]Item `json:"items"`
	}
	byte_body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error In HandlePostOrder while reading body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	json_body := body{}

	err = json.Unmarshal(byte_body, &json_body)

	if err != nil {
		log.Println("Error In HandlePostOrder while Unmarsheling body ", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}
	total := calculateTotal(json_body.Items)

	total_string := strconv.FormatFloat(total, 'f', 2, 64)

	user_uuid, err := uuid.Parse(json_body.UserId)

	if err != nil {
		log.Println("Error In HandlePostOrder while parsing uuid", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	our_user, err := apiconfig.DB.GetUserByID(r.Context(), uuid.NullUUID{
		UUID:  user_uuid,
		Valid: true,
	})

	order, err := apiconfig.DB.CreateOrder(r.Context(), database.CreateOrderParams{
		UserID: user_uuid,
		Total:  total_string,
		Status: "pending",
	})

	var wg sync.WaitGroup
	for _, item := range json_body.Items {
		wg.Add(1)
		go func(item Item) {
			defer wg.Done()
			price_str := strconv.FormatInt(int64(item.Price), 10)

			_, err := apiconfig.DB.CreateOrderItem(r.Context(), database.CreateOrderItemParams{
				OrderID:          order.ID,
				ItemID:           int32(item.ID),
				TakeawayQuantity: int32(item.TakeAwayQuantity),
				EatinQuantity:    int32(item.EatInQuantity),
				Price:            price_str,
			})
			if err != nil {
				log.Println("Error inserting order item into database:", err)
				WriteJSONError(w, 500, "Internal Server Error")
				return
			}
		}(item)
	}

	wg.Wait()

	if SocketHandler.master != nil {
		err := SocketHandler.master.WriteJSON(map[string]any{"user_id": json_body.UserId, "order_id": order.ID, "user_mobile": our_user.Phone, "user_name": our_user.Name})
		if err != nil {
			log.Println("Error In HandlePostOrder while sending data to socket client", err)
			return
		}
		log.Println("message sent success")
	} else {
		log.Println("Error In HandlePostOrder websocket conn not established")
	}

	WriteJSON(w, 201, map[string]any{"status": "ok", "order_id": order.ID})

}

func (apiconfig apiConfig) HandleGetPendingOrder(w http.ResponseWriter, r *http.Request) {
	orders, err := apiconfig.DB.GetPendingOrders(r.Context())

	if err != nil {
		log.Println("Error In HandleGetPendingOrder while fetching orders")
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 500, orders)

}

func calculateTotal(items map[string]Item) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Price) * float64(item.EatInQuantity+item.TakeAwayQuantity)
	}
	return total
}
