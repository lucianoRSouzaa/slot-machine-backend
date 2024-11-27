package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	httpInternal "slot-machine/internal/adapters/http"
	"slot-machine/internal/application/usecase"
	repository_in_memory "slot-machine/internal/infrastructure/repository/in_memory"
	"slot-machine/internal/infrastructure/security"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	playerRepo := repository_in_memory.NewInMemoryPlayerRepository()
	slotRepo := repository_in_memory.NewInMemorySlotMachineRepository()

	hasher := security.NewBcryptPasswordHasher(bcrypt.DefaultCost)

	playUC := usecase.NewPlayUseCase(playerRepo, slotRepo)
	createPlayerUC := usecase.NewCreatePlayerUseCase(playerRepo, hasher)
	createSlotMachineUC := usecase.NewCreateSlotMachineUseCase(slotRepo)
	getPlayerBalanceUC := usecase.NewGetPlayerBalanceUseCase(playerRepo)
	getSlotMachineBalanceUC := usecase.NewGetSlotMachineBalanceUseCase(slotRepo)

	handler := httpInternal.NewHandler(createPlayerUC, createSlotMachineUC, playUC, getPlayerBalanceUC, getSlotMachineBalanceUC)

	router := httpInternal.NewRouter(handler)

	corsAllowedOrigins := []string{"http://localhost:5173"}
	corsAllowedMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsAllowedHeaders := []string{"Content-Type", "Authorization"}

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins(corsAllowedOrigins),
		handlers.AllowedMethods(corsAllowedMethods),
		handlers.AllowedHeaders(corsAllowedHeaders),
	)

	loggingMiddleware := handlers.LoggingHandler(os.Stdout, router)

	finalHandler := corsMiddleware(loggingMiddleware)

	addr := ":8081"

	server := &http.Server{
		Addr:         addr,
		Handler:      finalHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para capturar erros do servidor
	serverErrors := make(chan error, 1)

	// Inicia o servidor em uma goroutine
	go func() {
		logger.Infof("Servidor rodando em %s", addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Canal para capturar sinais do sistema para shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Aguarda por erros ou sinais de shutdown
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	case sig := <-sigChan:
		logger.Infof("Recebido sinal %v, iniciando shutdown", sig)

		// Cria um contexto com timeout para o shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Tenta realizar o shutdown gracefull
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Erro durante shutdown: %v", err)
		}

		logger.Info("Servidor finalizado com sucesso")
	}
}
