<p align="center">
  <img src="assets/nickelgram-icon.svg" alt="NickelGram icon" width="144" height="144">
</p>

# NickelGram

[简体中文](README.md) · [Latest Release](https://github.com/raenut/NickelGram/releases/latest)

[![Latest Release](https://img.shields.io/github/v/release/raenut/NickelGram?label=release&color=788c96)](https://github.com/raenut/NickelGram/releases/latest) [![License](https://img.shields.io/github/license/raenut/NickelGram?color=788c96)](LICENSE)

NickelGram adds Telegram sharing to Kobo e-readers through [NickelMenu](https://github.com/pgaskin/NickelMenu).

Send selected text or the latest highlight, or export a book's highlights and annotations as a Markdown file.

## Supported devices

The current release package targets 32-bit ARM Linux Kobo readers.

## Quick start

### 1. Preparation

- Install NickelMenu on your Kobo e-reader.
- Prepare a Telegram Bot Token and the target chat's Chat ID.
- Back up `.kobo/KoboReader.sqlite`, any existing NickelGram configuration, and your NickelMenu configuration.


### 2. Installation

1. Download the archive ending in `-kobo-arm.zip` from the [Latest Release](https://github.com/raenut/NickelGram/releases/latest).
2. Extract the archive.
3. Merge its `.adds` folder into the root of Kobo storage, keeping your other existing files.
4. Copy `.adds/nickelgram/config.example.json` to `config.json` in the same folder.
5. Enter `telegram_bot_token` and `telegram_chat_id` in `config.json`.
6. Safely eject and restart your Kobo.

## Features

| Where | Menu entry | Result |
| --- | --- | --- |
| Text selection menu | Send to Telegram | Sends selected text with the book title and author |
| Reader menu | Send Latest Highlight | Sends the latest visible highlight or annotation for the current book |
| Reader menu | Export All Highlights | Exports the current book's highlights and annotations as Markdown and sends the file to Telegram |

Exported files are saved in `Highlights/` at the root of Kobo storage.

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

- Some selections or shares made immediately after switching books may be inaccurate.
- Compatibility across devices and firmware versions has not been fully verified. Testing so far has covered Kobo Libra Colour and Kobo Clara 2E.

Released under the [MIT License](LICENSE).

See [Local Testing](docs/LOCAL_TEST.md) for development and testing.

See [Third-party Components](docs/THIRD_PARTY.md) for bundled component details.
