package routes

import (
	"net/http"

	"app-api/api/handlers"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func HandleAuthRoutes(api *mux.Router, logger *zap.Logger) {

	authHandler := handlers.NewAuthHandler(logger)

	api.HandleFunc("/register", authHandler.RegisterUser).Methods(http.MethodPost)
	api.HandleFunc("/login", authHandler.LoginUser).Methods(http.MethodPost)
	api.HandleFunc("/check-token", authHandler.CheckToken).Methods(http.MethodPost)
}
