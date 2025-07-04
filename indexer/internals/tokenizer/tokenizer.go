package tokenizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

func Tokenize(content string) []string {
	var tokens []string

	content = strings.TrimSpace(content)
	if content == "" {
		return tokens
	}

	parsedContent := stopwords.CleanString(strings.ToLower(content), "en", true)

	// Tokenize using regex
	re := regexp.MustCompile(`[a-z0-9]+(?:-[a-z0-9]+)?`)
	rawTokens := re.FindAllString(parsedContent, -1)

	for _, token := range rawTokens {
		stemmed, err := snowball.Stem(token, "english", true)
		if err != nil {
			fmt.Println("stem error: word:", token, "error:", err)
			continue
		}
		tokens = append(tokens, stemmed)
	}

	return tokens
}
