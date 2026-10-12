package main

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

var indentationSpaces = regexp.MustCompile(`[ \t\x{00a0}\x{1680}\x{2000}-\x{200a}\x{202f}\x{205f}\x{3000}]+`)
var newlineRun = regexp.MustCompile(`\n+`)

func normalizeSelectionParagraphs(text string) string {
	text = strings.NewReplacer("\r\n", "\n", "\r", "\n", "\u2028", "\n", "\u2029", "\n").Replace(text)
	// Kobo stores paragraph starts in highlights as LF followed by a tab.
	// Keep the LF and remove the indentation marker as one line break.
	text = strings.ReplaceAll(text, "\n\t", "\n")
	text = indentationSpaces.ReplaceAllStringFunc(text, func(spaces string) string {
		if spaces == " " {
			return spaces
		}
		return "\n"
	})
	return newlineRun.ReplaceAllString(text, "\n")
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
