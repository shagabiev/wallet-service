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
	appPort := os.Getenv("APP_PORT")

	db, err := storage.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(db)
	h := &handler.WalletHandler{Service: svc}

	http.HandleFunc("/api/v1/wallet", h.UpdateBalance)
	http.HandleFunc("/api/v1/wallets/", h.GetBalance)

	log.Printf("Server started at :%s", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, nil))
}
