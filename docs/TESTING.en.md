# Local testing

See the [release notes](RELEASE_NOTES_0.1.4.md) for changes and limits. Local
tools use a read-only `KoboReader.sqlite` copy and the computer's `sqlite3` to
run NickelGram's production Go queries and formatters. Debug commands do not
connect to Telegram. Data and output stay under Git-ignored `local-test/`.

## Setup

Install `sqlite3` and Go 1.23 or later. Copy `KoboReader.sqlite` from a Kobo
backup, then run these commands from the project root:

```sh
KOBO_DB_SOURCE=/path/to/KoboReader.sqlite ./test.sh setup
./test.sh list-books
```

`setup` copies the database and leaves the source unchanged. Run it again when
the source changes. Set `SQLITE3=/path/to/sqlite3` or `GO=/path/to/go` if either
tool is not on `PATH`. Kobo uses the bundled ARM `sqlite3`; local testing uses
the computer version with the same SQL queries.

## Common commands

```sh
./test.sh list-books
./test.sh list-highlights
./test.sh current-book
./test.sh latest-highlight '<book identifier>'
./test.sh render-message '<BookmarkID>'
./test.sh export-md '<book identifier>'
```

Use `current` as the book identifier when supported. `render-message` writes
`local-test/output/message.txt`. `export-md` writes a Markdown file ordered by
book position, with chapter headings and `exported_at`. Kobo exports go to
`Highlights/` at the storage root.

To inspect paragraph breaks in a selection:

```sh
./test.sh render-selection current 'First paragraph

Second paragraph'
```

`RAW_INPUT_RAW`, `RAW_TEXT_RAW`, and `FINAL_MESSAGE_RAW` display line breaks as
`\n`. Numbered `LINE` entries show each text line and blank line.

Set `NICKELGRAM_TEST_CONFIG=/path/to/config.json` to preview with a configuration.
The tool reads settings but does not print the Token or call Telegram.

## Test limits

Local checks can inspect the database schema, book details, highlights,
annotations, timestamps, message formatting, Markdown, and filenames. The
snapshot may differ from the device, and paragraph behavior depends on the
book and supplied sample.

Before release, check title/author templates, a missing author, monospace
excerpts, custom labels, and linked and plain-text footers. Also check regular
selections, saved highlights, the latest highlight, and full-book exports on
Kobo.

Local checks cannot confirm NickelMenu's `{1}` input, reader UI behavior, the
Kobo ARM runtime, or Telegram delivery. Without paragraph markers, the program
cannot infer the book's paragraph breaks. An annotation absent from the database
cannot be retrieved.
