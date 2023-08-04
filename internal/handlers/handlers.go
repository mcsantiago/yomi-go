package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type JapaneseTokenizerRequest struct {
	Text string `json:"text"`
}

type JapaneseTokenizerResponse struct {
	KnownTokens   []Token `json:"known_tokens"`
	UnknownTokens []Token `json:"unknown_tokens"`
}

type Token struct {
	Surface  string   `json:"surface"`
	Features []string `json:"features"`
}

func HandleJapaneseTokenizerRequest(c *gin.Context) {
	var request JapaneseTokenizerRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response = kotoriTokenize(request.Text)
	c.JSON(http.StatusOK, response)
}

func kotoriTokenize(str string) JapaneseTokenizerResponse {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	var JapaneseTokenizerResponse JapaneseTokenizerResponse

	if err != nil {
		panic(err)
	}
	tokens := t.Tokenize(str)
	var knownResponseTokens []Token
	var unKnownResponseTokens []Token

	for _, token := range tokens {
		if token.Surface == "" {
			continue
		}

		if token.Class == tokenizer.KNOWN {
			knownResponseTokens = append(knownResponseTokens, Token{
				Surface:  token.Surface,
				Features: token.Features(),
			})
		}
		if token.Class == tokenizer.UNKNOWN {
			unKnownResponseTokens = append(unKnownResponseTokens, Token{
				Surface:  token.Surface,
				Features: token.Features(),
			})
		}
	}

	JapaneseTokenizerResponse.KnownTokens = knownResponseTokens
	JapaneseTokenizerResponse.UnknownTokens = unKnownResponseTokens

	return JapaneseTokenizerResponse
}
