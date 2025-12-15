package repositories

import (
	"app-api/database"
	"app-api/internal/dtos"
	"app-api/internal/models"
	"context"
	"net/http"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NoteRepository struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewNoteRepository(logger *zap.Logger) *NoteRepository {
	return &NoteRepository{
		DB:     database.GetDB(),
		Logger: logger,
	}
}

func (r *NoteRepository) GetMyNotes(ctx context.Context, userID uint) (dtos.StructuredResponse, error) {
	var notes []models.Note

	r.Logger.Info("GetMyNotes request received", zap.Uint("userID", userID))

	// Preload User and Comments (and comment user if present)
	if err := r.DB.
		Preload("User").
		Preload("User.Role").
		Preload("Comments").
		Preload("Comments.User").
		Preload("Comments.User.Role").
		Where("user_id = ?", userID).
		Find(&notes).Error; err != nil {
		r.Logger.Error("Failed to retrieve notes", zap.Error(err))
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve notes",
			Payload: nil,
		}, err
	}

	return dtos.StructuredResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Notes retrieved successfully",
		Payload: notes,
	}, nil
}

func (r *NoteRepository) GetNoteByID(ctx context.Context, noteID uint) (dtos.StructuredResponse, error) {
	var note models.Note

	r.Logger.Info("GetNoteByID request received", zap.Uint("noteID", noteID))

	if err := r.DB.Preload("User").
		Preload("Comments").
		Preload("Comments.User").
		First(&note, noteID).Error; err != nil {
		r.Logger.Error("Failed to retrieve note", zap.Error(err))
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve note",
			Payload: nil,
		}, err
	}

	return dtos.StructuredResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Note retrieved successfully",
		Payload: note,
	}, nil
}

func (r *NoteRepository) IsMyNote(ctx context.Context, noteID uint, userID uint) (bool, error) {
	var note models.Note

	r.Logger.Info("IsMyNote request received", zap.Uint("noteID", noteID), zap.Uint("userID", userID))

	if err := r.DB.First(&note, noteID).Error; err != nil {
		r.Logger.Error("Failed to retrieve note", zap.Error(err))
		return false, err
	}

	return note.UserID == userID, nil
}

func (r *NoteRepository) GetNotes(ctx context.Context) (dtos.StructuredResponse, error) {
	var notes []models.Note

	r.Logger.Info("GetNotes request received")

	if err := r.DB.Preload("User").
		Preload("Comments").
		Preload("Comments.User").
		Find(&notes).Error; err != nil {
		r.Logger.Error("Failed to retrieve notes", zap.Error(err))
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve notes",
			Payload: nil,
		}, err
	}

	return dtos.StructuredResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Notes retrieved successfully",
		Payload: notes,
	}, nil
}

func (r *NoteRepository) CreateNote(ctx context.Context, noteDto dtos.CreateNoteDto) (dtos.StructuredResponse, error) {
	// Convert DTO to model
	note := models.Note{
		Title:       noteDto.Title,
		Description: noteDto.Description,
		Content:     noteDto.Content,
		Public:      noteDto.Public,
		UserID:      noteDto.UserID,
	}

	if err := r.DB.Create(&note).Error; err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		}, err
	}

	return dtos.StructuredResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Note created successfully",
		Payload: note,
	}, nil
}

func (r *NoteRepository) UpdateNote(ctx context.Context, noteDto dtos.UpdateNoteDto) (dtos.StructuredResponse, error) {

	var note models.Note

	if err := r.DB.First(&note, noteDto.ID).Error; err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "Note not found",
			Payload: nil,
		}, nil
	}

	note.Title = noteDto.Title
	note.Description = noteDto.Description
	note.Content = noteDto.Content
	note.IsVisited = noteDto.IsVisited

	if err := r.DB.Save(&note).Error; err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		}, err
	}

	return dtos.StructuredResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Note updated successfully",
		Payload: note,
	}, nil
}

func (r *NoteRepository) DeleteNote(ctx context.Context, noteDto dtos.DeleteNoteDto) (dtos.StructuredResponse, error) {
	var note models.Note

	if err := r.DB.First(&note, noteDto.ID).Error; err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "Note not found",
			Payload: nil,
		}, nil
	}

	if err := r.DB.Delete(&note).Error; err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		}, err
	}

	return dtos.StructuredResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Note deleted successfully",
		Payload: nil,
	}, nil
}
