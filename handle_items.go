package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Qu-Ack/kanteen_api/internal/database"
)

func (apiconfig apiConfig) HandleGetItems(w http.ResponseWriter, r *http.Request) {
	items, err := apiconfig.DB.GetItems(r.Context())
	if err != nil {
		log.Println("Error In HandleGetItems while reading items from DB", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}
	WriteJSON(w, 200, items)

}

func (apiconfig apiConfig) HandleUpdateItem(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Name       string `json:"name"`
		Price      int    `json:"price"`
		Stock      int    `json:"stock"`
		CategoryID int    `json:"category_id"`
		ItemID     string `json:"item_id"`
	}

	byte_body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error In HandleUpdateItem while reading body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	json_body := body{}

	err = json.Unmarshal(byte_body, &json_body)

	if err != nil {
		log.Println("Error In HandleUpdateItem while unmarsheling body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	my_item_id, err := strconv.Atoi(json_body.ItemID)
	if err != nil {
		log.Println("Error In HandleUpdateItem while converting item id to int", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	err = apiconfig.DB.UpdateItem(r.Context(), database.UpdateItemParams{
		Name:       json_body.Name,
		Price:      int32(json_body.Price),
		Stock:      int32(json_body.Stock),
		CategoryID: int32(json_body.CategoryID),
		ID:         int32(my_item_id),
	})

	if err != nil {
		log.Println("Error In HandleUpdateItem while writing item to db", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})

}
func (apiconfig apiConfig) HandleDeleteItem(w http.ResponseWriter, r *http.Request) {
	type body struct {
		ItemID int `json:"item_id"`
	}

	byte_body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error In HandleDeleteItem while reading body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	json_body := body{}

	err = json.Unmarshal(byte_body, &json_body)

	if err != nil {
		log.Println("Error In HandleDeleteItem while unmarsheling body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	err = apiconfig.DB.DeleteItem(r.Context(), int32(json_body.ItemID))

	if err != nil {
		log.Println("Error In HandleDeleteItem while writing item to db", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})

}
func (apiconfig apiConfig) HandleCreateItem(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Name       string `json:"name"`
		Price      int    `json:"price"`
		Stock      int    `json:"stock"`
		CategoryID int    `json:"category_id"`
	}

	byte_body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error In HandleCreateItem while reading body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	json_body := body{}

	err = json.Unmarshal(byte_body, &json_body)

	if err != nil {
		log.Println("Error In HandleCreateItem while unmarsheling body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	_, err = apiconfig.DB.CreateItem(r.Context(), database.CreateItemParams{
		Name:       json_body.Name,
		Price:      int32(json_body.Price),
		Stock:      int32(json_body.Stock),
		CategoryID: int32(json_body.CategoryID),
	})

	if err != nil {
		log.Println("Error In HandleCreateItem while writing item to db", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})

}
