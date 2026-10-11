package main

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCurrentBookUsesOpenEventBeforeLastReadUpdates(t *testing.T) {
	db := filepath.Join(t.TempDir(), "KoboReader.sqlite")
	schema := `CREATE TABLE content(ContentID TEXT, ContentType TEXT, Title TEXT, Attribution TEXT, DateLastRead TEXT);
CREATE TABLE Event(EventType INTEGER, ContentID TEXT, LastOccurrence TEXT);
CREATE TABLE Bookmark(VolumeID TEXT, Type TEXT, Text TEXT, Hidden BOOL, DateCreated TEXT, DateModified TEXT);
INSERT INTO content VALUES('old','6','旧书','旧作者','2026-10-06T10:00:00Z');
INSERT INTO content VALUES('new','6','新书','新作者',NULL);
INSERT INTO Event VALUES(1021,'old','2026-10-06T10:00:00.000');
INSERT INTO Event VALUES(1020,'new','2026-10-06T10:01:00.000');`
	if out, err := exec.Command("/usr/bin/sqlite3", db, schema).CombinedOutput(); err != nil {
		t.Fatalf("fixture: %v: %s", err, out)
	}
	book, err := readCurrentBook("/usr/bin/sqlite3", "", db)
	if err != nil || book.Title != "新书" || book.Author != "新作者" {
		t.Fatalf("open event: %+v, %v", book, err)
	}
	if out, err := exec.Command("/usr/bin/sqlite3", db, `INSERT INTO Event VALUES(1021,'new','2026-10-06T10:02:00.000');`).CombinedOutput(); err != nil {
		t.Fatalf("fixture: %v: %s", err, out)
	}
	book, err = readCurrentBook("/usr/bin/sqlite3", "", db)
	if err != nil || book.Title != "旧书" {
		t.Fatalf("closed event must fall back: %+v, %v", book, err)
	}
}
