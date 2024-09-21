package main

import (
	"encoding/json"
	"fmt"
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

	for key, value := range json_body.Items {
		fmt.Println(fmt.Sprintf("%v: %v", key, value.Name))
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})

}
