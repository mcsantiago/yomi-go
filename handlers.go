package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type JapaneseTokenizerRequest struct {
	Text string `json:"text"`
}

func HandleJapaneseTokenizerRequest(c *gin.Context) {
	var request JapaneseTokenizerRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var tokens = tokenize(request.Text)

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}

func tokenize(str string) []string {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())

	if err != nil {
		panic(err)
	}
	tokens := t.Tokenize(str)
	var surfaces []string

	for _, token := range tokens {
		features := strings.Join(token.Features(), ",")
		surfaces = append(surfaces, token.Surface)
		fmt.Printf("%s\t%v\n", token.Surface, features)
	}

	return surfaces
}
