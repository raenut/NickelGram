package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func safeMarkdownName(title, volume string) string {
	name := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || strings.ContainsRune(`/\:*?"<>|`, r) {
			return '_'
		}
		return r
	}, title)
	name = strings.Trim(name, " .")
	if name == "" {
		name = "book"
	}
	runes := []rune(name)
	if len(runes) > 100 {
		name = string(runes[:100])
	}
	id := sha256.Sum256([]byte(volume))
	return fmt.Sprintf("%s-%x.md", name, id[:4])
}

func yamlString(value string) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func bookMarkdown(book RecentBook, marks []Highlight, footer string, tags []string) string {
	var out strings.Builder
	out.WriteString("---\n")
	out.WriteString("title: " + yamlString(strings.TrimSpace(book.Title)) + "\n")
	out.WriteString("author: " + yamlString(strings.TrimSpace(book.Author)) + "\n")
	out.WriteString("source: kobo\n")
	cleanTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag = strings.TrimSpace(tag); tag != "" {
			cleanTags = append(cleanTags, tag)
		}
	}
	if len(cleanTags) == 0 {
		out.WriteString("tags: []\n")
	} else {
		out.WriteString("tags:\n")
		for _, tag := range cleanTags {
			out.WriteString("  - " + yamlString(tag) + "\n")
		}
	}
	out.WriteString("---\n")
	for _, mark := range marks {
		out.WriteString("\n<hr>\n\n")
		lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(mark.Text), "\r\n", "\n"), "\n")
		for _, line := range lines {
			out.WriteString("> " + line + "\n")
		}
		if note := strings.TrimSpace(mark.Annotation); note != "" {
			out.WriteString("\n批注：\n" + note + "\n")
		}
	}
	if footer = strings.TrimSpace(footer); footer != "" {
		out.WriteString("\n<hr>\n\n" + footer + "\n")
	}
	return out.String()
}

func runExportBook(root string) error {
	lock, err := acquireLock()
	if err != nil {
		return errors.New("另一个 NickelGram 操作正在运行。")
	}
	defer func() { _ = lock.Close(); _ = os.Remove(lockPath) }()
	c, err := readConfig(root)
	if err != nil {
		return err
	}
	tg, err := newTelegram(root)
	if err != nil {
		return err
	}
	if err = verifySQLiteRuntime(root); err != nil {
		return err
	}
	binary, lib := filepath.Join(root, "sqlite3"), filepath.Join(root, "lib")
	book, err := readCurrentBook(binary, lib, koboDB)
	if err != nil {
		return err
	}
	marks, err := readBookHighlights(binary, lib, koboDB, book)
	if err != nil {
		return err
	}
	if len(marks) == 0 {
		return errors.New("当前书没有可用的高亮或批注；未发送。")
	}
	dir := filepath.Join(root, "exports")
	if err = os.MkdirAll(dir, 0700); err != nil {
		return errors.New("无法创建导出目录；未发送。")
	}
	name := safeMarkdownName(book.Title, book.Volume)
	path := filepath.Join(dir, name)
	if err = os.WriteFile(path, []byte(bookMarkdown(book, marks, c.MDFooter, c.MDTags)), 0600); err != nil {
		return errors.New("无法写入 Markdown 文件；未发送。")
	}
	outcome := tg.SendDocument(c, path, name)
	if !outcome.OK {
		return errors.New(outcome.Message)
	}
	fmt.Println("发送成功：" + name)
	return nil
}
