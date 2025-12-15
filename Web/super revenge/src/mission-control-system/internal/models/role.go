package models

var (
	RoleAdmin   = "admin"
	RoleMember  = "member"
	DefaultRole = RoleMember
)

type Role struct {
	ID   uint   `gorm:"primaryKey;column:id" json:"id"`
	Name string `gorm:"size:50;not null;unique;column:name" json:"name"`
}

func (Role) TableName() string {
	return "Roles"
}
