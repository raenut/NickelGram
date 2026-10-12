package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"unicode/utf16"
)

func TestSendMessageLinksOnlyFooterNameWithoutPreview(t *testing.T) {
	footer := "📖 Kobo × NickelGram"
	selection := "NickelGram 普通选区测试\n下一段"
	message := selectionMessage(RecentBook{Title: "九诗心：暗夜里的文学启明", Author: "黄晓丹"}, selection, footer)
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var payload struct {
			Text               string `json:"text"`
			LinkPreviewOptions struct {
				IsDisabled bool `json:"is_disabled"`
			} `json:"link_preview_options"`
			Entities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
			} `json:"entities"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload.Text != message || !payload.LinkPreviewOptions.IsDisabled || len(payload.Entities) != 2 {
			t.Errorf("payload = %+v", payload)
		} else {
			quoteEntity, linkEntity := payload.Entities[0], payload.Entities[1]
			quoteStart := strings.Index(message, "❝\n") + len("❝\n")
			if quoteEntity.Type != "code" || quoteEntity.Offset != len(utf16.Encode([]rune(message[:quoteStart]))) || quoteEntity.Length != len(utf16.Encode([]rune(selection))) {
				t.Errorf("quote entity = %+v", quoteEntity)
			}
			wantLinkOffset := len(utf16.Encode([]rune(strings.TrimSuffix(message, "NickelGram"))))
			if linkEntity.Type != "text_link" || linkEntity.Offset != wantLinkOffset || linkEntity.Length != len("NickelGram") || linkEntity.URL != "https://github.com/raenut/NickelGram" {
				t.Errorf("link entity = %+v", linkEntity)
			}
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"ok":true,"result":{"message_id":1}}`)), Header: make(http.Header)}, nil
	})}
	tg := Telegram{Client: client, BaseURL: "https://example.test"}
	config := Config{BotToken: "1:token", ChatID: "123", Footer: footer, FooterURL: "https://github.com/raenut/NickelGram", FooterLinkText: "NickelGram"}
	if outcome := tg.SendMessage(config, message); !outcome.OK {
		t.Fatalf("send: %+v", outcome)
	}
}

func TestQuoteEntityExcludesAnnotationAndUsesCustomLabel(t *testing.T) {
	config := Config{QuoteLabel: "摘录", AnnotationLabel: "备注："}
	message := configuredHighlightMessage(
		RecentBook{Title: "书名", Author: "作者"},
		Highlight{Text: "第一段\n第二段", Annotation: "我的批注"},
		"",
		config,
	)
	entities := telegramMessageEntities(config, message)
	if len(entities) != 1 {
		t.Fatalf("entities = %+v", entities)
	}
	quoteStart := strings.Index(message, "摘录\n") + len("摘录\n")
	if entities[0]["type"] != "code" || entities[0]["offset"] != utf16Length(message[:quoteStart]) || entities[0]["length"] != utf16Length("第一段\n第二段") {
		t.Fatalf("quote entity = %+v", entities[0])
	}
}

func TestQuoteEntityDoesNotIncludeTrailingWhitespace(t *testing.T) {
	config := Config{}
	message := configuredSelectionMessage(RecentBook{Title: "书名"}, "书摘正文\n", "", config)
	entities := telegramMessageEntities(config, message)
	if len(entities) != 1 || entities[0]["length"] != utf16Length("书摘正文") {
		t.Fatalf("entities = %+v", entities)
	}
}
