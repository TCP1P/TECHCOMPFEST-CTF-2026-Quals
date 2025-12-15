package handlers

import (
	"app-api/internal/dtos"
	"app-api/internal/services"
	"app-api/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type NoteHandler struct {
	BaseHandler
	service *services.NoteService
}

func NewNoteHandler(logger *zap.Logger) *NoteHandler {
	return &NoteHandler{
		BaseHandler: BaseHandler{
			Logger: logger,
		},
		service: services.NewNoteService(logger),
	}
}

// @Summary Get current user's Notes
// @Description Get all Notes created by the authenticated user
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dtos.StructuredResponse "Notes retrieved successfully"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /notes/get-my-notes [get]
func (h *NoteHandler) GetMyNotes(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("GetMyNotes request received")

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		h.Logger.Error("Failed to get user ID from context", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	h.Logger.Debug("Fetching notes for user", zap.Uint("userID", userID))
	response, err := h.service.GetMyNotes(r.Context(), userID)
	if err != nil {
		h.Logger.Error("Failed to get user notes", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	h.Logger.Info("User notes retrieved successfully")
	h.ReturnJSONResponse(w, response)
}

// @Summary Get Note by ID
// @Description Get a Note by its ID, if it is public or belongs to the user
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param noteId path int true "Note ID"
// @Success 200 {object} dtos.StructuredResponse "Note retrieved successfully"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 403 {object} dtos.StructuredResponse "Forbidden"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /notes/get-note-by-id/{noteId} [get]
func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	noteIDParam := mux.Vars(r)["noteId"]

	noteIDUint64, err := strconv.ParseUint(noteIDParam, 10, 64)
	if err != nil {
		h.Logger.Warn("Invalid note ID", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid note ID",
			Payload: nil,
		})
		return
	}

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		h.Logger.Error("Failed to get user ID from context", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	response, err := h.service.GetNoteByID(r.Context(), uint(noteIDUint64), userID)
	if err != nil {
		h.Logger.Error("Failed to fetch note", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	h.Logger.Info("Note retrieved successfully")
	h.ReturnJSONResponse(w, response)
}

// @Summary Get all Notes ( admin only, not ready yet )
// @Description Get all Notes from the database
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dtos.StructuredResponse "Notes retrieved successfully"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /notes/get-notes [get]
func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("GetNotes request received")

	// userId, err := utils.GetUserIDFromContext(r.Context())
	// if err != nil {
	// 	h.Logger.Error("Failed to get user ID from context", zap.Error(err))
	// 	h.ReturnJSONResponse(w, dtos.StructuredResponse{
	// 		Success: false,
	// 		Status:  http.StatusInternalServerError,
	// 		Message: err.Error(),
	// 		Payload: nil,
	// 	})
	// 	return
	// }

	// if !utils.IsUserAdmin(r.Context()) {
	// 	h.Logger.Warn("Unauthorized access to GetNotes", zap.Uint("userID", userId))
	// 	h.ReturnJSONResponse(w, dtos.StructuredResponse{
	// 		Success: false,
	// 		Status:  http.StatusForbidden,
	// 		Message: "Unauthorized access",
	// 		Payload: nil,
	// 	})
	// 	return
	// }

	response, err := h.service.GetNotes(r.Context())

	if err != nil {
		h.Logger.Error("Failed to get notes", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	// The repository now guarantees that response. Payload will be a slice (possibly empty)
	h.Logger.Info("Notes retrieved successfully")

	// Return the actual response, not a hardcoded message
	h.ReturnJSONResponse(w, response)
}

// @Summary Create a new Note
// @Description Create a new Note with the provided details
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param note body dtos.CreateNoteDto true "Note data"
// @Success 200 {object} dtos.StructuredResponse "Note created successfully"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /notes/create-note [post]
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("CreateNote request received")

	var req dtos.CreateNoteDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Failed to decode request body", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}
	defer r.Body.Close()

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		h.Logger.Error("Failed to get user ID from context", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	req.UserID = userID

	h.Logger.Debug("Creating note", zap.String("title", req.Title))
	response, err := h.service.CreateNote(r.Context(), req)

	if err != nil {
		h.Logger.Error("Failed to create note", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	h.Logger.Info("Note created successfully")
	h.ReturnJSONResponse(w, response)
}

// @Summary Update an existing Note
// @Description Update an existing Note with the provided details
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param note body dtos.UpdateNoteDto true "Note update data"
// @Success 200 {object} dtos.StructuredResponse "Note updated successfully"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /notes/update-note [put]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("UpdateNote request received")

	var req dtos.UpdateNoteDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Failed to decode request body", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}
	defer r.Body.Close()

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		h.Logger.Error("Failed to get user ID from context", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	req.UserID = userID
	h.Logger.Debug("Updating note", zap.Uint("id", req.ID))
	response, err := h.service.UpdateNote(r.Context(), req)

	if err != nil {
		h.Logger.Error("Failed to update note", zap.Error(err))
		if response.Message != "" {
			h.ReturnJSONResponse(w, response)
			return
		}
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	h.ReturnJSONResponse(w, response)
}

// @Summary Delete an existing Note
// @Description Delete an existing Note by ID
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param note body dtos.DeleteNoteDto true "Note deletion data"
// @Success 200 {object} dtos.StructuredResponse "Note deleted successfully"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /notes/delete-note [delete]
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("DeleteNote request received")

	var req dtos.DeleteNoteDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("Failed to decode request body", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}
	defer r.Body.Close()

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		h.Logger.Error("Failed to get user ID from context", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	req.UserID = userID
	h.Logger.Debug("Deleting note", zap.Uint("id", req.ID))
	response, err := h.service.DeleteNote(r.Context(), req)

	if err != nil {

		h.Logger.Error("Failed to delete note", zap.Error(err))
		if response.Message != "" {
			h.ReturnJSONResponse(w, response)
			return
		}
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}

	h.Logger.Info("Note deleted successfully")
	h.ReturnJSONResponse(w, response)
}
