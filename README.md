<p align="center">
  <img src="assets/nickelgram-icon.svg" alt="NickelGram icon" width="144" height="144">
</p>

# NickelGram

[简体中文](README.zh-CN.md) · [Latest Release](https://github.com/raenut/NickelGram/releases/tag/v0.1.2)

[![Latest Release](https://img.shields.io/github/v/release/raenut/NickelGram?label=release&color=788c96)](https://github.com/raenut/NickelGram/releases/latest) [![License](https://img.shields.io/github/license/raenut/NickelGram?color=788c96)](LICENSE)

NickelGram adds Telegram sharing to Kobo e-readers through [NickelMenu](https://github.com/pgaskin/NickelMenu).

Send selected text or the latest highlight, or export a book's highlights and annotations as a Markdown file.

## Features

| Where | Menu entry | Result |
| --- | --- | --- |
| Selection | Send&nbsp;to&nbsp;Telegram | Sends selected text with the book title and author |
| Reader | Send&nbsp;Latest&nbsp;Highlight | Sends the latest highlight or annotation |
| Reader | Export&nbsp;All&nbsp;Highlights | Exports and sends all highlights and annotations as Markdown |

Exported files are saved in `Highlights/` at the root of Kobo storage.

## Supported Devices

The current release package targets 32-bit ARM Linux Kobo readers.

## Install

- Install NickelMenu on your Kobo e-reader.
- Prepare a Telegram Bot Token and the target chat's Chat ID.
- Back up `.kobo/KoboReader.sqlite`, any existing NickelGram / NickelMenu configuration.

1. Download `NickelGram-0.1.2-kobo-arm.zip` from the [v0.1.2 release](https://github.com/raenut/NickelGram/releases/tag/v0.1.2) and extract it.
2. Merge the extracted `.adds` folder into the root of Kobo storage, keeping your other existing files.
3. Copy `.adds/nickelgram/config.example.json` to `config.json` in the same folder.
4. Edit `config.json` using the configuration table below.
5. Safely eject and restart your Kobo.

## Configuration

Set these parameters in `.adds/nickelgram/config.json`:

| Parameter | Description | Example |
| --- | --- | --- |
| `telegram_bot_token` | Telegram Bot Token | `"YOUR_BOT_TOKEN"` |
| `telegram_chat_id` | Target chat's Chat ID | `"YOUR_CHAT_ID"` |
| `footer` | Add a footer to Telegram messages | `true` |
| `footer_text` | Telegram message footer | `"📖 From Kobo"` |
| `md_footer` | Add a footer to exported Markdown | `true` |
| `md_footer_text` | Markdown footer text | `"*📖 From Kobo*"` |
| `md_tags` | Tags for exported files | `["highlights", "literature"]` |

Keep `config.json` private because it contains your Bot Token and Chat ID.

## Limitations

- Direct selections spanning paragraphs can lose paragraph breaks or carry a leading space from Kobo/NickelMenu. NickelGram cannot reliably reconstruct breaks absent from the selected text. Saved highlights use Kobo's stored text when matched.
- Chapter names in Markdown depend on Kobo's table of contents and the book file being readable. If an anchor cannot be matched, the heading falls back to a file name or “Unknown Chapter”. Chapter resolution reads the relevant book section; it does not reconstruct selected text.
- Tested on Kobo Libra Colour and Kobo Clara 2E. The final v0.1.2 package has not had a complete device retest; other devices and firmware versions remain unverified.

See [Testing](docs/TESTING.md) for development and testing.

See [Third-Party Licenses](docs/THIRD_PARTY_LICENSES.md) for bundled component details.

## License

Licensed under the [MIT License](LICENSE). Copyright © 2026 NickelGram contributors.

---

Made with 🤍 by [raenut](https://github.com/raenut) 👩🏻‍💻
