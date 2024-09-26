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

	order, err := apiconfig.DB.CreateOrder(r.Context(), database.CreateOrderParams{
		UserID: user_uuid,
		Total:  total_string,
		Status: "pending",
	})

	// Use a WaitGroup to wait for all insertions to complete
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

	// sending data to the web client
	if SocketHandler.master != nil {
		err := SocketHandler.master.WriteJSON(json_body)
		if err != nil {
			log.Println("Error In HandlePostOrder while sending data to socket client", err)
			return
		}
		log.Println("message sent success")
	} else {
		log.Println("Error In HandlePostOrder websocket conn not established")
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})

}

func calculateTotal(items map[string]Item) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Price) * float64(item.EatInQuantity+item.TakeAwayQuantity)
	}
	return total
}
