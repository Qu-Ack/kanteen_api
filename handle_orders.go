package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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
		Items map[string]Item `json:"items"`
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
