# NickelGram

English | [中文](README.zh-CN.md)

NickelGram is a lightweight tool for Kobo e-readers and NickelMenu. It can send selected text, the latest highlight or annotation to Telegram, or export all highlights and annotations from a book as a Markdown file and send it to Telegram.

The project is still at an early stage. Features, compatibility, and usage may change in future releases.

## Features

### v0.1.1 update

- **Book order**: Export highlights in book content order and by their position within that content, rather than by creation time. Highlights without matching position data follow those with known positions.
- **Chapters and export time**: Use Markdown chapter headings that can be folded in editors which support heading folding, grouping consecutive highlights from the same content section while keeping an `<hr>` separator for each highlight. Show only the filename for HTML paths, such as `part0043.xhtml`. YAML front matter now includes `exported_at` to the second without a time zone; there is no `Highlights` heading.
- **Menu and location**: The Reader menu entry is now **Export All Highlights**. The `.md` file is saved in `Highlights/` at the root of Kobo storage and is still sent to Telegram. “All” means every visible highlight and annotation in the current book.

### v0.1.0

- **Share Selection**: Send the currently selected text as a regular Telegram message from the selection menu, together with the book title, author, and an optional `footer`.

- **Latest Highlight or Annotation**: Send the latest visible highlight or annotation detected for the current book from the Reader menu.

- **Export Book Highlights**: Export all detected visible highlights and annotations from the current book to a `.md` file and send it to Telegram. The exported Markdown includes YAML front matter, `<hr>` separators, configurable tags, and a separate Markdown `footer`.

## Installation

A Kobo e-reader with NickelMenu installed is required.

The current release package is built for 32-bit ARM Linux and has not been verified on other architectures. Before installation, it is recommended to back up `.kobo/KoboReader.sqlite`, any existing `.adds/nickelgram/config.json`, and your NickelMenu configuration.

1. Download `NickelGram-0.1.1-kobo-arm.zip` from the [v0.1.1 Release](https://github.com/raenut/NickelGram/releases/tag/v0.1.1).

2. Extract the archive and merge the included `.adds` directory into the root of the Kobo storage. Do not remove other existing contents inside `.adds`. The bundled `.adds/nm/nickelgram` adds three NickelMenu entries.

3. Copy `.adds/nickelgram/config.example.json` to `config.json` in the same directory, then enter your Telegram Bot Token and target Chat ID. Your actual configuration remains on the device and is not included in the Release package.

4. Safely eject the Kobo and restart it, then confirm that the NickelGram menu entries appear. The Telegram Bot must also have permission to send messages and files to the target chat or channel.

## Configuration

Example configuration:

```json
{
  "telegram_bot_token": "YOUR_BOT_TOKEN",
  "telegram_chat_id": "YOUR_CHAT_ID",
  "footer": true,
  "footer_text": "📖 From Kobo",
  "md_footer": true,
  "md_footer_text": "*📖 From Kobo*",
  "md_tags": []
}
```

The placeholder values above cannot be used directly.

`footer` controls the `footer` appended to regular Telegram messages, while `md_footer` independently controls the `footer` in exported Markdown files. Set the corresponding option to `false` to disable it.

`md_footer_text` may contain Markdown formatting, and `md_tags` can be configured with values such as `["highlights", "literature"]`.

The legacy `"footer": "text"` format is still supported for regular messages.

## Usage

| Entry | Action |
| --- | --- |
| Text selection menu | **Send to Telegram** — Send the selected text |
| Reader menu | **Send Latest Highlight** — Send the latest highlight or annotation |
| Reader menu | **Export All Highlights** — Export and send all visible highlights and annotations from the current book as `.md` |

Exported Markdown files are stored in:

```text
Highlights/
```

The filename consists of the book title and a short identifier derived from the Kobo `ContentID`. Re-exporting the same book uses the same filename.

For local testing and repeatable test commands, see [Local Testing](docs/LOCAL_TEST.md).

The Release also includes the source archive `NickelGram-0.1.1-source.zip` and `SHA256SUMS`.

## Known Issues & Limitations

- **Book metadata immediately after opening a new book**: Kobo's `DateLastRead` may still point to the previously opened book. If a new highlight or annotation has already been written to the database, NickelGram prioritizes the associated Bookmark `VolumeID` when identifying the book.

  However, when text is merely selected and has not yet been saved as a highlight or annotation, the message may still show the title and author of the previous book. The current-book context provided by the Reader menu also cannot be fully verified using only a local database snapshot.

- **Line breaks in selections**: NickelGram preserves the line breaks actually passed by NickelMenu. However, it has not yet been fully verified on-device whether Kobo's selection menu preserves all paragraph breaks when a selection spans multiple paragraphs.

- **Re-sharing existing highlights**: Re-sharing an existing highlight by opening it again is not currently considered a supported workflow. If the menu only passes part of a sentence, the resulting message may be truncated.

  The three currently supported workflows are regular text selection sharing, sharing a saved highlight or annotation from the Reader menu, and exporting all highlights from a book.

- **v0.1.1 on-device regression testing**: Local database queries, formatting logic, and the ARM build have been verified. However, the three complete menu workflows, actual Telegram delivery, and compatibility across different firmware versions still require further on-device verification for the [v0.1.1 Release](https://github.com/raenut/NickelGram/releases/tag/v0.1.1) package.

## Risk Disclaimer

NickelGram v0.1.1 is an early release and may still contain bugs. Possible issues include changes in Kobo database state, firmware or NickelMenu compatibility problems, network or Telegram delivery failures, incorrect book identification, and unexpected data loss during installation or use.

Back up your device data and configuration before installation and decide whether the software is appropriate for your own use. This project is provided "as is" under the terms described in [LICENSE](LICENSE), without warranty of fitness for any particular purpose.

## Tested Environment

The author's current on-device test environment is limited to:

- Kobo Libra Colour
- Kobo Clara 2E
- The latest firmware available for those devices at the time of testing
- NickelMenu 0.6.0

Other Kobo models, firmware versions, and NickelMenu versions have not yet been verified. The v0.1.1 Release package has also not yet completed full end-to-end regression testing on the devices listed above.

## Third-party Components & Licenses

See [THIRD_PARTY.md](docs/THIRD_PARTY.md) for information about third-party components and their licenses.
