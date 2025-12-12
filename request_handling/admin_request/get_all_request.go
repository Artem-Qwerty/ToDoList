package admin_request

import (
	"ToDoList/interfaces/task_interface"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllTasks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var tasks []task_interface.Task
	result := db.Find(&tasks)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}
