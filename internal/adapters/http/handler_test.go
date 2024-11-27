package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	httpGo "net/http"
	"net/http/httptest"
	"slot-machine/internal/adapters/http"
	"slot-machine/internal/application/usecase"
	"slot-machine/internal/domain/model"
	repository_in_memory "slot-machine/internal/infrastructure/repository/in_memory"
	"slot-machine/internal/infrastructure/security"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SetupHandler() *http.Handler {
	playerRepo := repository_in_memory.NewInMemoryPlayerRepository()
	hasher := security.NewBcryptPasswordHasher(bcrypt.DefaultCost)

	createPlayerUC := usecase.NewCreatePlayerUseCase(playerRepo, hasher)

	handler := &http.Handler{
		CreatePlayerUseCase: createPlayerUC,
	}

	return handler
}

func TestCreatePlayerHandler(t *testing.T) {
	handler := SetupHandler()

	router := httpGo.NewServeMux()
	router.HandleFunc("/players", handler.CreatePlayer)

	t.Run("CreatePlayer_Success", func(t *testing.T) {
		reqBody := usecase.CreatePlayerRequest{
			Balance:  1000,
			Email:    "email",
			Password: "abc",
		}

		jsonBody, err := json.Marshal(reqBody)
		assert.NoError(t, err, "Erro ao serializar a requisição")

		req, err := httpGo.NewRequest("POST", "/players", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err, "Erro ao criar a requisição HTTP")

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, httpGo.StatusCreated, rr.Code, "Status code deve ser 201 Created")

		var resp usecase.CreatePlayerResponse
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err, "Erro ao deserializar a resposta")

		assert.Equal(t, reqBody.Email, resp.Player.Email, "ID do jogador deve corresponder ao solicitado")
		assert.Equal(t, reqBody.Balance, resp.Player.Balance, "Saldo do jogador deve corresponder ao solicitado")

		storedPlayer, err := handler.CreatePlayerUseCase.PlayerRepo.GetPlayer(context.Background(), resp.Player.ID)
		assert.NoError(t, err, "Erro ao recuperar o jogador do repositório")
		assert.Equal(t, resp.Player, *storedPlayer, "Jogador armazenado deve corresponder à resposta")
	})

	t.Run("CreatePlayer_AlreadyExists", func(t *testing.T) {
		initialPlayer := &model.Player{
			ID:       "player456",
			Balance:  500,
			Email:    "email",
			Password: "aaa",
		}

		err := handler.CreatePlayerUseCase.PlayerRepo.CreatePlayer(context.Background(), initialPlayer)
		assert.NoError(t, err, "Erro ao criar o jogador inicial no repositório")

		reqBody := usecase.CreatePlayerRequest{
			Balance:  1500,
			Email:    "email",
			Password: "aaa",
		}

		jsonBody, err := json.Marshal(reqBody)
		assert.NoError(t, err, "Erro ao serializar a requisição")

		req, err := httpGo.NewRequest("POST", "/players", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err, "Erro ao criar a requisição HTTP")

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, httpGo.StatusConflict, rr.Code, "Status code deve ser 409 Conflict")

		expectedError := HTTPError{
			Code:    httpGo.StatusConflict,
			Message: "Player already exists",
		}

		var actualError HTTPError
		err = json.Unmarshal(rr.Body.Bytes(), &actualError)
		assert.NoError(t, err, "Erro ao deserializar a resposta de erro")
		assert.Equal(t, expectedError, actualError, "Mensagem de erro deve corresponder ao esperado")

		storedPlayer, err := handler.CreatePlayerUseCase.PlayerRepo.GetPlayer(context.Background(), initialPlayer.ID)
		assert.NoError(t, err, "Erro ao recuperar o jogador do repositório")
		assert.Equal(t, initialPlayer.Balance, storedPlayer.Balance, "Saldo do jogador não deve ser alterado")
	})

	t.Run("CreatePlayer_InvalidPayload", func(t *testing.T) {
		invalidJSON := `{"id": "player789", "balance": "not_a_number"}`

		req, err := httpGo.NewRequest("POST", "/players", bytes.NewBufferString(invalidJSON))
		assert.NoError(t, err, "Erro ao criar a requisição HTTP")

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, httpGo.StatusBadRequest, rr.Code, "Status code deve ser 400 Bad Request")

		expectedError := HTTPError{
			Code:    httpGo.StatusBadRequest,
			Message: "Invalid request payload",
		}

		var actualError HTTPError
		err = json.Unmarshal(rr.Body.Bytes(), &actualError)
		assert.NoError(t, err, "Erro ao deserializar a resposta de erro")
		assert.Equal(t, expectedError, actualError, "Mensagem de erro deve corresponder ao esperado")
	})
}
