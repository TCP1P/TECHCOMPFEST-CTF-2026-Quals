package dtos

type CreateCommentDto struct {
	Content string `json:"content"`
	NoteID  uint   `json:"noteId"`
	UserID  uint   `json:"-"`
}

type UpdateCommentDto struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`

	UserID uint `json:"-"`
}

type DeleteCommentDto struct {
	ID uint `json:"id"`

	UserID uint `json:"-"`
}
