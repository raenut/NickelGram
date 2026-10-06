package main

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

var blankLines = regexp.MustCompile(`\n{2,}`)

func normalizeSelectionParagraphs(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.NewReplacer("\r", "\n", "\u2028", "\n", "\u2029", "\n\n").Replace(text)
	return blankLines.ReplaceAllString(text, "\n\n")
}

func validateSelection(text string) (string, error) {
	if !utf8.ValidString(text) || strings.ContainsRune(text, 0) {
		return "", errors.New("选区不是有效 UTF-8 文本；未发送。")
	}
	if strings.TrimSpace(text) == "" {
		return "", errors.New("选区为空；未发送。")
	}
	return text, nil
}
