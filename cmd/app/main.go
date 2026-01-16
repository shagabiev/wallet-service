package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shagabiev/wallet-service/internal/handler"
	"github.com/shagabiev/wallet-service/internal/service"
	"github.com/shagabiev/wallet-service/internal/storage"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	db, err := storage.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(db)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/wallet", h.UpdateWallet)
	mux.HandleFunc("GET /api/v1/wallets/{id}", h.GetBalance)

	log.Printf("listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
