package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(handler *Handler) http.Handler {
	r := mux.NewRouter()

	// Rotas para jogadores
	r.HandleFunc("/players", handler.CreatePlayer).Methods("POST")
	r.HandleFunc("/players/balance", handler.GetPlayerBalance).Methods("GET")

	// Rotas para máquinas de slot
	r.HandleFunc("/machines", handler.CreateSlotMachine).Methods("POST")
	r.HandleFunc("/machines/balance", handler.GetSlotMachineBalance).Methods("GET")

	// Rota para jogar na máquina de slot
	r.HandleFunc("/play", handler.PlaySlotMachine).Methods("POST")

	return r
}
