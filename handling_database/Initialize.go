package handling_database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;size:100;not null"`
	Password string `gorm:"size:255;not null"`
	Tasks    []Task `gorm:"constraint:OnDelete:CASCADE"`
}

type Task struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null;index"`
	TaskNumber int    `gorm:"not null"`
	Title      string `gorm:"size:255;not null"`
	Completed  bool   `gorm:"default:false"`

	// Unique(user_id, task_number)
	_ struct{} `gorm:"uniqueIndex:idx_user_tasknumber"`
}

type Quote struct {
	ID      int    `gorm:"column:id;primaryKey"`
	Text    string `gorm:"column:text"`
	Current bool   `gorm:"column:current"`
}

func AutoMigrateAndSeed(db *gorm.DB) {
	// Run migrations
	if err := db.AutoMigrate(&User{}, &Task{}, &Quote{}); err != nil {
		log.Fatal("Migration failed:", err)
	}
	hash1, err := bcrypt.GenerateFromPassword([]byte("ilovemary"), 14)
	if err != nil {
		log.Fatal(err)
	}
	hash2, err := bcrypt.GenerateFromPassword([]byte("iamdead"), 14)
	if err != nil {
		log.Fatal(err)
	}
	password1 := string(hash1)
	password2 := string(hash2)
	// Insert users (ON CONFLICT DO NOTHING equivalent)
	users := []User{
		{
			Username: "James",
			Password: password1,
		},
		{
			Username: "Mary",
			Password: password2,
		},
	}

	for _, u := range users {
		db.Where(User{Username: u.Username}).FirstOrCreate(&u)
	}

	// Insert tasks
	tasks := []Task{
		{UserID: 1, TaskNumber: 1, Title: "Escape Silent Hill", Completed: false},
		{UserID: 1, TaskNumber: 2, Title: "Find Mary", Completed: false},
		{UserID: 1, TaskNumber: 3, Title: "Headshot a monster", Completed: true},

		{UserID: 2, TaskNumber: 1, Title: "Give a proper location", Completed: false},
		{UserID: 2, TaskNumber: 2, Title: "Die", Completed: true},
	}

	for _, t := range tasks {
		db.Where(Task{UserID: t.UserID, TaskNumber: t.TaskNumber}).FirstOrCreate(&t)
	}
	quotes := []Quote{
		{ID: 1, Text: "MyQuote1", Current: false},
		{ID: 2, Text: "MyQuote2", Current: true},
		{ID: 3, Text: "MyQuote3", Current: false},
	}
	for _, q := range quotes {
		db.Where(Quote{ID: q.ID, Text: q.Text, Current: q.Current}).FirstOrCreate(&q)
	}

	// Verification (optional logging)
	var userCount, taskCount, quoteCount int64
	db.Model(&User{}).Count(&userCount)
	db.Model(&Task{}).Count(&taskCount)
	db.Model(&Quote{}).Count(&quoteCount)
	log.Printf("Users created: %d\n", userCount)
	log.Printf("Tasks created: %d\n", taskCount)
	log.Printf("Quotes created: %d\n", quoteCount)
}
