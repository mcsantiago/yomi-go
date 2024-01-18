package japanese_tokenizer

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "cloud.google.com/go/translate/apiv3/translatepb"
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
	TokenData     tokenizer.TokenData         `json:"token_data"`
	JmdictEntries []japanese_dict.JmdictEntry `json:"jmdict_entries"`
}

func HandleJapaneseTokenizerRequest(c *gin.Context, dict map[string][]japanese_dict.JmdictEntry) {
	var request JapaneseTokenizerRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	translation, err := translateText("scriptsight-396805", "ja", "en-us", request.Text)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := JapaneseTokenizerResponse{
		Translation: translation,
		Tokens:      kotoriTokenize(request.Text, dict),
	}

	c.JSON(http.StatusOK, response)
}

/*
*
Serves a HTML page with the Japanese Tokenizer Response
*/
func HandleJapaneseAnalyzerRequest(c *gin.Context, dict map[string][]japanese_dict.JmdictEntry, tpl *template.Template) {
	fmt.Println("HandleJapaneseAnalyzerRequest")
	text := c.PostForm("text")
	fmt.Println(c.PostForm("text"))
	tokens := kotoriTokenize(text, dict)
	translation, err := translateText("scriptsight-396805", "ja", "en-us", text)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tpl.ExecuteTemplate(c.Writer, "tokenCards.html", JapaneseTokenizerResponse{
		Translation: translation,
		Tokens:      tokens,
	})
}

func kotoriTokenize(str string, dict map[string][]japanese_dict.JmdictEntry) []TokenExt {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())

	if err != nil {
		panic(err)
	}
	tokens := t.Tokenize(str)
	var tokenExts []TokenExt

	for _, token := range tokens {
		entries, ok := dict[token.Surface]
		if !ok {
			entries = []japanese_dict.JmdictEntry{}
		}
		tokenExts = append(tokenExts, TokenExt{
			TokenData:     tokenizer.NewTokenData(token),
			JmdictEntries: entries,
		})
	}

	return tokenExts
}

func translateText(projectID string, sourceLang string, targetLang string, text string) (string, error) {
	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return "", fmt.Errorf("NewTranslationClient: %v", err)
	}
	defer client.Close()

	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		SourceLanguageCode: sourceLang,
		TargetLanguageCode: targetLang,
		MimeType:           "text/plain", // Mime types: "text/plain", "text/html"
		Contents:           []string{text},
	}

	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return "", fmt.Errorf("TranslateText: %v", err)
	}

	var translations []string
	// Display the translation for each input text provided
	for _, translation := range resp.GetTranslations() {
		translations = append(translations, translation.GetTranslatedText())
		fmt.Printf("Translated text: %v\n", translation.GetTranslatedText())
	}

	return strings.Join(translations, ""), nil
}
