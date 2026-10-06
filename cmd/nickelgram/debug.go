package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// runDebug is only entered through the explicit local debug command. It never
// reads Telegram credentials or constructs a Telegram client.
func runDebug(args []string) error {
	if len(args) == 0 {
		return errors.New("需要命令：list-books、list-highlights、latest-highlight、render-message、export-md")
	}
	db, binary, outDir := os.Getenv("NICKELGRAM_TEST_DB"), os.Getenv("NICKELGRAM_TEST_SQLITE"), os.Getenv("NICKELGRAM_TEST_OUTPUT")
	if db == "" || binary == "" || outDir == "" {
		return errors.New("请通过 ./test.sh 运行本地测试")
	}
	footer := os.Getenv("NICKELGRAM_TEST_FOOTER")
	mdFooter := footer
	var mdTags []string
	if configPath := os.Getenv("NICKELGRAM_TEST_CONFIG"); configPath != "" {
		c, err := readConfigFile(configPath)
		if err != nil {
			return err
		}
		footer = c.Footer
		mdFooter, mdTags = c.MDFooter, c.MDTags
	}
	query := func(sql string, dest any) error {
		data, err := sqliteQuery(binary, "", db, sql)
		if err != nil {
			return err
		}
		if len(strings.TrimSpace(string(data))) == 0 {
			data = []byte("[]")
		}
		return json.Unmarshal(data, dest)
	}
	bookByID := func(id string) (RecentBook, error) {
		if id == "current" {
			return readCurrentBook(binary, "", db)
		}
		var rows []RecentBook
		key := strings.ToUpper(hex.EncodeToString([]byte(id)))
		err := query(`SELECT hex(ContentID) AS volume, hex(COALESCE(Title,'')) AS title,
 hex(COALESCE(Attribution,'')) AS author FROM content
 WHERE CAST(ContentType AS TEXT)='6' AND hex(ContentID)='`+key+`' LIMIT 1;`, &rows)
		if err != nil || len(rows) != 1 {
			return RecentBook{}, errors.New("找不到书籍标识")
		}
		book := rows[0]
		if book.Volume, err = decodeHexField(book.Volume); err != nil {
			return RecentBook{}, err
		}
		if book.Title, err = decodeHexField(book.Title); err != nil {
			return RecentBook{}, err
		}
		if book.Author, err = decodeHexField(book.Author); err != nil {
			return RecentBook{}, err
		}
		return book, nil
	}
	markByID := func(id string) (Highlight, error) {
		var rows []Highlight
		key := strings.ToUpper(hex.EncodeToString([]byte(id)))
		err := query(`SELECT `+highlightColumns+` FROM Bookmark WHERE hex(BookmarkID)='`+key+`'
 AND `+visibleHighlightCondition+` LIMIT 1;`, &rows)
		if err != nil || len(rows) != 1 {
			return Highlight{}, errors.New("找不到可见的高亮或批注标识")
		}
		return decodeHighlight(rows[0])
	}
	if err := os.MkdirAll(outDir, 0700); err != nil {
		return err
	}
	switch args[0] {
	case "list-books":
		if len(args) != 1 {
			return errors.New("list-books 不需要参数")
		}
		var rows []struct{ Volume, Title, Author, LastRead string }
		err := query(`SELECT ContentID AS volume, COALESCE(Title,'') AS title,
 COALESCE(Attribution,'') AS author, COALESCE(DateLastRead,'') AS lastread
 FROM content WHERE CAST(ContentType AS TEXT)='6'
 ORDER BY DateLastRead DESC, rowid DESC;`, &rows)
		if err != nil {
			return err
		}
		for _, row := range rows {
			fmt.Printf("%s\t%s\t%s\t%s\n", row.Volume, row.Title, row.Author, row.LastRead)
		}
	case "list-highlights":
		if len(args) > 2 {
			return errors.New("用法：list-highlights [书籍标识]")
		}
		filter := ""
		if len(args) == 2 {
			filter = ` AND hex(VolumeID)='` + strings.ToUpper(hex.EncodeToString([]byte(args[1]))) + `'`
		}
		var rows []Highlight
		err := query(`SELECT `+highlightColumns+` FROM Bookmark WHERE `+visibleHighlightCondition+filter+`
 ORDER BY DateCreated, BookmarkID;`, &rows)
		if err != nil {
			return err
		}
		for _, encoded := range rows {
			mark, err := decodeHighlight(encoded)
			if err != nil {
				return err
			}
			fmt.Printf("%s\t%s\t%s\t%s\n", mark.ID, mark.Volume, mark.Modified, oneLine(mark.Text))
		}
	case "latest-highlight":
		if len(args) != 2 {
			return errors.New("用法：latest-highlight <书籍标识>")
		}
		book, err := bookByID(args[1])
		if err != nil {
			return err
		}
		mark, err := readLatestHighlight(binary, "", db, book)
		if err != nil {
			return err
		}
		printFields(book, mark, footer)
	case "current-book":
		if len(args) != 1 {
			return errors.New("current-book 不需要参数")
		}
		book, err := readCurrentBook(binary, "", db)
		if err != nil {
			return err
		}
		fmt.Printf("BOOK_ID=%s\nTITLE=%s\nAUTHOR=%s\n", book.Volume, book.Title, book.Author)
	case "render-message":
		if len(args) != 2 {
			return errors.New("用法：render-message <BookmarkID>")
		}
		mark, err := markByID(args[1])
		if err != nil {
			return err
		}
		book, err := bookByID(mark.Volume)
		if err != nil {
			return err
		}
		message := highlightMessage(book, mark, footer)
		path := filepath.Join(outDir, "message.txt")
		if err := os.WriteFile(path, []byte(message), 0600); err != nil {
			return err
		}
		printFields(book, mark, footer)
		printTextDebug("RAW_TEXT", mark.Text)
		printTextDebug("RAW_ANNOTATION", mark.Annotation)
		printTextDebug("FINAL_MESSAGE", message)
		fmt.Println("FILE=" + path)
	case "render-selection":
		if len(args) != 3 {
			return errors.New("用法：render-selection <书籍标识|current> <选中文本>")
		}
		book, err := bookByID(args[1])
		if err != nil {
			return err
		}
		selection, err := validateSelection(args[2])
		if err != nil {
			return err
		}
		printTextDebug("RAW_INPUT", selection)
		if full := readMatchingBookmarkText(binary, "", db, book, selection); full != "" {
			selection = full
		}
		message := selectionMessage(book, selection, footer)
		path := filepath.Join(outDir, "selection-message.txt")
		if err := os.WriteFile(path, []byte(message), 0600); err != nil {
			return err
		}
		fmt.Printf("BOOK_ID=%s\nTITLE=%s\nAUTHOR=%s\nTEXT=%s\nFOOTER=%s\n", book.Volume, book.Title, book.Author, oneLine(selection), oneLine(footer))
		printTextDebug("FINAL_MESSAGE", message)
		fmt.Println("FILE=" + path)
	case "export-md":
		if len(args) != 2 {
			return errors.New("用法：export-md <书籍标识>")
		}
		book, err := bookByID(args[1])
		if err != nil {
			return err
		}
		marks, err := readBookHighlights(binary, "", db, book)
		if err != nil {
			return err
		}
		if len(marks) == 0 {
			return errors.New("这本书没有可见的高亮或批注")
		}
		name := safeMarkdownName(book.Title, book.Volume)
		path := filepath.Join(outDir, name)
		content := bookMarkdown(book, marks, mdFooter, mdTags, time.Now())
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			return err
		}
		fmt.Printf("TITLE=%s\nAUTHOR=%s\nMD_FOOTER=%s\nFILENAME=%s\nFILE=%s\n", book.Title, book.Author, mdFooter, name, path)
		printTextDebug("FINAL_MARKDOWN", content)
	default:
		return errors.New("未知本地测试命令")
	}
	return nil
}

func oneLine(s string) string {
	return strings.NewReplacer("\r", "\\r", "\n", "\\n", "\t", "\\t").Replace(s)
}

func printFields(book RecentBook, mark Highlight, footer string) {
	fmt.Printf("BOOK_ID=%s\nBOOKMARK_ID=%s\nTITLE=%s\nAUTHOR=%s\nTEXT=%s\nANNOTATION=%s\nDATE=%s\nFOOTER=%s\n", book.Volume, mark.ID, book.Title, book.Author, oneLine(mark.Text), oneLine(mark.Annotation), mark.Modified, oneLine(footer))
}

func printTextDebug(label, value string) {
	fmt.Printf("%s_RAW=%q\n", label, value)
	for i, line := range strings.Split(value, "\n") {
		fmt.Printf("%s_LINE_%d=%q\n", label, i+1, line)
	}
}
