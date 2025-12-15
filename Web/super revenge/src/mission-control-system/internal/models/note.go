package models

import "time"

type Note struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	Title       string    `gorm:"size:255;not null;column:title" json:"title"`
	Description string    `gorm:"size:255;null;column:description" json:"description"`
	Content     string    `gorm:"type:text;not null;column:content" json:"content"`
	IsVisited   bool      `gorm:"default:false;column:isVisited" json:"isVisited"`
	Public      bool      `gorm:"default:false;column:public" json:"public"`
	Comments    []Comment `gorm:"foreignKey:NoteID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
	CreatedAt   time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	UserID      uint      `gorm:"column:user_id;not null" json:"userId"`
	User        User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
}

// TableName overrides the table name used by TodoItem to `todos`
func (Note) TableName() string {
	return "Notes"
}
