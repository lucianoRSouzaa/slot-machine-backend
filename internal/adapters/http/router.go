package http

import (
	"net/http"

	_ "slot-machine/docs"

	"slot-machine/internal/adapters/http/handler"
	"slot-machine/internal/adapters/http/middleware"
	"slot-machine/internal/domain/ports"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handler *handler.Handler, jwtManager ports.JWTManager) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.HandleFunc("/refresh", handler.Refresh).Methods("POST")
	r.HandleFunc("/players", handler.CreatePlayer).Methods("POST")

	secure := r.PathPrefix("/").Subrouter()
	secure.Use(middleware.JWTMiddleware(jwtManager))

	secure.HandleFunc("/players/balance", handler.GetPlayerBalance).Methods("GET")
	secure.HandleFunc("/play", handler.PlaySlotMachine).Methods("POST")

	admin := r.PathPrefix("/").Subrouter()
	admin.Use(middleware.AdminMiddleware(jwtManager))

	admin.HandleFunc("/machines", handler.CreateSlotMachine).Methods("POST")
	admin.HandleFunc("/machines/balance", handler.GetSlotMachineBalance).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
