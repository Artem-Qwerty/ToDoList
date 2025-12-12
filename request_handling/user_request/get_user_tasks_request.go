package user_request

import (
	"ToDoList/interfaces/task_interface"
	"ToDoList/interfaces/user_interface"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserTasks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	username := c.Param("username")

	// Find the user
	var user user_interface.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get tasks ordered by task_number
	var tasks []task_interface.Task
	result := db.Where("user_id = ?", user.ID).Order("task_number ASC").Find(&tasks)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"tasks":    tasks,
		"count":    len(tasks),
	})
}
