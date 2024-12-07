package handler

import (
	"encoding/json"
	"net/http"
	handler_error "slot-machine/internal/adapters/http/handler/error"
	"slot-machine/internal/adapters/http/middleware"
	"slot-machine/internal/application/usecase"
	"slot-machine/internal/domain/repository"
)

type Handler struct {
	CreatePlayerUseCase          *usecase.CreatePlayerUseCase
	CreateSlotMachineUseCase     *usecase.CreateSlotMachineUseCase
	PlayUseCase                  *usecase.PlayUseCase
	GetPlayerBalanceUseCase      *usecase.GetPlayerBalanceUseCase
	GetSlotMachineBalanceUseCase *usecase.GetSlotMachineBalanceUseCase
	loginUseCase                 *usecase.LoginUseCase
}

func NewHandler(
	cpUC *usecase.CreatePlayerUseCase,
	csmUC *usecase.CreateSlotMachineUseCase,
	pUC *usecase.PlayUseCase,
	gpUC *usecase.GetPlayerBalanceUseCase,
	gsmUC *usecase.GetSlotMachineBalanceUseCase,
	loginUC *usecase.LoginUseCase,
) *Handler {
	return &Handler{
		CreatePlayerUseCase:          cpUC,
		CreateSlotMachineUseCase:     csmUC,
		PlayUseCase:                  pUC,
		GetPlayerBalanceUseCase:      gpUC,
		GetSlotMachineBalanceUseCase: gsmUC,
		loginUseCase:                 loginUC,
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
// @Failure 400 {object} handler_error.HTTPError "Payload inválido"
// @Failure 404 {object} handler_error.HTTPError "Máquina caça-níqueis não encontrada"
// @Failure 422 {object} handler_error.HTTPError "Saldo insuficiente"
// @Failure 500 {object} handler_error.HTTPError "Erro interno do servidor"
// @Router /play [post]
// @Security BearerAuth
func (h *Handler) PlaySlotMachine(w http.ResponseWriter, r *http.Request) {
	var req usecase.PlayRequest
	w.Header().Set("Content-Type", "application/json")

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(handler_error.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handler_error.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
		return
	}

	req.PlayerID = userID

	response, err := h.PlayUseCase.Execute(r.Context(), &req)
	if err != nil {
		handler_error.HandleError(w, err)
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
// @Failure 400 {object} handler_error.HTTPError "Payload inválido"
// @Failure 409 {object} handler_error.HTTPError "Jogador já existe"
// @Failure 500 {object} handler_error.HTTPError "Erro interno do servidor"
// @Router /players [post]
func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreatePlayerRequest
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handler_error.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})

		return
	}

	resp, err := h.CreatePlayerUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == usecase.ErrPlayerAlreadyExists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(handler_error.HTTPError{
				Code:    http.StatusConflict,
				Message: "Player already exists",
			})

			return
		} else if err == usecase.ErrValidate {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(handler_error.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "email and password must be provided",
			})

			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handler_error.HTTPError{
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
// @Failure 400 {object} handler_error.HTTPError "Payload inválido ou parâmetros inválidos"
// @Failure 500 {object} handler_error.HTTPError "Erro interno do servidor"
// @Router /machines [post]
// @Security AdminAuth
func (h *Handler) CreateSlotMachine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req usecase.CreateSlotMachineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handler_error.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})

		return
	}

	resp, err := h.CreateSlotMachineUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == usecase.ErrUnauthorized {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(handler_error.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			return
		}
		if err == usecase.ErrSlotMachineAlreadyExists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(handler_error.HTTPError{
				Code:    http.StatusConflict,
				Message: "Slot machine already exists",
			})

			return
		} else if err == usecase.ErrValidate {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(handler_error.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "level, multiple gain, and description must be provided",
			})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handler_error.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to create slot machine",
		})

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// LoginHandler realiza a autenticação do usuário e retorna um token JWT.
// @Summary Login
// @Description Autentica um usuário e retorna um token JWT.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param loginRequest body usecase.LoginRequest true "Dados de autenticação"
// @Success 200 {object} usecase.LoginResponse
// @Failure 400 {object} handler_error.HTTPError "Requisição inválida"
// @Failure 401 {object} handler_error.HTTPError "Credenciais inválidas"
// @Failure 500 {object} handler_error.HTTPError "Erro interno do servidor"
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req usecase.LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.loginUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetPlayerBalance retorna o saldo do jogador.
// @Summary Obter saldo do jogador
// @Description Retorna o saldo do jogador especificado.
// @Tags Player
// @Accept json
// @Produce json
// @Success 200 {object} usecase.GetPlayerBalanceResponse "Saldo do jogador"
// @Failure 401 {object} handler_error.HTTPError "Não autorizado"
// @Failure 404 {object} handler_error.HTTPError "Jogador não encontrado"
// @Failure 500 {object} handler_error.HTTPError "Erro interno do servidor"
// @Router /players/balance [get]
// @Security BearerAuth
func (h *Handler) GetPlayerBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(handler_error.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	req := usecase.GetPlayerBalanceRequest{
		PlayerID: userID,
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

// GetSlotMachineBalance retorna o saldo da máquina caça-níqueis.
// @Summary Obter saldo da máquina caça-níqueis
// @Description Retorna o saldo da máquina caça-níqueis especificada.
// @Tags SlotMachine
// @Accept json
// @Produce json
// @Param machine_id query string true "ID da máquina caça-níqueis"
// @Success 200 {object} usecase.GetSlotMachineBalanceResponse "Saldo da máquina caça-níqueis"
// @Failure 401 {object} handler_error.HTTPError "Não autorizado"
// @Failure 400 {object} handler_error.HTTPError "ID da máquina caça-níqueis é obrigatório"
// @Failure 404 {object} handler_error.HTTPError "Máquina caça-níqueis não encontrada"
// @Failure 500 {object} handler_error.HTTPError "Erro interno do servidor"
// @Router /machines/balance [get]
// @Security AdminAuth
func (h *Handler) GetSlotMachineBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	machineID := r.URL.Query().Get("machine_id")
	if machineID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handler_error.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "machine_id is required",
		})

		return
	}

	req := usecase.GetSlotMachineBalanceRequest{
		MachineID: machineID,
	}

	resp, err := h.GetSlotMachineBalanceUseCase.Execute(r.Context(), &req)
	if err != nil {
		if err == usecase.ErrUnauthorized {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(handler_error.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
			return
		}
		if err == repository.ErrSlotMachineNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(handler_error.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Slot machine not found",
			})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handler_error.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to get slot machine balance",
		})

		return
	}

	json.NewEncoder(w).Encode(resp)
}
