package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func main() {
	r := gin.Default()

	// Serve your HTML file
	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")
	r.LoadHTMLGlob("./web/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Endpoint for sending messages
	r.POST("/llm/api", func(c *gin.Context) {
		var chatRequest ChatRequest
		if err := c.BindJSON(&chatRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := CallChatAPI(chatRequest.Model, chatRequest.Messages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": string(response)})
	})

	r.Run() // By default, it serves on :8080
}

// Message represents a single message in the chat stream.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request payload for the chat API.
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// CallChatAPI asynchronously sends a request to the chat API and returns the final response.
func CallChatAPI(model string, messages []Message) ([]byte, error) {
	requestBody, err := json.Marshal(ChatRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error sending request to chat API: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return responseData, nil
}
