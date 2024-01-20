package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	note "github.com/m-c-frank/note/api"
)

func main() {
	router := gin.Default()

	router.Any("/proxy/*path", proxyHandler)

	err := router.Run("localhost:3020")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Model             string `json:"model"`
	CreatedAt         string `json:"created_at"`
	Message           *Message `json:"message,omitempty"`
	Done              bool   `json:"done"`
	TotalDuration     *int64 `json:"total_duration,omitempty"`
	LoadDuration      *int64 `json:"load_duration,omitempty"`
	PromptEvalCount   *int   `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration *int64 `json:"prompt_eval_duration,omitempty"`
	EvalCount         *int   `json:"eval_count,omitempty"`
	EvalDuration      *int64 `json:"eval_duration,omitempty"`
}

type ChatRequest struct {
	Model    string `json:"model"`
	Messages []Message `json:"messages"`
}

func handleChatRequest(chatRequest []byte) error {
	var req ChatRequest
	err := json.Unmarshal(chatRequest, &req)
	if err != nil {
		return err
	}
	fmt.Println("Handling Chat Request:", req)
	if len(req.Messages) <= 2 {
		noteContent := note.TakeNote(req.Messages[0].Content, "llm/app")
		noteRepoPath, err := note.DetermineRepositoryPath("")
		if err != nil {
			return err
		}
		err = note.PersistNote(noteContent, noteRepoPath)
		if err != nil {
			return err
		}
		err = sendPostRequest(req.Messages)
		if err != nil {
			return err
		}
		return err
	}
	return err
}

func sendPostRequest(messages []Message) error {
	fmt.Println("sending messages now to vectorstore")
	// Convert your messages to JSON
	jsonData, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("Something went wrong marshalling")
		return err
	}

	// Create a new POST request
	url := "http://localhost:8000/api/vectorstore/add"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Something went wrong posting")
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Something went wrong requesting")
		return err
	}
	defer response.Body.Close()


	return nil
}

func proxyHandler(c *gin.Context) {
	requestBody, _ := io.ReadAll(c.Request.Body)

	fmt.Println("Incoming Request JSON:", string(requestBody))

	requestURI := strings.Replace(c.Request.RequestURI, "/proxy", "", -1)
	remoteURL := "http://localhost:11434/api" + requestURI

	if requestURI == "/chat" {
		handleChatRequest(requestBody)
	}


	proxyReq, _ := http.NewRequest(c.Request.Method, remoteURL, bytes.NewBuffer(requestBody))
	proxyReq.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Response JSON:", string(responseBody))

	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Set(key, value)
		}
	}
	c.Writer.WriteHeader(resp.StatusCode)
	_, err = c.Writer.Write(responseBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write response"})
		return
	}
}
