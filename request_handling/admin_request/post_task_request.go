package admin_request

import (
	"ToDoList/interfaces/task_interface"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostTask(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var newTask task_interface.Task

	// Bind JSON from request body
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Insert into database
	result := db.Create(&newTask)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task: " + result.Error.Error()})
		return
	}

	// Return the created task
	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task":    newTask,
	})
}
