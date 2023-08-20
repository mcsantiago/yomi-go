package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scriptsight/tokenizerservice/internal/japanese_dict"
	"github.com/scriptsight/tokenizerservice/internal/japanese_tokenizer"
)

var tpl *template.Template

func main() {
	tpl, _ = template.ParseGlob("templates/*.html")
	filepath := "./data/JMdict_e"
	reader, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		return
	}

	dict, err := japanese_dict.LoadJmdictMap(reader)
	if err != nil {
		fmt.Printf("Error loading dictionary: %s", err)
		return
	}

	fmt.Println("Loaded dictionary")
	fmt.Println(dict["日本"][0].Kanji[0].Expression)

	r := gin.Default()
	// Configuring CORS for all origins for simplicity, adjust for production needs
	r.Use(cors.Default())
	r.GET("/", func(c *gin.Context) {
		handleIndexPage(c)
	})
	r.POST("/japanese/kotori/tokenize", func(c *gin.Context) {
		japanese_tokenizer.HandleJapaneseTokenizerRequest(c, dict)
	})
	r.POST("/japanese/kotori/analyze", func(c *gin.Context) {
		japanese_tokenizer.HandleJapaneseAnalyzerRequest(c, dict, tpl)
	})
	r.Run() // Listen and serve on 0.0.0.0:8080
}

func handleIndexPage(c *gin.Context) {
	tpl.ExecuteTemplate(c.Writer, "index.html", nil)
}
