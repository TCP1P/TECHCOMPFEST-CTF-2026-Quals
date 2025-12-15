package models

type User struct {
	ID           uint   `gorm:"primaryKey;column:id" json:"id"`
	Email        string `gorm:"column:email;not null;unique" json:"email"`
	Name         string `gorm:"column:name;not null" json:"name"`
	PasswordHash string `gorm:"column:passwordHash;not null" json:"-"` // Using json:"-" to exclude from JSON responses
	RoleID       uint   `gorm:"column:role_id;not null;default:1" json:"roleId"`
	Role         Role   `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	Notes        []Note `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"notes"`
}

func (User) TableName() string {
	return "Users"
}
