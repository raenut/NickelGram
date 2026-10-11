package main

import (
	"archive/zip"
	"encoding/hex"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
)

type tocEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// A Kobo content row can contain many TOC chapters. Resolve each bookmark's
// Kobo span to the preceding TOC anchor in that content file.
func resolveChapterTitles(binary, lib, db string, book RecentBook, marks []Highlight) {
	dbPath := strings.TrimPrefix(db, "file:")
	dbPath, _, _ = strings.Cut(dbPath, "?")
	volumeHex := strings.ToUpper(hex.EncodeToString([]byte(book.Volume)))
	out, err := sqliteQuery(binary, lib, db, `SELECT hex(ContentID) AS id, hex(COALESCE(Title,'')) AS title
 FROM content WHERE hex(BookID)='`+volumeHex+`' AND CAST(ContentType AS TEXT)='899';`)
	if err != nil {
		return
	}
	var encoded []tocEntry
	if json.Unmarshal(out, &encoded) != nil {
		return
	}
	toc := make([]tocEntry, 0, len(encoded))
	for _, row := range encoded {
		id, idErr := decodeHexField(row.ID)
		title, titleErr := decodeHexField(row.Title)
		if idErr == nil && titleErr == nil && title != "" {
			toc = append(toc, tocEntry{ID: id, Title: title})
		}
	}
	archive, err := zip.OpenReader(filepath.Join(filepath.Dir(dbPath), "kepub", book.Volume))
	if err != nil {
		return
	}
	defer archive.Close()
	files := make(map[string]*zip.File, len(archive.File))
	for _, file := range archive.File {
		files[file.Name] = file
	}
	cache := map[string]string{}
	for i := range marks {
		base, _, ok := strings.Cut(marks[i].ChapterID, "!!")
		if !ok || base != book.Volume {
			continue
		}
		member := strings.SplitN(strings.SplitN(marks[i].ChapterID, "!!", 2)[1], "#", 2)[0]
		fragment := strings.ReplaceAll(marks[i].StartPath, `\.`, ".")
		if n := strings.LastIndexByte(fragment, '#'); n >= 0 {
			fragment = fragment[n+1:]
		}
		if fragment == "" {
			continue
		}
		content, cached := cache[member]
		if !cached {
			file := files[member]
			if file == nil || file.UncompressedSize64 > 4<<20 {
				continue
			}
			reader, openErr := file.Open()
			if openErr != nil {
				continue
			}
			bytes, readErr := io.ReadAll(io.LimitReader(reader, 4<<20))
			reader.Close()
			if readErr != nil {
				continue
			}
			content = string(bytes)
			cache[member] = content
		}
		position := htmlIDPosition(content, fragment)
		if position < 0 {
			continue
		}
		bestPosition, best := -1, tocEntry{}
		for _, entry := range toc {
			prefix := book.Volume + "!!" + member + "#"
			if !strings.HasPrefix(entry.ID, prefix) {
				continue
			}
			anchor := strings.TrimPrefix(entry.ID, prefix)
			anchorPosition := htmlIDPosition(content, anchor)
			if anchorPosition < 0 {
				if n := strings.LastIndexByte(anchor, '-'); n >= 0 {
					anchorPosition = htmlIDPosition(content, anchor[:n])
				}
			}
			if anchorPosition <= position && anchorPosition > bestPosition {
				bestPosition, best = anchorPosition, entry
			}
		}
		if bestPosition >= 0 {
			marks[i].Chapter = best.Title
			marks[i].ChapterID = best.ID
		}
	}
}

func htmlIDPosition(content, id string) int {
	if id == "" {
		return -1
	}
	if pos := strings.Index(content, `id="`+id+`"`); pos >= 0 {
		return pos
	}
	return strings.Index(content, "id='"+id+"'")
}
