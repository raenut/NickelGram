package main

import "strings"

const (
	defaultHeaderFormat    = "{title}, {author}"
	defaultQuoteLabel      = "❝"
	defaultAnnotationLabel = "批注："
)

func messageConfig(c Config) Config {
	if c.HeaderFormat == "" {
		c.HeaderFormat = defaultHeaderFormat
	}
	if c.QuoteLabel == "" {
		c.QuoteLabel = defaultQuoteLabel
	}
	if c.AnnotationLabel == "" {
		c.AnnotationLabel = defaultAnnotationLabel
	}
	return c
}

func bookHeader(book RecentBook, format string) string {
	if strings.TrimSpace(book.Author) == "" && format == defaultHeaderFormat {
		return book.Title
	}
	return strings.NewReplacer("{title}", book.Title, "{author}", book.Author).Replace(format)
}

func selectionMessage(book RecentBook, selection, footer string) string {
	return configuredSelectionMessage(book, selection, footer, Config{})
}

func configuredSelectionMessage(book RecentBook, selection, footer string, config Config) string {
	config = messageConfig(config)
	blocks := []string{bookHeader(book, config.HeaderFormat), config.QuoteLabel + "\n" + normalizeSelectionParagraphs(selection)}
	if strings.TrimSpace(footer) != "" {
		blocks = append(blocks, footer)
	}
	return strings.Join(blocks, "\n\n")
}

func highlightMessage(book RecentBook, mark Highlight, footer string) string {
	return configuredHighlightMessage(book, mark, footer, Config{})
}

func configuredHighlightMessage(book RecentBook, mark Highlight, footer string, config Config) string {
	config = messageConfig(config)
	blocks := []string{bookHeader(book, config.HeaderFormat), config.QuoteLabel + "\n" + normalizeSelectionParagraphs(mark.Text)}
	if strings.TrimSpace(mark.Annotation) != "" {
		blocks = append(blocks, config.AnnotationLabel+"\n"+mark.Annotation)
	}
	if strings.TrimSpace(footer) != "" {
		blocks = append(blocks, footer)
	}
	return strings.Join(blocks, "\n\n")
}
