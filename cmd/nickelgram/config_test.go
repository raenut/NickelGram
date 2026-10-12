package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFooterConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name, fields, want, wantURL, wantMD string
	}{
		{"enabled", `"footer":true,"footer_text":"📖 Kobo × NickelGram","footer_url":"https://github.com/raenut/NickelGram","footer_link_text":"NickelGram","md_footer":true,"md_footer_text":"*📖 From Kobo*","md_tags":["kobo"]`, "📖 Kobo × NickelGram", "https://github.com/raenut/NickelGram", "*📖 From Kobo*"},
		{"disabled", `"footer":false,"footer_text":"📖 来自 Kobo","md_footer":false,"md_footer_text":"*📖 摘自 Kobo*"`, "", "", ""},
		{"legacy", `"footer":"旧版文字"`, "旧版文字", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			data := `{"telegram_bot_token":"1:localtest","telegram_chat_id":"1",` + tc.fields + `}`
			if err := os.WriteFile(path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			config, err := readConfigFile(path)
			if err != nil || config.Footer != tc.want || config.FooterURL != tc.wantURL || config.MDFooter != tc.wantMD {
				t.Fatalf("footer = %q, footer_url = %q, md_footer = %q, error = %v", config.Footer, config.FooterURL, config.MDFooter, err)
			}
		})
	}
}

func TestMessageFormattingConfiguration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	data := `{"telegram_bot_token":"1:localtest","telegram_chat_id":"1","header_format":"{author}｜《{title}》","quote_label":"摘录","annotation_label":"备注："}`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	config, err := readConfigFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.HeaderFormat != "{author}｜《{title}》" || config.QuoteLabel != "摘录" || config.AnnotationLabel != "备注：" {
		t.Fatalf("message formatting config = %+v", config)
	}
	got := configuredHighlightMessage(RecentBook{Title: "书名", Author: "作者"}, Highlight{Text: "正文", Annotation: "我的批注"}, "", config)
	if want := "作者｜《书名》\n\n摘录\n正文\n\n备注：\n我的批注"; got != want {
		t.Fatalf("message = %q, want %q", got, want)
	}
}

func TestInvalidHeaderFormatRejected(t *testing.T) {
	for _, format := range []string{"{author}", "{title} {publisher}", "{title\n{author}"} {
		path := filepath.Join(t.TempDir(), "config.json")
		encoded, _ := json.Marshal(format)
		data := `{"telegram_bot_token":"1:localtest","telegram_chat_id":"1","header_format":` + string(encoded) + `}`
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := readConfigFile(path); err == nil {
			t.Errorf("accepted invalid header_format %q", format)
		}
	}
}
