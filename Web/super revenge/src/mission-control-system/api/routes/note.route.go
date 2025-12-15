package routes

import (
	"app-api/api/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func HandleNoteRoutes(api *mux.Router, logger *zap.Logger) {
	// Debug route to confirm this handler is being registered
	logger.Info("Note routes registered")

	// Root note endpoint
	noteHandler := handlers.NewNoteHandler(logger)

	// Public routes (if any)
	// No public routes for now

	// Protected routes (require authentication)
	protectedRouter := ApplyAuthMiddleware(api, logger)
	adminRouter := ApplyAdminMiddleware(api, logger)

	adminRouter.HandleFunc("/get-notes", noteHandler.GetNotes).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/get-my-notes", noteHandler.GetMyNotes).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/get-note-by-id/{noteId}", noteHandler.GetNoteByID).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/create-note", noteHandler.CreateNote).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/update-note", noteHandler.UpdateNote).Methods(http.MethodPut)
	protectedRouter.HandleFunc("/delete-note", noteHandler.DeleteNote).Methods(http.MethodDelete)
}
