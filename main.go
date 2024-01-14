package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    note "github.com/m-c-frank/note"
)

func main() {
    r := gin.Default()

    // Serve your HTML file
    r.LoadHTMLGlob("templates/*")

    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{})
    })

    // Endpoint for sending messages
    r.POST("/api/note", func(c *gin.Context) {
        note.api()
        // Handle your chat logic here
        // For example, you can read the message sent by the user and process it
    })

    r.Run() // By default, it serves on :8080
}

