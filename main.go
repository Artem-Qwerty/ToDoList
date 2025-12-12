package main

import (
	"ToDoList/authorisation"
	"ToDoList/handling_database"
	"ToDoList/internal/cache"
	"ToDoList/registration"
	"ToDoList/request_handling/quote_of_the_day"
	"ToDoList/request_handling/user_request"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"

	_ "net/http"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("dsn")

	db := handling_database.ConnectDB(dsn)
	handling_database.AutoMigrateAndSeed(db)
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		panic("Cannot connect to Postgres")
	}
	fmt.Println("Connected to PostgreSQL successfully!")

	redisClient := cache.NewRedisClient("localhost:6379")

	r := gin.Default()

	// Inject DB into context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("redis", redisClient)
		c.Next()
	})
	quote_of_the_day.StartQuoteRotator(db, redisClient)
	// Public routes
	r.GET("/quote", quote_of_the_day.GetTheQuote)
	r.POST("/users/register", registration.RegisterUser)
	r.POST("/users/login", authorisation.LoginUser)

	// Protected routes - require authentication
	protected := r.Group("/tasks")
	protected.Use(authorisation.AuthMiddleware())
	{

		protected.GET("/:username", authorisation.ValidateUsernameMatch(), user_request.GetUserTasks)
		protected.POST("/:username", authorisation.ValidateUsernameMatch(), user_request.CreateUserTask)
		protected.PATCH("/:username/:id", authorisation.ValidateUsernameMatch(), user_request.UpdateUserTask)
	}

	r.Run(":8085")
}
