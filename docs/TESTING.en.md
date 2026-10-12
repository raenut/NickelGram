# Local testing

For release changes and known limitations, see the [v0.1.3 release notes](RELEASE_NOTES_0.1.3.md).

The local tools use a read-only snapshot of your `KoboReader.sqlite` and the computer's `sqlite3` to run NickelGram's production Go queries and formatters. Debug commands do not create a Telegram client or call the Telegram API. The database, configuration, and generated output stay under Git-ignored `local-test/`.

## First run

Install `sqlite3` and Go 1.23 or later. Copy `KoboReader.sqlite` from a Kobo backup, then run these commands from the project root:

```sh
KOBO_DB_SOURCE=/path/to/KoboReader.sqlite ./test.sh setup
./test.sh list-books
```

`setup` copies a local snapshot and leaves the original database unchanged. Run it again after the source database changes. If `sqlite3` or Go is not on `PATH`, set `SQLITE3=/path/to/sqlite3` or `GO=/path/to/go`. The Kobo package uses its bundled ARM `sqlite3`; computer tests use the computer's `sqlite3` with the same query SQL.

```text
local-test/
├── KoboReader.sqlite       # private database snapshot; do not commit
├── footer.txt              # test footer when no config is specified
├── bin/nickelgram-debug    # computer build
├── cache/                  # Go build cache
└── output/                 # rendered messages, Markdown, and diagnostics
```

To preview with your own configuration, set `NICKELGRAM_TEST_CONFIG=/path/to/config.json`. The program parses the configuration and uses only the footer and tags for previews. It does not send Telegram requests or print the Bot Token. Without this variable, it uses `local-test/footer.txt`.

## Common commands

```sh
./test.sh list-books
./test.sh list-highlights
./test.sh latest-highlight '<book identifier from list-books>'
./test.sh render-message '<BookmarkID from list-highlights>'
./test.sh export-md '<book identifier from list-books>'
```

For `latest-highlight` and `export-md`, use `current` as the book identifier to select the book inferred by production code. `./test.sh current-book` prints that inferred book. `render-message` writes `output/message.txt`; `export-md` writes a Markdown file named with the title and a short identifier. The preview is ordered by book position and includes chapter headings and a timezone-free `exported_at`. On Kobo, exports are saved under `Highlights/` at the storage root.

To inspect paragraph breaks across a selection, run:

```sh
./test.sh render-selection current 'First paragraph

Second paragraph'
```

The `RAW_INPUT_RAW`, `RAW_TEXT_RAW`, and `FINAL_MESSAGE_RAW` values show line breaks as `\n`; numbered `LINE` entries show each line and blank line.

## Limits of local testing

Local tests can inspect the snapshot's SQLite schema, book metadata, highlight text, annotations, timestamps, footer, messages, Markdown, and filenames. A snapshot may differ from the device's current state, and multi-paragraph behavior depends on the supplied sample.

Local tests cannot confirm the text passed by NickelMenu's `{1}` placeholder, selection UI behavior, reader events, when Kobo writes bookmarks, reader-menu state, Kobo ARM runtime behavior, or actual Telegram delivery. These require device verification. If the input lacks paragraph markers, NickelGram cannot infer the book's paragraph boundaries. An annotation that has not been saved to the Kobo database cannot be retrieved by the read-only query.
