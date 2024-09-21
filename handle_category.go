package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Qu-Ack/kanteen_api/internal/database"
)

func (apiconfig apiConfig) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := apiconfig.DB.GetCategories(r.Context())
	if err != nil {
		log.Println("Error In HandleGetCategories while getting categories from DB", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 200, categories)
}

func (apiconfig apiConfig) HandlePostCategory(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Name string `json:"name"`
	}

	byte_body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error In HandlePostCategory while reading body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	json_body := body{}

	err = json.Unmarshal(byte_body, &json_body)

	if err != nil {
		log.Println("Error In HandlePostCategory while unmarsheling body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	_, err = apiconfig.DB.CreateCategory(r.Context(), json_body.Name)

	if err != nil {
		log.Println("Error In HandlePostCategory while creating category in Db", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})
}
func (apiconfig apiConfig) HandleUpdateCategory(w http.ResponseWriter, r *http.Request) {

	type body struct {
		Name       string `json:"name"`
		CategoryID int    `json:"category_id"`
	}

	byte_body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error In HandleUpdateCategory while reading body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	json_body := body{}

	err = json.Unmarshal(byte_body, &json_body)

	if err != nil {
		log.Println("Error In HandleUpdateCategory while unmarsheling body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	err = apiconfig.DB.UpdateCategory(r.Context(), database.UpdateCategoryParams{
		Name: json_body.Name,
		ID:   int32(json_body.CategoryID),
	})

	if err != nil {
		log.Println("Error In HandleUpdateCategory while creating category in Db", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})
}
func (apiconfig apiConfig) HandleDeleteCategory(w http.ResponseWriter, r *http.Request) {

	type body struct {
		CategoryID int `json:"category_id"`
	}

	byte_body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error In HandleDeleteCategory while reading body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	json_body := body{}

	err = json.Unmarshal(byte_body, &json_body)

	if err != nil {
		log.Println("Error In HandleDeleteCategory while unmarsheling body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	err = apiconfig.DB.DeleteCategory(r.Context(), int32(json_body.CategoryID))

	if err != nil {
		log.Println("Error In HandleUpdateCategory while creating category in Db", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})
}
