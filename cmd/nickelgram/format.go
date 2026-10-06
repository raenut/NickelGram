package main

import "strings"

func selectionMessage(book RecentBook, selection, footer string) string {
	heading := book.Title
	if strings.TrimSpace(book.Author) != "" {
		heading += ", " + book.Author
	}
	blocks := []string{heading, "❝\n" + normalizeSelectionParagraphs(selection)}
	if strings.TrimSpace(footer) != "" {
		blocks = append(blocks, footer)
	}
	return strings.Join(blocks, "\n\n")
}

func highlightMessage(book RecentBook, mark Highlight, footer string) string {
	heading := book.Title
	if strings.TrimSpace(book.Author) != "" {
		heading += ", " + book.Author
	}
	blocks := []string{heading, "❝\n" + mark.Text}
	if strings.TrimSpace(mark.Annotation) != "" {
		blocks = append(blocks, "批注：\n"+mark.Annotation)
	}
	if strings.TrimSpace(footer) != "" {
		blocks = append(blocks, footer)
	}
	return strings.Join(blocks, "\n\n")
}
