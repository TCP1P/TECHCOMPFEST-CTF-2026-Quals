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

type CommentHandler struct {
	BaseHandler
	service *services.CommentService
}

func NewCommentHandler(logger *zap.Logger) *CommentHandler {
	return &CommentHandler{
		BaseHandler: BaseHandler{Logger: logger},
		service:     services.NewCommentService(logger),
	}
}

// @Summary Get comments for a note
// @Description Retrieve all comments belonging to the specified note
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param noteId path int true "Note ID"
// @Success 200 {object} dtos.StructuredResponse "Comments retrieved successfully"
// @Failure 400 {object} dtos.StructuredResponse "Invalid note ID"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /comments/note/{noteId} [get]
func (h *CommentHandler) GetCommentsByNote(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.service.GetCommentsByNote(r.Context(), uint(noteIDUint64))
	if err != nil {
		h.Logger.Error("Failed to fetch comments", zap.Error(err))
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

// @Summary Create a comment
// @Description Create a new comment for a note
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment body dtos.CreateCommentDto true "Comment data"
// @Success 200 {object} dtos.StructuredResponse "Comment created successfully"
// @Failure 400 {object} dtos.StructuredResponse "Invalid payload"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /comments/new [post]
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateCommentDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("Failed to decode comment payload", zap.Error(err))
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
		h.Logger.Error("Failed to get user ID", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}
	req.UserID = userID

	response, err := h.service.CreateComment(r.Context(), req)
	if err != nil {
		h.Logger.Error("Failed to create comment", zap.Error(err))
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

// @Summary Update a comment
// @Description Update an existing comment
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment body dtos.UpdateCommentDto true "Comment update data"
// @Success 200 {object} dtos.StructuredResponse "Comment updated successfully"
// @Failure 400 {object} dtos.StructuredResponse "Invalid payload"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 404 {object} dtos.StructuredResponse "Comment not found"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /comments/edit [put]
func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var req dtos.UpdateCommentDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("Failed to decode comment payload", zap.Error(err))
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
		h.Logger.Error("Failed to get user ID", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}
	req.UserID = userID

	response, err := h.service.UpdateComment(r.Context(), req)
	if err != nil {
		h.Logger.Error("Failed to update comment", zap.Error(err))
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

// @Summary Delete a comment
// @Description Delete an existing comment
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment body dtos.DeleteCommentDto true "Comment deletion data"
// @Success 200 {object} dtos.StructuredResponse "Comment deleted successfully"
// @Failure 400 {object} dtos.StructuredResponse "Invalid payload"
// @Failure 401 {object} dtos.StructuredResponse "Unauthorized"
// @Failure 500 {object} dtos.StructuredResponse "Internal server error"
// @Router /comments/delete [delete]
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	var req dtos.DeleteCommentDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("Failed to decode comment payload", zap.Error(err))
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
		h.Logger.Error("Failed to get user ID", zap.Error(err))
		h.ReturnJSONResponse(w, dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		})
		return
	}
	req.UserID = userID

	response, err := h.service.DeleteComment(r.Context(), req)
	if err != nil {
		h.Logger.Error("Failed to delete comment", zap.Error(err))
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
