package japanese_tokenizer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/scriptsight/tokenizerservice/internal/japanese_dict"
)

type JapaneseTokenizerRequest struct {
	Text string `json:"text"`
}

type JapaneseTokenizerResponse struct {
	Translation string     `json:"translation"`
	Tokens      []TokenExt `json:"tokens"`
}

type TokenExt struct {
	TokenData tokenizer.TokenData `json:"token_data"`
	Glossary  string              `json:"glossary"`
}

func HandleJapaneseTokenizerRequest(c *gin.Context, dict map[string][]japanese_dict.JmdictEntry) {
	var request JapaneseTokenizerRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, kotoriTokenize(request.Text, dict))
}

func kotoriTokenize(str string, dict map[string][]japanese_dict.JmdictEntry) JapaneseTokenizerResponse {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())

	if err != nil {
		panic(err)
	}
	tokens := t.Tokenize(str)
	var tokenExts []TokenExt

	for _, token := range tokens {
		entries := dict[token.Surface]
		translation := ""
		if len(entries) > 0 {
			translation = entries[0].Sense[0].Glossary[0].Content
		}
		tokenExts = append(tokenExts, TokenExt{
			TokenData: tokenizer.NewTokenData(token),
			Glossary:  translation,
		})
	}

	return JapaneseTokenizerResponse{
		Translation: "",
		Tokens:      tokenExts,
	}
}