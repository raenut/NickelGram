package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const koboDB = "/mnt/onboard/.kobo/KoboReader.sqlite"

type RecentBook struct {
	Volume string `json:"volume"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Highlight struct {
	ID         string `json:"id"`
	Volume     string `json:"volume"`
	Text       string `json:"text"`
	Annotation string `json:"annotation"`
	Chapter    string `json:"chapter"`
	ChapterID  string `json:"chapter_id"`
	Modified   string `json:"modified"`
}

const highlightColumns = `BookmarkID AS id, hex(VolumeID) AS volume,
 hex(Text) AS text, hex(COALESCE(Annotation,'')) AS annotation,
 CASE WHEN COALESCE(DateModified,'') > COALESCE(DateCreated,'')
      THEN DateModified ELSE DateCreated END AS modified`

const visibleHighlightCondition = `Type IN ('highlight','note')
   AND length(trim(COALESCE(Text,'')))>0
   AND lower(COALESCE(CAST(Hidden AS TEXT),'false')) NOT IN ('1','true')`

const recentBookSQL = `SELECT hex(ContentID) AS volume,
 hex(COALESCE(Title,'')) AS title,
 hex(COALESCE(Attribution,'')) AS author
 FROM content
 WHERE CAST(ContentType AS TEXT)='6'
   AND DateLastRead IS NOT NULL
   AND length(trim(CAST(DateLastRead AS TEXT)))>0
 ORDER BY DateLastRead DESC, rowid DESC
 LIMIT 1;`

// A newly created highlight can identify its book before DateLastRead advances.
// If no visible bookmark is newer than the last read timestamp, retain the
// existing DateLastRead behavior.
const currentBookSQL = `SELECT hex(c.ContentID) AS volume,
 hex(COALESCE(c.Title,'')) AS title,
 hex(COALESCE(c.Attribution,'')) AS author
 FROM content c
 WHERE c.ContentID = COALESCE(
   (SELECT b.VolumeID FROM Bookmark b
    INNER JOIN content bc ON bc.ContentID=b.VolumeID AND CAST(bc.ContentType AS TEXT)='6'
    WHERE b.Type IN ('highlight','note')
      AND length(trim(COALESCE(b.Text,'')))>0
      AND lower(COALESCE(CAST(b.Hidden AS TEXT),'false')) NOT IN ('1','true')
      AND julianday(CASE WHEN COALESCE(b.DateModified,'') > COALESCE(b.DateCreated,'')
                         THEN b.DateModified ELSE b.DateCreated END) >
          COALESCE((SELECT max(julianday(DateLastRead)) FROM content
                    WHERE CAST(ContentType AS TEXT)='6'), 0)
    ORDER BY julianday(CASE WHEN COALESCE(b.DateModified,'') > COALESCE(b.DateCreated,'')
                            THEN b.DateModified ELSE b.DateCreated END) DESC,
             b.BookmarkID DESC
    LIMIT 1),
   (SELECT ContentID FROM content
    WHERE CAST(ContentType AS TEXT)='6'
      AND DateLastRead IS NOT NULL
      AND length(trim(CAST(DateLastRead AS TEXT)))>0
    ORDER BY DateLastRead DESC, rowid DESC LIMIT 1))
 LIMIT 1;`

func sqliteQuery(binary, lib, db, query string) ([]byte, error) {
	limit := 1500 * time.Millisecond
	if db == ":memory:" {
		limit = 800 * time.Millisecond
	}
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, "-init", "/dev/null", "-batch", "-bail", "-readonly", "-json", db, ".timeout 1000", "PRAGMA query_only=ON;", query)
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, "LD_LIBRARY_PATH=") {
			cmd.Env = append(cmd.Env, v)
		}
	}
	cmd.Env = append(cmd.Env, "LD_LIBRARY_PATH="+lib)
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.New("SQLite 只读查询失败；未发送。")
	}
	return out, nil
}

func verifySQLiteRuntime(root string) error {
	binary := filepath.Join(root, "sqlite3")
	info, err := os.Stat(binary)
	if err != nil || info.Mode()&0111 == 0 {
		return errors.New("sqlite3 缺失或无执行权限；请重新安装 NickelGram。")
	}
	if _, err = os.Stat(filepath.Join(root, "lib/libsqlite3.so.0")); err != nil {
		return errors.New("SQLite 动态库缺失；请重新安装 NickelGram。")
	}
	if _, err = os.Stat("/lib/ld-linux-armhf.so.3"); err != nil {
		return errors.New("当前 Kobo 缺少所需 ARM 加载器；未发送。")
	}
	_, err = sqliteQuery(binary, filepath.Join(root, "lib"), ":memory:", "SELECT sqlite_version();")
	return err
}

func decodeHexField(s string) (string, error) {
	b, err := hex.DecodeString(s)
	if err != nil || !utf8.Valid(b) {
		return "", errors.New("数据库文本编码无效；未发送。")
	}
	return string(b), nil
}

func readRecentBook(binary, lib, db string) (RecentBook, error) {
	return readBookQuery(binary, lib, db, recentBookSQL)
}

func readCurrentBook(binary, lib, db string) (RecentBook, error) {
	return readBookQuery(binary, lib, db, currentBookSQL)
}

func readBookQuery(binary, lib, db, query string) (RecentBook, error) {
	out, err := sqliteQuery(binary, lib, db, query)
	if err != nil {
		return RecentBook{}, err
	}
	var rows []RecentBook
	if len(strings.TrimSpace(string(out))) == 0 || json.Unmarshal(out, &rows) != nil || len(rows) != 1 {
		return RecentBook{}, errors.New("未找到最近打开的书籍；未发送。")
	}
	b := rows[0]
	if b.Volume, err = decodeHexField(b.Volume); err != nil {
		return RecentBook{}, err
	}
	if b.Title, err = decodeHexField(b.Title); err != nil {
		return RecentBook{}, err
	}
	if b.Author, err = decodeHexField(b.Author); err != nil {
		return RecentBook{}, err
	}
	if strings.TrimSpace(b.Volume) == "" || strings.TrimSpace(b.Title) == "" {
		return RecentBook{}, errors.New("最近打开的书籍缺少可用书名；未发送。")
	}
	return b, nil
}

func readLatestHighlight(binary, lib, db string, book RecentBook) (Highlight, error) {
	volumeHex := strings.ToUpper(hex.EncodeToString([]byte(book.Volume)))
	query := `SELECT ` + highlightColumns + ` FROM Bookmark
 WHERE hex(VolumeID)='` + volumeHex + `'
   AND ` + visibleHighlightCondition + `
   AND (length(trim(COALESCE(DateCreated,'')))>0 OR length(trim(COALESCE(DateModified,'')))>0)
 ORDER BY modified DESC
 LIMIT 2;`
	out, err := sqliteQuery(binary, lib, db, query)
	if err != nil {
		return Highlight{}, err
	}
	var rows []Highlight
	if json.Unmarshal(out, &rows) != nil || len(rows) == 0 {
		return Highlight{}, errors.New("当前书没有可用的高亮或批注；未发送。")
	}
	if len(rows) > 1 && rows[0].Modified == rows[1].Modified {
		return Highlight{}, errors.New("当前书有多条相同最新时间的高亮；未发送。")
	}
	return decodeHighlight(rows[0])
}

func decodeHighlight(mark Highlight) (Highlight, error) {
	var err error
	if mark.Volume, err = decodeHexField(mark.Volume); err != nil {
		return Highlight{}, err
	}
	if mark.Text, err = decodeHexField(mark.Text); err != nil {
		return Highlight{}, err
	}
	if mark.Annotation, err = decodeHexField(mark.Annotation); err != nil {
		return Highlight{}, err
	}
	return mark, nil
}

// readBookHighlights follows the book's content order, then each mark's position
// within that content. Marks without a matching content row sort last.
func readBookHighlights(binary, lib, db string, book RecentBook) ([]Highlight, error) {
	volumeHex := strings.ToUpper(hex.EncodeToString([]byte(book.Volume)))
	query := `SELECT hex(b.Text) AS text,
	 hex(COALESCE(b.Annotation,'')) AS annotation,
	 hex(COALESCE(chapter.Title,'')) AS chapter,
	 hex(b.ContentID) AS chapter_id
	 FROM Bookmark b
	 LEFT JOIN content chapter ON chapter.ContentID=b.ContentID
	   AND chapter.BookID=b.VolumeID AND CAST(chapter.ContentType AS TEXT)='9'
	 WHERE hex(b.VolumeID)='` + volumeHex + `'
	   AND b.Type IN ('highlight','note')
	   AND length(trim(COALESCE(b.Text,'')))>0
	   AND lower(COALESCE(CAST(b.Hidden AS TEXT),'false')) NOT IN ('1','true')
	 ORDER BY chapter.VolumeIndex IS NULL, CAST(chapter.VolumeIndex AS INTEGER),
	          b.ChapterProgress, b.BookmarkID;`
	out, err := sqliteQuery(binary, lib, db, query)
	if err != nil {
		return nil, err
	}
	var rows []Highlight
	if json.Unmarshal(out, &rows) != nil {
		return nil, errors.New("Bookmark 输出无法解析；未发送。")
	}
	for i := range rows {
		if rows[i].Text, err = decodeHexField(rows[i].Text); err != nil {
			return nil, err
		}
		if rows[i].Annotation, err = decodeHexField(rows[i].Annotation); err != nil {
			return nil, err
		}
		if rows[i].Chapter, err = decodeHexField(rows[i].Chapter); err != nil {
			return nil, err
		}
		if rows[i].ChapterID, err = decodeHexField(rows[i].ChapterID); err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func bookmarkMatchKey(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func sentenceBoundaryPrefix(full, selected string) bool {
	if utf8.RuneCountInString(selected) < 8 || !strings.HasPrefix(full, selected) || len(full) == len(selected) {
		return false
	}
	last, _ := utf8.DecodeLastRuneInString(selected)
	next, _ := utf8.DecodeRuneInString(full[len(selected):])
	return strings.ContainsRune("。.!?！？", last) || strings.ContainsRune("。.!?！？", next)
}

func readMatchingBookmarkText(binary, lib, db string, book RecentBook, selection string) string {
	volumeHex := strings.ToUpper(hex.EncodeToString([]byte(book.Volume)))
	query := `SELECT hex(Text) AS text FROM Bookmark
 WHERE VolumeID=CAST(X'` + volumeHex + `' AS TEXT)
   AND Type IN ('highlight','note')
   AND length(trim(COALESCE(Text,'')))>0
   AND lower(COALESCE(CAST(Hidden AS TEXT),'false')) NOT IN ('1','true');`
	out, err := sqliteQuery(binary, lib, db, query)
	if err != nil {
		return ""
	}
	var rows []struct {
		Text string `json:"text"`
	}
	if json.Unmarshal(out, &rows) != nil {
		return ""
	}
	key := bookmarkMatchKey(selection)
	match, count := "", 0
	for _, row := range rows {
		full, err := decodeHexField(row.Text)
		if err != nil {
			return ""
		}
		fullKey := bookmarkMatchKey(full)
		if fullKey == key || sentenceBoundaryPrefix(fullKey, key) {
			match, count = full, count+1
		}
	}
	if count != 1 {
		return ""
	}
	return match
}
