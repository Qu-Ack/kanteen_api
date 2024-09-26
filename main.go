package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Qu-Ack/kanteen_api/internal/database"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB     *database.Queries
	master *websocket.Conn
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

	serve_mux.HandleFunc("/ws", apiconfig.HandleWebSocketConn)

	serve_mux.HandleFunc("GET /category", apiconfig.HandleGetCategories)
	serve_mux.HandleFunc("POST /category", apiconfig.HandlePostCategory)
	serve_mux.HandleFunc("PUT /category", apiconfig.HandleUpdateCategory)
	serve_mux.HandleFunc("DELETE /category", apiconfig.HandleDeleteCategory)

	serve_mux.HandleFunc("GET /item", apiconfig.HandleGetItems)
	serve_mux.HandleFunc("POST /item", apiconfig.HandleCreateItem)
	serve_mux.HandleFunc("PUT /item", apiconfig.HandleUpdateItem)
	serve_mux.HandleFunc("DELETE /item", apiconfig.HandleDeleteItem)

	serve_mux.HandleFunc("POST /otp", apiconfig.HandleCreateOTP)
	serve_mux.HandleFunc("POST /verifyotp", apiconfig.HandleVerifyOTP)

	serve_mux.HandleFunc("POST /order", apiconfig.HandlePostOrder)

	// Wrap the mux with the CORS middleware
	corsHandler := enableCORS(serve_mux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: corsHandler,
	}

	log.Println(fmt.Sprintf("Server started on PORT %v", PORT))
	server.ListenAndServe()
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust to your needs
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
