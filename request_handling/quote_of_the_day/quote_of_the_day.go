package quote_of_the_day

import (
	"ToDoList/interfaces/quote_interface"
	"ToDoList/internal/cache"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTheQuote(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	redis := c.MustGet("redis").(*cache.RedisClient)

	// 1. Пробуем Redis
	val, err := redis.Client.Get(cache.Ctx, "current_quote").Result()
	if err == nil {
		c.Data(http.StatusOK, "application/json", []byte(val))
		return
	}

	// 2. Redis пуст — берем из БД
	var quote quote_interface.Quote
	if err := db.Where("current = true").First(&quote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No quotes available"})
		return
	}

	// 3. Кладем строку
	jsonData, _ := json.Marshal(quote)
	redis.Client.Set(cache.Ctx, "current_quote", string(jsonData), 10*time.Second)

	c.JSON(http.StatusOK, quote)
}
func StartQuoteRotator(db *gorm.DB, redis *cache.RedisClient) {
	go func() {
		for {
			time.Sleep(10 * time.Second)

			var current quote_interface.Quote
			if err := db.Where("current = true").First(&current).Error; err != nil {
				continue
			}

			// ищем следующую
			var next quote_interface.Quote
			if err := db.Where("id > ?", current.ID).First(&next).Error; err != nil {
				db.First(&next)
			}

			// переключаем
			db.Model(&quote_interface.Quote{}).Where("id = ?", current.ID).Update("current", false)
			db.Model(&quote_interface.Quote{}).Where("id = ?", next.ID).Update("current", true)

			// перезагружаем next уже с current=true
			db.Where("id = ?", next.ID).First(&next)

			// кладём в redis
			jsonData, _ := json.Marshal(next)
			redis.Client.Set(cache.Ctx, "current_quote", string(jsonData), 10*time.Second)
		}
	}()
}
