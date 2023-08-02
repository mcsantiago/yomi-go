package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/japanese/kotori/tokenize", HandleJapaneseTokenizerRequest)
	r.Run() // Listen and serve on 0.0.0.0:8080
}
