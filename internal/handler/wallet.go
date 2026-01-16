package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/shagabiev/wallet-service/internal/service"
)

type WalletHandler struct {
	Service *service.Service
}

type WalletRequest struct {
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        int64     `json:"amount"`
}

type WalletResponse struct {
	WalletID uuid.UUID `json:"walletId"`
	Balance  int64     `json:"balance"`
}

func (h *WalletHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateBalance(r.Context(), req.WalletID, req.OperationType, req.Amount); err != nil {
		if err == service.ErrInsufficientFunds {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, _ := h.Service.GetBalance(r.Context(), req.WalletID)
	json.NewEncoder(w).Encode(WalletResponse{WalletID: req.WalletID, Balance: balance})
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	walletIDStr := r.URL.Path[len("/api/v1/wallets/"):]
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balance, err := h.Service.GetBalance(r.Context(), walletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(WalletResponse{WalletID: walletID, Balance: balance})
}
