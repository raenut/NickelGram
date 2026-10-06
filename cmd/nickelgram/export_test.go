package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestBookExportQueryAndMarkdown(t *testing.T) {
	db := filepath.Join(t.TempDir(), "KoboReader.sqlite")
	schema := `CREATE TABLE Bookmark(VolumeID TEXT, Type TEXT, Text TEXT, Annotation TEXT, Hidden BOOL, DateCreated TEXT, BookmarkID TEXT);
INSERT INTO Bookmark VALUES('book-1','highlight','second','','false','2026-01-02','b');
INSERT INTO Bookmark VALUES('book-2','highlight','other book','','false','2026-01-01','c');
INSERT INTO Bookmark VALUES('book-1','note','first' || char(10) || 'line','my note','false','2026-01-01','a');
INSERT INTO Bookmark VALUES('book-1','highlight','hidden','','true','2026-01-03','d');`
	if out, err := exec.Command("/usr/bin/sqlite3", db, schema).CombinedOutput(); err != nil {
		t.Fatalf("fixture: %v: %s", err, out)
	}
	book := RecentBook{Volume: "book-1", Title: "A/B: Book", Author: "作者"}
	marks, err := readBookHighlights("/usr/bin/sqlite3", "", db, book)
	if err != nil {
		t.Fatal(err)
	}
	if len(marks) != 2 || marks[0].Text != "first\nline" || marks[0].Annotation != "my note" || marks[1].Text != "second" {
		t.Fatalf("marks: %+v", marks)
	}
	want := "---\ntitle: \"A/B: Book\"\nauthor: \"作者\"\nsource: kobo\ntags:\n  - \"kobo\"\n---\n\n<hr>\n\n> first\n> line\n\n批注：\nmy note\n\n<hr>\n\n> second\n\n<hr>\n\n*footer*\n"
	if got := bookMarkdown(book, marks, "*footer*", []string{"kobo"}); got != want {
		t.Fatalf("markdown = %q, want %q", got, want)
	}
	if got := bookMarkdown(book, marks, "  ", nil); strings.HasSuffix(got, "\n\n") || strings.Contains(got, "footer") || !strings.Contains(got, "tags: []\n") {
		t.Fatalf("empty footer left extra space: %q", got)
	}
	if got := safeMarkdownName(book.Title, book.Volume); !strings.HasPrefix(got, "A_B_ Book-") || !strings.HasSuffix(got, ".md") {
		t.Fatalf("filename: %q", got)
	}
}

func TestSendDocumentUploadsMarkdownFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "book.md")
	if err := os.WriteFile(path, []byte("# book\n"), 0600); err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/sendDocument") {
			t.Errorf("path: %s", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Error(err)
		}
		f, header, err := r.FormFile("document")
		if err != nil {
			t.Error(err)
		} else {
			defer f.Close()
			b, _ := io.ReadAll(f)
			if header.Filename != "book.md" || string(b) != "# book\n" || r.FormValue("chat_id") != "123" {
				t.Errorf("upload: %q %q %q", header.Filename, b, r.FormValue("chat_id"))
			}
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"ok":true,"result":{"message_id":1}}`)), Header: make(http.Header)}, nil
	})}
	tg := Telegram{Client: client, BaseURL: "https://example.test"}
	if outcome := tg.SendDocument(Config{BotToken: "1:token", ChatID: "123"}, path, "book.md"); !outcome.OK {
		t.Fatalf("send: %+v", outcome)
	}
}
