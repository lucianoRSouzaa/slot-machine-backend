package http

import (
	"encoding/json"
	"net/http"
	"slot-machine/internal/application/usecase"
	"slot-machine/internal/domain/repository"
)

type Handler struct {
	CreatePlayerUseCase          *usecase.CreatePlayerUseCase
	CreateSlotMachineUseCase     *usecase.CreateSlotMachineUseCase
	PlayUseCase                  *usecase.PlayUseCase
	GetPlayerBalanceUseCase      *usecase.GetPlayerBalanceUseCase
	GetSlotMachineBalanceUseCase *usecase.GetSlotMachineBalanceUseCase
}

func NewHandler(
	cpUC *usecase.CreatePlayerUseCase,
	csmUC *usecase.CreateSlotMachineUseCase,
	pUC *usecase.PlayUseCase,
	gpUC *usecase.GetPlayerBalanceUseCase,
	gsmUC *usecase.GetSlotMachineBalanceUseCase,
) *Handler {
	return &Handler{
		CreatePlayerUseCase:          cpUC,
		CreateSlotMachineUseCase:     csmUC,
		PlayUseCase:                  pUC,
		GetPlayerBalanceUseCase:      gpUC,
		GetSlotMachineBalanceUseCase: gsmUC,
	}
}

func (h *Handler) PlaySlotMachine(w http.ResponseWriter, r *http.Request) {
	var req usecase.PlayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response, err := h.PlayUseCase.Execute(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.CreatePlayerUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == usecase.ErrPlayerAlreadyExists {
			http.Error(w, "Player already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Unable to create player", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) CreateSlotMachine(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreateSlotMachineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.CreateSlotMachineUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == usecase.ErrSlotMachineAlreadyExists {
			http.Error(w, "Slot machine already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Unable to create slot machine", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetPlayerBalance(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	req := usecase.GetPlayerBalanceRequest{
		PlayerID: playerID,
	}

	resp, err := h.GetPlayerBalanceUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == repository.ErrPlayerNotFound {
			http.Error(w, "Player not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to get player balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetSlotMachineBalance(w http.ResponseWriter, r *http.Request) {
	machineID := r.URL.Query().Get("machine_id")
	if machineID == "" {
		http.Error(w, "machine_id is required", http.StatusBadRequest)
		return
	}

	req := usecase.GetSlotMachineBalanceRequest{
		MachineID: machineID,
	}

	resp, err := h.GetSlotMachineBalanceUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == repository.ErrSlotMachineNotFound {
			http.Error(w, "Slot machine not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to get slot machine balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleError(w http.ResponseWriter, err error) {
	switch err {
	case usecase.ErrInsufficientBalance:
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
	case repository.ErrPlayerNotFound, repository.ErrSlotMachineNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
