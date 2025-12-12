package user_request

import (
	"ToDoList/interfaces/task_interface"
	"ToDoList/interfaces/user_interface"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateUserTask(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	username := c.Param("username")
	taskNumberStr := c.Param("id") // This is now task_number, not database ID

	taskNumber, err := strconv.Atoi(taskNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task number"})
		return
	}

	// Find the user
	var user user_interface.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Find the task by user_id and task_number
	var task task_interface.Task
	result := db.Where("user_id = ? AND task_number = ?", user.ID, taskNumber).First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Update fields
	if input.Title != nil {
		if *input.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
			return
		}
		task.Title = *input.Title
	}

	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	if err := db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"task":    task,
	})
}
