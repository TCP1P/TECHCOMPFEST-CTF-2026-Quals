package services

import (
	"app-api/internal/dtos"
	"app-api/internal/models"
	"app-api/internal/repositories"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type NoteService struct {
	noteRepository *repositories.NoteRepository
}

func NewNoteService(logger *zap.Logger) *NoteService {
	return &NoteService{
		noteRepository: repositories.NewNoteRepository(logger),
	}
}

func (s *NoteService) GetNotes(ctx context.Context) (dtos.StructuredResponse, error) {

	response, err := s.noteRepository.GetNotes(ctx)
	if err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  500,
			Message: "Failed to retrieve notes",
			Payload: nil,
		}, err
	}
	return response, nil
}

func (s *NoteService) GetNoteByID(ctx context.Context, noteID uint, userID uint) (dtos.StructuredResponse, error) {

	notes, err := s.noteRepository.GetNoteByID(ctx, noteID)
	if err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  500,
			Message: "Failed to retrieve note",
			Payload: nil,
		}, err
	}
	isMine, err := s.noteRepository.IsMyNote(ctx, noteID, userID)
	if err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  500,
			Message: "Failed to verify note ownership",
			Payload: nil,
		}, fmt.Errorf("failed to verify note ownership")
	}
	if !notes.Payload.(models.Note).Public && !isMine {
		return dtos.StructuredResponse{
			Success: false,
			Status:  403,
			Message: "Unauthorized to access this note",
			Payload: nil,
		}, fmt.Errorf("unauthorized access to note")
	}
	return notes, nil
}

func (s *NoteService) CreateNote(ctx context.Context, note dtos.CreateNoteDto) (dtos.StructuredResponse, error) {
	return s.noteRepository.CreateNote(ctx, note)
}

func (s *NoteService) GetMyNotes(ctx context.Context, userID uint) (dtos.StructuredResponse, error) {
	return s.noteRepository.GetMyNotes(ctx, userID)
}

func (s *NoteService) UpdateNote(ctx context.Context, noteDto dtos.UpdateNoteDto) (dtos.StructuredResponse, error) {
	is_valid, err := s.noteRepository.IsMyNote(ctx, noteDto.ID, noteDto.UserID)
	if err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  500,
			Message: "Failed to verify note ownership",
			Payload: nil,
		}, fmt.Errorf("failed to verify note ownership")
	}
	if !is_valid {
		return dtos.StructuredResponse{
			Success: false,
			Status:  403,
			Message: "Unauthorized to update this note",
			Payload: nil,
		}, fmt.Errorf("unauthorized userto update this note")
	}
	return s.noteRepository.UpdateNote(ctx, noteDto)
}

func (s *NoteService) DeleteNote(ctx context.Context, noteDto dtos.DeleteNoteDto) (dtos.StructuredResponse, error) {
	is_valid, err := s.noteRepository.IsMyNote(ctx, noteDto.ID, noteDto.UserID)
	if err != nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  500,
			Message: "Failed to verify note ownership",
			Payload: nil,
		}, fmt.Errorf("failed to verify note ownership")
	}
	if !is_valid {
		return dtos.StructuredResponse{
			Success: false,
			Status:  403,
			Message: "Unauthorized to delete this note",
			Payload: nil,
		}, fmt.Errorf("unauthorized user to delete this note")
	}
	return s.noteRepository.DeleteNote(ctx, noteDto)
}
