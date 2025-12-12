package registration

import (
	"ToDoList/authorisation"
	"ToDoList/interfaces/user_interface"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := authorisation.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := user_interface.User{
		Username: input.Username,
		Password: hashedPassword,
	}

	result := db.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists or invalid"})
		return
	}

	// Generate token for immediate login
	token, err := authorisation.GenerateToken(newUser.ID, newUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User created but failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "User registered successfully",
		"username": newUser.Username,
		"user_id":  newUser.ID,
		"token":    token,
	})
}
