package task_interface

type Task struct {
	ID         int    `gorm:"column:id;primaryKey"`
	UserID     int    `gorm:"column:user_id;not null;index"`
	TaskNumber int    `gorm:"column:task_number;not null"` // New field
	Title      string `gorm:"column:title;not null"`
	Completed  bool   `gorm:"column:completed;default:false"`

	// Optional: Add relationship (helps with joins)
	//User user_interface.User `gorm:"foreignKey:UserID"`
}

func (Task) TableName() string {
	return "tasks"
}
