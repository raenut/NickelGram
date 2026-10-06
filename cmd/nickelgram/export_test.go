package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestBookExportQueryAndMarkdown(t *testing.T) {
	db := filepath.Join(t.TempDir(), "KoboReader.sqlite")
	schema := `CREATE TABLE content(ContentID TEXT, BookID TEXT, ContentType TEXT, Title TEXT, VolumeIndex INTEGER);
	INSERT INTO content VALUES('part-a','book-1','9','OEBPS/Text/part0043.xhtml',0);
	INSERT INTO content VALUES('part-b','book-1','9','第二章',1);
	CREATE TABLE Bookmark(VolumeID TEXT, ContentID TEXT, Type TEXT, Text TEXT, Annotation TEXT, Hidden BOOL, DateCreated TEXT, BookmarkID TEXT, ChapterProgress REAL);
	INSERT INTO Bookmark VALUES('book-1','part-a','highlight','later','','false','2026-01-01','b',0.8);
	INSERT INTO Bookmark VALUES('book-2','part-a','highlight','other book','','false','2026-01-01','c',0.1);
	INSERT INTO Bookmark VALUES('book-1','part-a','note','first' || char(10) || 'line','my note','false','2026-01-03','a',0.1);
	INSERT INTO Bookmark VALUES('book-1','part-b','highlight','second chapter','','false','2026-01-02','e',0.2);
	INSERT INTO Bookmark VALUES('book-1','missing','highlight','orphan','','false','2026-01-04','f',0.9);
	INSERT INTO Bookmark VALUES('book-1','part-b','highlight','hidden','','true','2026-01-03','d',0.3);`
	if out, err := exec.Command("/usr/bin/sqlite3", db, schema).CombinedOutput(); err != nil {
		t.Fatalf("fixture: %v: %s", err, out)
	}
	book := RecentBook{Volume: "book-1", Title: "A/B: Book", Author: "作者"}
	marks, err := readBookHighlights("/usr/bin/sqlite3", "", db, book)
	if err != nil {
		t.Fatal(err)
	}
	if len(marks) != 4 || marks[0].Text != "first\nline" || marks[0].Annotation != "my note" || marks[0].Chapter != "OEBPS/Text/part0043.xhtml" || marks[0].ChapterID != "part-a" || marks[1].Text != "later" || marks[2].Text != "second chapter" || marks[2].Chapter != "第二章" || marks[3].Text != "orphan" || marks[3].Chapter != "" {
		t.Fatalf("marks: %+v", marks)
	}
	exportedAt := time.Date(2026, 10, 6, 22, 2, 45, 0, time.FixedZone("test", 8*3600))
	want := "---\ntitle: \"A/B: Book\"\nauthor: \"作者\"\nsource: kobo\nexported_at: \"2026-10-06 22:02:45\"\ntags:\n  - \"kobo\"\n---\n\n<hr>\n\n## part0043.xhtml\n\n> first\n> line\n\n批注：\nmy note\n\n<hr>\n\n> later\n\n<hr>\n\n## 第二章\n\n> second chapter\n\n<hr>\n\n## 未知章节\n\n> orphan\n\n<hr>\n\n*footer*\n"
	if got := bookMarkdown(book, marks, "*footer*", []string{"kobo"}, exportedAt); got != want {
		t.Fatalf("markdown = %q, want %q", got, want)
	}
	if got := bookMarkdown(book, marks, "  ", nil, exportedAt); strings.HasSuffix(got, "\n\n") || strings.Contains(got, "footer") || !strings.Contains(got, "tags: []\n") {
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
