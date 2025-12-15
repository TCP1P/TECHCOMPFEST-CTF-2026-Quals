package services

import (
	"app-api/internal/dtos"
	"app-api/internal/repositories"
	"context"

	"go.uber.org/zap"
)

type CommentService struct {
	repo *repositories.CommentRepository
}

func NewCommentService(logger *zap.Logger) *CommentService {
	return &CommentService{
		repo: repositories.NewCommentRepository(logger),
	}
}

func (s *CommentService) GetCommentsByNote(ctx context.Context, noteID uint) (dtos.StructuredResponse, error) {
	return s.repo.GetCommentsByNote(ctx, noteID)
}

func (s *CommentService) CreateComment(ctx context.Context, dto dtos.CreateCommentDto) (dtos.StructuredResponse, error) {
	return s.repo.CreateComment(ctx, dto)
}

func (s *CommentService) UpdateComment(ctx context.Context, dto dtos.UpdateCommentDto) (dtos.StructuredResponse, error) {
	is_valid, err := s.repo.IsMyComment(ctx, dto.ID, dto.UserID)
	if err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  500,
			Message: "Failed to verify comment ownership",
			Payload: nil,
		}, err
	}
	if !is_valid {
		return dtos.StructuredResponse{
			Success: false,
			Status:  403,
			Message: "Unauthorized to update this comment",
			Payload: nil,
		}, nil
	}
	return s.repo.UpdateComment(ctx, dto)
}

func (s *CommentService) DeleteComment(ctx context.Context, dto dtos.DeleteCommentDto) (dtos.StructuredResponse, error) {
	is_valid, err := s.repo.IsMyComment(ctx, dto.ID, dto.UserID)
	if err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  500,
			Message: "Failed to verify comment ownership",
			Payload: nil,
		}, err
	}
	if !is_valid {
		return dtos.StructuredResponse{
			Success: false,
			Status:  403,
			Message: "Unauthorized to update this comment",
			Payload: nil,
		}, nil
	}
	return s.repo.DeleteComment(ctx, dto)
}
