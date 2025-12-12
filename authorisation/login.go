package authorisation

import (
	"ToDoList/interfaces/user_interface"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("=== LOGIN DEBUG ===")
	fmt.Println("Looking for username:", input.Username)

	// Find user
	var user user_interface.User
	result := db.Where("username = ?", input.Username).First(&user)

	fmt.Println("Query error:", result.Error)
	fmt.Println("Rows affected:", result.RowsAffected)
	fmt.Println("Found user ID:", user.ID)
	fmt.Println("Found username:", user.Username)
	fmt.Println("Found password hash:", user.Password)
	fmt.Println("==================")

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": user.Username,
		"user_id":  user.ID,
	})
}
