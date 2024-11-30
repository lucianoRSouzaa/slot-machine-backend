package http

import (
	"net/http"

	_ "slot-machine/docs"

	"slot-machine/internal/adapters/http/handler"
	"slot-machine/internal/adapters/middleware"
	"slot-machine/internal/domain/ports"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handler *handler.Handler, jwtManager ports.JWTManager) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/players", handler.CreatePlayer).Methods("POST")

	secure := r.PathPrefix("/").Subrouter()
	secure.Use(middleware.JWTMiddleware(jwtManager))

	secure.HandleFunc("/players/balance", handler.GetPlayerBalance).Methods("GET")
	secure.HandleFunc("/machines", handler.CreateSlotMachine).Methods("POST")
	secure.HandleFunc("/machines/balance", handler.GetSlotMachineBalance).Methods("GET")
	secure.HandleFunc("/play", handler.PlaySlotMachine).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
