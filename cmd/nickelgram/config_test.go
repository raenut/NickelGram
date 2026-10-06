package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFooterConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name, fields, want, wantMD string
	}{
		{"enabled", `"footer":true,"footer_text":"📖 来自 Kobo","md_footer":true,"md_footer_text":"*📖 摘自 Kobo*","md_tags":["kobo"]`, "📖 来自 Kobo", "*📖 摘自 Kobo*"},
		{"disabled", `"footer":false,"footer_text":"📖 来自 Kobo","md_footer":false,"md_footer_text":"*📖 摘自 Kobo*"`, "", ""},
		{"legacy", `"footer":"旧版文字"`, "旧版文字", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			data := `{"telegram_bot_token":"1:localtest","telegram_chat_id":"1",` + tc.fields + `}`
			if err := os.WriteFile(path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			config, err := readConfigFile(path)
			if err != nil || config.Footer != tc.want || config.MDFooter != tc.wantMD {
				t.Fatalf("footer = %q, md_footer = %q, error = %v", config.Footer, config.MDFooter, err)
			}
		})
	}
}
