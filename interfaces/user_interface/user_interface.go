package user_interface

type User struct {
	ID       int    `gorm:"column:id;primaryKey"`
	Username string `gorm:"column:username;unique;not null"`
	Password string `gorm:"column:password;not null"` // Store hashed password
}

func (User) TableName() string {
	return "users"
}
