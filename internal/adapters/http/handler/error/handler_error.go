package handler_error

import (
	"encoding/json"
	"net/http"
	"slot-machine/internal/application/usecase"
	"slot-machine/internal/domain/repository"
)

func HandleError(w http.ResponseWriter, err error) {
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

// HTTPError representa um erro retornado pela API.
// @Description Estrutura para representar erros na API.
// @Description Contém a mensagem de erro e um código opcional.
// @Description Pode ser expandida conforme necessário.
type HTTPError struct {
	Code    int    `json:"code"`    // Código do erro HTTP
	Message string `json:"message"` // Mensagem descritiva do erro
}
