package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var db = make(map[string]string)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": "OK"})
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": "OK"})
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	type MessageFrom struct {
		ID           int32  `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	}

	type MessageChat struct {
		ID        int32  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	}

	type MessageEntity struct {
		Type   string `json:"type"`
		Offset int32  `json:"offset"`
		Length int32  `json:"length"`
	}

	type Message struct {
		MessageID int32           `json:"message_id"`
		From      MessageFrom     `json:"from"`
		Chat      MessageChat     `json:"chat"`
		Date      int32           `json:"date"`
		Text      string          `json:"text"`
		Entities  []MessageEntity `json:"entities"`
	}

	type Update struct {
		UpdateID int32   `json:"update_id"`
		Message  Message `json:"message"`
	}

	r.POST("/messages", func(c *gin.Context) {
		var update Update
		err := c.BindJSON(&update)
		failOnError(err, "Failed to unmarshal data")
		sendMessage(update.Message.Text, update.Message.Chat.ID)
		c.JSON(http.StatusAccepted, gin.H{
			"message": update.Message.Text,
		})
	})

	return r
}

func sendMessage(message string, chatID int32) {
	postURL := os.Getenv("TELEGRAM_ENDPOINT") + "/sendMessage"
	fmt.Println(postURL)
	http.PostForm(postURL, url.Values{
		"chat_id": {strconv.Itoa(int(chatID))},
		"text":    {message},
	})
}

func main() {
	godotenv.Load()
	r := setupRouter()
	r.Run(":4321")
}
