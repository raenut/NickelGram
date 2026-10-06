package main

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectionMessage(t *testing.T) {
	book := RecentBook{Title: "书名", Author: "作者"}
	got := selectionMessage(book, "选中的原文", "footer")
	want := "书名, 作者\n\n❝\n选中的原文\n\nfooter"
	if got != want {
		t.Fatalf("message = %q, want %q", got, want)
	}
}

func TestSelectionParagraphSeparators(t *testing.T) {
	book := RecentBook{Title: "书名"}
	for _, tc := range []struct{ input, want string }{
		{"单段原文", "单段原文"},
		{"段内\n换行", "段内\n换行"},
		{"第一段\n\n第二段", "第一段\n\n第二段"},
		{"第一段\r\n\r\n第二段", "第一段\n\n第二段"},
		{"第一段\r第二段", "第一段\n第二段"},
		{"第一段\u2028第二段", "第一段\n第二段"},
		{"第一段\u2029第二段", "第一段\n\n第二段"},
	} {
		got := selectionMessage(book, tc.input, "")
		if want := "书名\n\n❝\n" + tc.want; got != want {
			t.Errorf("input %q: got %q, want %q", tc.input, got, want)
		}
	}
}

func TestSelectionNewlinesSurviveShellAndTelegramJSON(t *testing.T) {
	selection := "第一段\n\n第二段's"
	quoted := "'" + strings.ReplaceAll(selection, "'", "'\"'\"'") + "'"
	passed, err := exec.Command("/bin/sh", "-c", "printf '%s' "+quoted).Output()
	if err != nil || string(passed) != selection {
		t.Fatalf("shell argument: %q, %v", passed, err)
	}
	want := selectionMessage(RecentBook{Title: "书名"}, string(passed), "")
	body, err := json.Marshal(map[string]any{"chat_id": "1", "text": want})
	if err != nil {
		t.Fatal(err)
	}
	var payload struct{ Text string `json:"text"` }
	if err := json.Unmarshal(body, &payload); err != nil || payload.Text != want {
		t.Fatalf("Telegram JSON text = %q, err = %v", payload.Text, err)
	}
}

func TestReopenedHighlightUsesCompleteBookmarkText(t *testing.T) {
	db := filepath.Join(t.TempDir(), "KoboReader.sqlite")
	schema := `CREATE TABLE Bookmark(VolumeID TEXT, Type TEXT, Text TEXT, Hidden BOOL);
INSERT INTO Bookmark VALUES('book-1','highlight','这是第一句足够长的内容。' || char(10) || char(10) || '这是完整高亮的第二句。',0);
INSERT INTO Bookmark VALUES('book-2','highlight','这是第一句足够长的内容。别的书。',0);`
	if out, err := exec.Command("/usr/bin/sqlite3", db, schema).CombinedOutput(); err != nil {
		t.Fatalf("fixture: %v: %s", err, out)
	}
	book := RecentBook{Volume: "book-1", Title: "书名"}
	want := "这是第一句足够长的内容。\n\n这是完整高亮的第二句。"
	for _, selected := range []string{"这是第一句足够长的内容。", "这是第一句足够长的内容"} {
		if got := readMatchingBookmarkText("/usr/bin/sqlite3", "", db, book, selected); got != want {
			t.Fatalf("selection %q matched %q", selected, got)
		}
	}
	if got := selectionMessage(book, want, "footer"); !strings.Contains(got, "内容。\n\n这是完整") {
		t.Fatalf("complete message lost second sentence: %q", got)
	}
	if got := readMatchingBookmarkText("/usr/bin/sqlite3", "", db, book, "无关选区。"); got != "" {
		t.Fatalf("unrelated selection matched %q", got)
	}
	if out, err := exec.Command("/usr/bin/sqlite3", db, `INSERT INTO Bookmark VALUES('book-1','highlight','这是第一句足够长的内容。另一种后文。',0);`).CombinedOutput(); err != nil {
		t.Fatalf("fixture: %v: %s", err, out)
	}
	if got := readMatchingBookmarkText("/usr/bin/sqlite3", "", db, book, "这是第一句足够长的内容。"); got != "" {
		t.Fatalf("ambiguous selection matched %q", got)
	}
}

func TestLatestHighlightStaysInBookAndUsesModification(t *testing.T) {
	db := filepath.Join(t.TempDir(), "KoboReader.sqlite")
	schema := `CREATE TABLE Bookmark(VolumeID TEXT, Type TEXT, Text TEXT, Annotation TEXT, DateCreated TEXT, DateModified TEXT, Hidden BOOL, BookmarkID TEXT DEFAULT '');
INSERT INTO Bookmark(VolumeID,Type,Text,Annotation,DateCreated,DateModified,Hidden) VALUES('book-1','highlight','older creation','written later','2026-01-01T00:00:00Z','2026-01-04T00:00:00Z',0);
INSERT INTO Bookmark(VolumeID,Type,Text,Annotation,DateCreated,DateModified,Hidden) VALUES('book-1','highlight','newer creation','','2026-01-03T00:00:00Z',NULL,0);
INSERT INTO Bookmark(VolumeID,Type,Text,Annotation,DateCreated,DateModified,Hidden) VALUES('book-2','note','other book','private','2026-01-05T00:00:00Z','2026-01-05T00:00:00Z',0);
INSERT INTO Bookmark(VolumeID,Type,Text,Annotation,DateCreated,DateModified,Hidden) VALUES('book-1','highlight','hidden','','2026-01-06T00:00:00Z',NULL,1);`
	if out, err := exec.Command("/usr/bin/sqlite3", db, schema).CombinedOutput(); err != nil {
		t.Fatalf("fixture: %v: %s", err, out)
	}
	book := RecentBook{Volume: "book-1", Title: "书名", Author: "作者"}
	mark, err := readLatestHighlight("/usr/bin/sqlite3", "", db, book)
	if err != nil {
		t.Fatal(err)
	}
	if mark.Text != "older creation" || mark.Annotation != "written later" {
		t.Fatalf("wrong highlight: %+v", mark)
	}
	if got := highlightMessage(book, mark, "footer"); got != "书名, 作者\n\n❝\nolder creation\n\n批注：\nwritten later\n\nfooter" {
		t.Fatalf("message: %q", got)
	}
}
