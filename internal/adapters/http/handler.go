package http

import (
	"encoding/json"
	"net/http"
	"slot-machine/internal/application/usecase"
	"slot-machine/internal/domain/repository"
)

// HTTPError representa um erro retornado pela API.
// @Description Estrutura para representar erros na API.
// @Description Contém a mensagem de erro e um código opcional.
// @Description Pode ser expandida conforme necessário.
type HTTPError struct {
	Code    int    `json:"code"`    // Código do erro HTTP
	Message string `json:"message"` // Mensagem descritiva do erro
}

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

// PlaySlotMachine permite que o jogador jogue na máquina caça-níqueis.
// @Summary Jogar na máquina caça-níqueis
// @Description Permite que o jogador faça uma aposta e jogue na máquina caça-níqueis especificada.
// @Tags SlotMachine
// @Accept json
// @Produce json
// @Param playRequest body usecase.PlayRequest true "Dados da jogada"
// @Success 200 {object} usecase.PlayResponse "Jogada realizada com sucesso"
// @Failure 400 {object} HTTPError "Payload inválido"
// @Failure 404 {object} HTTPError "Máquina caça-níqueis não encontrada"
// @Failure 422 {object} HTTPError "Saldo insuficiente"
// @Failure 500 {object} HTTPError "Erro interno do servidor"
// @Router /play [post]
func (h *Handler) PlaySlotMachine(w http.ResponseWriter, r *http.Request) {
	var req usecase.PlayRequest
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})

		return
	}

	response, err := h.PlayUseCase.Execute(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(response)
}

// CreatePlayer permite a criação de um novo jogador.
// @Summary Criar um novo jogador
// @Description Permite a criação de um novo jogador com um saldo inicial.
// @Tags Player
// @Accept json
// @Produce json
// @Param createPlayerRequest body usecase.CreatePlayerRequest true "Dados do jogador a ser criado"
// @Success 201 {object} usecase.CreatePlayerResponse "Jogador criado com sucesso"
// @Failure 400 {object} HTTPError "Payload inválido"
// @Failure 409 {object} HTTPError "Jogador já existe"
// @Failure 500 {object} HTTPError "Erro interno do servidor"
// @Router /players [post]
func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreatePlayerRequest
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})

		return
	}

	resp, err := h.CreatePlayerUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == usecase.ErrPlayerAlreadyExists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(HTTPError{
				Code:    http.StatusConflict,
				Message: "Player already exists",
			})

			return
		} else if err == usecase.ErrValidate {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HTTPError{
				Code:    http.StatusBadRequest,
				Message: "email and password must be provided",
			})

			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to create player",
		})

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// CreateSlotMachine permite a criação de uma nova máquina caça-níqueis.
// @Summary Criar uma nova máquina caça-níqueis
// @Description Permite a criação de uma nova máquina caça-níqueis com os parâmetros especificados.
// @Tags SlotMachine
// @Accept json
// @Produce json
// @Param createSlotMachineRequest body usecase.CreateSlotMachineRequest true "Dados da máquina caça-níqueis a ser criada"
// @Success 201 {object} usecase.CreateSlotMachineResponse "Máquina criada com sucesso"
// @Failure 400 {object} HTTPError "Payload inválido ou parâmetros inválidos"
// @Failure 500 {object} HTTPError "Erro interno do servidor"
// @Router /machines [post]
func (h *Handler) CreateSlotMachine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req usecase.CreateSlotMachineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})

		return
	}

	resp, err := h.CreateSlotMachineUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == usecase.ErrSlotMachineAlreadyExists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(HTTPError{
				Code:    http.StatusConflict,
				Message: "Slot machine already exists",
			})

			return
		} else if err == usecase.ErrValidate {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HTTPError{
				Code:    http.StatusBadRequest,
				Message: "level, multiple gain, and description must be provided",
			})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to create slot machine",
		})

		return
	}

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
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Insufficient balance",
		})
	case repository.ErrPlayerNotFound, repository.ErrSlotMachineNotFound:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(HTTPError{
			Code:    http.StatusNotFound,
			Message: "Resource not found",
		})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
	}
}
