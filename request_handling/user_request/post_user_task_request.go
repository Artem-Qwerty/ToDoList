package user_request

import (
	"ToDoList/interfaces/task_interface"
	"ToDoList/interfaces/user_interface"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUserTask(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	username := c.Param("username")

	// Find the user
	var user user_interface.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}

	// Get the highest task_number for this user
	var maxTaskNumber int
	db.Model(&task_interface.Task{}).
		Where("user_id = ?", user.ID).
		Select("COALESCE(MAX(task_number), 0)").
		Scan(&maxTaskNumber)

	// Create new task with next task_number
	newTask := task_interface.Task{
		UserID:     user.ID,
		TaskNumber: maxTaskNumber + 1,
		Title:      input.Title,
		Completed:  input.Completed,
	}

	result := db.Create(&newTask)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task":    newTask,
	})
}
