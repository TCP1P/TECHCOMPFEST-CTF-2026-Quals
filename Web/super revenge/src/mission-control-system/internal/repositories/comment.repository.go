package repositories

import (
	"app-api/database"
	"app-api/internal/dtos"
	"app-api/internal/models"
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentRepository struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewCommentRepository(logger *zap.Logger) *CommentRepository {
	return &CommentRepository{
		DB:     database.GetDB(),
		Logger: logger,
	}
}

func (r *CommentRepository) GetCommentsByNote(ctx context.Context, noteID uint) (dtos.StructuredResponse, error) {
	var comments []models.Comment

	r.Logger.Info("GetCommentsByNote request received", zap.Uint("noteID", noteID))

	if err := r.DB.WithContext(ctx).Where("note_id = ?", noteID).Find(&comments).Error; err != nil {
		r.Logger.Error("Failed to retrieve comments", zap.Error(err))
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve comments",
			Payload: nil,
		}, err
	}

	return dtos.StructuredResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Comments retrieved successfully",
		Payload: comments,
	}, nil
}

func (r *CommentRepository) CreateComment(ctx context.Context, dto dtos.CreateCommentDto) (dtos.StructuredResponse, error) {
	comment := models.Comment{
		NoteID: dto.NoteID,
		UserID: dto.UserID,
		Text:   dto.Content,
	}

	if err := r.DB.WithContext(ctx).Create(&comment).Error; err != nil {
		r.Logger.Error("Failed to create comment", zap.Error(err))
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
		Message: "Comment created successfully",
		Payload: comment,
	}, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, dto dtos.UpdateCommentDto) (dtos.StructuredResponse, error) {
	var comment models.Comment

	if err := r.DB.WithContext(ctx).First(&comment, dto.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dtos.StructuredResponse{
				Success: false,
				Status:  http.StatusNotFound,
				Message: "Comment not found",
				Payload: nil,
			}, nil
		}
		r.Logger.Error("Failed to fetch comment", zap.Error(err))
		return dtos.StructuredResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Payload: nil,
		}, err
	}

	comment.Text = dto.Content

	if err := r.DB.WithContext(ctx).Save(&comment).Error; err != nil {
		r.Logger.Error("Failed to update comment", zap.Error(err))
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
		Message: "Comment updated successfully",
		Payload: comment,
	}, nil
}

func (r *CommentRepository) IsMyComment(ctx context.Context, commentID uint, userID uint) (bool, error) {
	var comment models.Comment

	if err := r.DB.WithContext(ctx).First(&comment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		r.Logger.Error("Failed to fetch comment", zap.Error(err))
		return false, err
	}

	isOwner := comment.UserID == userID
	return isOwner, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, dto dtos.DeleteCommentDto) (dtos.StructuredResponse, error) {
	if err := r.DB.WithContext(ctx).Delete(&models.Comment{}, dto.ID).Error; err != nil {
		r.Logger.Error("Failed to delete comment", zap.Error(err))
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
		Message: "Comment deleted successfully",
		Payload: nil,
	}, nil
}
