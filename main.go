package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Configuring CORS for all origins for simplicity, adjust for production needs
	r.Use(cors.Default())
	r.POST("/japanese/kotori/tokenize", HandleJapaneseTokenizerRequest)
	r.Run() // Listen and serve on 0.0.0.0:8080
}
