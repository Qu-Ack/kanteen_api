package main

import (
	"log"
	"net/http"
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
