package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Qu-Ack/kanteen_api/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error while loading env variables")
	}

	PORT := os.Getenv("PORT")
	DB_STRING := os.Getenv("DB_STRING")

	serve_mux := http.NewServeMux()

	db, err := sql.Open("postgres", DB_STRING)

	if err != nil {
		log.Println("Error while opening DB")
		return
	}

	dbQueries := database.New(db)

	apiconfig := apiConfig{
		DB: dbQueries,
	}

	serve_mux.HandleFunc("GET /category", apiconfig.HandleGetCategories)
	serve_mux.HandleFunc("POST /category", apiconfig.HandlePostCategory)
	serve_mux.HandleFunc("PUT /category", apiconfig.HandleUpdateCategory)
	serve_mux.HandleFunc("DELETE /category", apiconfig.HandleDeleteCategory)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: serve_mux,
	}

	log.Println(fmt.Sprintf("Server started on PORT %v", PORT))
	server.ListenAndServe()
}
