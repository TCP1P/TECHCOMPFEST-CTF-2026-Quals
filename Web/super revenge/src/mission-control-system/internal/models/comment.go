package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	NoteID    uint      `gorm:"column:note_id;not null" json:"noteId"`
	UserID    uint      `gorm:"column:user_id;not null" json:"userId"`
	Text      string    `gorm:"column:text;type:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
}

func (Comment) TableName() string {
	return "Comments"
}
