package dtos

import "time"

// NoteDto represents a note in the system
// @Description A note with all its details
type NoteDto struct {
	// Unique identifier
	// @example 1
	ID uint `json:"id" example:"1"`
	// Title of the note
	// @example My First Note
	Title string `json:"title" example:"Buy groceries"`
	// Optional description with details
	// @example Milk, eggs, bread, and cheese
	Description string `json:"description" example:"Milk, eggs, bread, and cheese"`
	// Content of the note
	// @example Don't forget to check expiration dates
	Content string `json:"content" example:"Don't forget to check expiration dates"`
	// Whether the note is visited
	// @example false
	IsVisited bool `json:"isVisited" example:"false"`
	// When the note was created
	// @example 2025-06-10T10:30:00Z
	CreatedAt time.Time `json:"createdAt" example:"2025-06-10T10:30:00Z"`
	// When the note was last updated
	// @example 2025-06-10T10:30:00Z
	UpdatedAt time.Time `json:"updatedAt" example:"2025-06-10T10:30:00Z"`
}

// GetNoteDto represents the data needed to retrieve a specific note
// @Description Data for retrieving a specific note by ID
type GetNoteDto struct {
	// ID of the note to retrieve
	// @example 1
	ID uint `json:"id" example:"1"`
}

// CreateNoteDto represents the data needed to create a new note
// @Description Data for creating a new note
type CreateNoteDto struct {
	// Title of the note (3-255 characters)
	// @example My First Note
	Title string `json:"title" validate:"required;min=3;max=255" example:"My First Note"`
	// Optional description with details (max 255 characters)
	// @example This is a note about my tasks
	Description string `json:"description" validate:"max=255" example:"This is a note about my tasks"`
	// Content of the note
	// @example Don't forget to check expiration dates
	Content string `json:"content" validate:"required" example:"Don't forget to check expiration dates"`
	// Public status of the note
	// @example false
	Public bool `json:"public" example:"false"`

	// User ID associated with the note
	// @example 1
	UserID uint `json:"userId" example:"1"`
}

// UpdateNoteDto represents the data needed to update an existing note
// @Description Data for updating an existing note
type UpdateNoteDto struct {
	// ID of the note to update
	// @example 1
	ID uint `json:"id" example:"1"`
	// Updated title (3-255 characters)
	// @example My Updated Note
	Title string `json:"title" validate:"min=3;max=255" example:"My Updated Note"`
	// Updated description (max 255 characters)
	// @example This is an updated note about my tasks
	Description string `json:"description" validate:"max=255" example:"This is an updated note about my tasks"`
	// Updated content of the note
	// @example Remember to buy fruits and vegetables
	Content string `json:"content" validate:"required" example:"Remember to buy fruits and vegetables"`
	// Updated visited status
	// @example true
	IsVisited bool `json:"isVisited" example:"true"`

	// User ID associated with the note
	// @example 1
	UserID uint `json:"userId" example:"1"`
}

// DeleteNoteDto represents the data needed to delete an existing note
// @Description Data for deleting an existing note
type DeleteNoteDto struct {
	// ID of the note to delete
	// @example 1
	ID uint `json:"id" example:"1"`

	// User ID associated with the note
	// @example 1
	UserID uint `json:"userId" example:"1"`
}
