package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error while loading env variables")
	}

	PORT := os.Getenv("PORT")
	_ = os.Getenv("DB_STRING")

	serve_mux := http.NewServeMux()

	serve_mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: serve_mux,
	}

	log.Println(fmt.Sprintf("Server started on PORT %v", PORT))
	server.ListenAndServe()
}
