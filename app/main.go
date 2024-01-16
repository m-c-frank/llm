package main

import "github.com/gin-gonic/gin"

func main() {
    router := gin.Default()

    // Serve static files from the Svelte build directory
    router.Static("/", "./web/build")

    router.Run() // By default, it serves on :8080
}

