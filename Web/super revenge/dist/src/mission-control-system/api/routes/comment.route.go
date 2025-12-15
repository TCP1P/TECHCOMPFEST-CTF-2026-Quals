package routes

import (
	"app-api/api/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func HandleCommentRoutes(api *mux.Router, logger *zap.Logger) {
	logger.Info("Comment routes registered")

	commentHandler := handlers.NewCommentHandler(logger)

	protectedRouter := ApplyAuthMiddleware(api, logger)
	protectedRouter.HandleFunc("/note/{noteId}", commentHandler.GetCommentsByNote).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/new", commentHandler.CreateComment).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/edit", commentHandler.UpdateComment).Methods(http.MethodPut)
	protectedRouter.HandleFunc("/delete", commentHandler.DeleteComment).Methods(http.MethodDelete)
}
