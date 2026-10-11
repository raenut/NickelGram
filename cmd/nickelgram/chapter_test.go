package main

import (
	"archive/zip"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestChapterWithinSharedKepubFile(t *testing.T) {
	root := t.TempDir()
	db := filepath.Join(root, "KoboReader.sqlite")
	schema := `CREATE TABLE content(ContentID TEXT, BookID TEXT, ContentType TEXT, Title TEXT, VolumeIndex INTEGER);
CREATE TABLE Bookmark(VolumeID TEXT, ContentID TEXT, Type TEXT, Text TEXT, Annotation TEXT, Hidden BOOL, ChapterProgress REAL, BookmarkID TEXT, StartContainerPath TEXT);
INSERT INTO content VALUES('book-1!!Text/part.xhtml','book-1','9','Text/part.xhtml',1);
INSERT INTO content VALUES('book-1!!Text/part.xhtml-1','book-1','899','卷首',2);
INSERT INTO content VALUES('book-1!!Text/part.xhtml#chap-seven-1','book-1','899','第七章',2);
INSERT INTO content VALUES('book-1!!Text/part.xhtml#chap-eight-1','book-1','899','不规则的章名',3);
INSERT INTO Bookmark VALUES('book-1','book-1!!Text/part.xhtml','highlight','intro','','false',0.1,'b','span#kobo\.1\.1');
INSERT INTO Bookmark VALUES('book-1','book-1!!Text/part.xhtml','highlight','sample','','false',0.8,'a','span#kobo\.405\.1');`
	if out, err := exec.Command("/usr/bin/sqlite3", db, schema).CombinedOutput(); err != nil {
		t.Fatalf("fixture: %v: %s", err, out)
	}
	dir := filepath.Join(root, "kepub")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	file, err := os.Create(filepath.Join(dir, "book-1"))
	if err != nil {
		t.Fatal(err)
	}
	archive := zip.NewWriter(file)
	member, err := archive.Create("Text/part.xhtml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = member.Write([]byte(`<span id="kobo.1.1">intro</span><h1 id="chap-seven">Seven</h1><h1 id="chap-eight">Eight</h1><span id="kobo.405.1">sample</span>`)); err != nil {
		t.Fatal(err)
	}
	if err = archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err = file.Close(); err != nil {
		t.Fatal(err)
	}
	marks, err := readBookHighlights("/usr/bin/sqlite3", "", db, RecentBook{Volume: "book-1"})
	if err != nil || len(marks) != 2 || marks[0].Chapter != "卷首" || marks[1].Chapter != "不规则的章名" {
		t.Fatalf("chapter: %+v, %v", marks, err)
	}
}
