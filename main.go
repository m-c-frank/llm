package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Any("/proxy/*path", proxyHandler)

	err := router.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func proxyHandler(c *gin.Context) {
	// Read and log the incoming request
	requestBody, _ := io.ReadAll(c.Request.Body)

	fmt.Println("Incoming Request JSON:", string(requestBody))

	// Forward the request to the remote URL
	remoteURL := "http://192.168.2.177:11434/api" + strings.Replace(c.Request.RequestURI, "/proxy", "", -1)
	fmt.Println(remoteURL)
	proxyReq, _ := http.NewRequest(c.Request.Method, remoteURL, bytes.NewBuffer(requestBody))
	proxyReq.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read and log the response from the remote URL
	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Response JSON:", string(responseBody))

	// Relay the response back to the original requester
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
