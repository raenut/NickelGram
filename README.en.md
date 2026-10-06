<p align="center">
  <img src="assets/nickelgram-icon.svg" alt="NickelGram icon" width="144" height="144">
</p>

# NickelGram

[简体中文](README.md) · [Download the latest release](https://github.com/raenut/NickelGram/releases/latest)

[![Latest release](https://img.shields.io/github/v/release/raenut/NickelGram?label=release&color=788c96)](https://github.com/raenut/NickelGram/releases/latest) [![License](https://img.shields.io/github/license/raenut/NickelGram?color=788c96)](LICENSE)

NickelGram adds Telegram sharing to Kobo e-readers through [NickelMenu](https://github.com/pgaskin/NickelMenu). Send selected text or the latest highlight, or export a book's highlights and annotations as a Markdown file.

## 🚀 Quick start

You need a Kobo e-reader with NickelMenu installed, a Telegram Bot Token, and the target chat's Chat ID. The release package targets 32-bit ARM Linux Kobo readers. Before installing, back up `.kobo/KoboReader.sqlite`, any existing NickelGram configuration, and your NickelMenu configuration.

1. Download the Kobo ARM ZIP from the [latest release](https://github.com/raenut/NickelGram/releases/latest).
2. Extract it and merge its `.adds` folder into the root of Kobo storage, keeping your other existing files.
3. Copy `.adds/nickelgram/config.example.json` to `config.json` in the same folder. Fill in `telegram_bot_token` and `telegram_chat_id`.
4. Safely eject and restart your Kobo. Use a NickelGram sharing command from the text selection or reader menu.

## 📖 Features

| Where | Menu entry | Result |
| --- | --- | --- |
| Text selection menu | **Send to Telegram** | Sends selected text with the book title and author |
| Reader menu | **Send Latest Highlight** | Sends the latest visible highlight or annotation for the current book |
| Reader menu | **Export All Highlights** | Exports the current book's highlights and annotations as Markdown and sends the file to Telegram |

Exported files are saved in `Highlights/` at the root of Kobo storage. Exporting the same book again uses the same filename.

## ⚙️ Configuration

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

`footer` controls the footer on Telegram messages; `md_footer` controls the footer in exported Markdown. Set either to `false` to turn it off. Add tags through `md_tags`, for example `["highlights", "literature"]`. Keep the device's `config.json` private because it contains your Bot Token and Chat ID.

## ⚠️ Limitations

- Immediately after switching books, sharing only selected text may show the previous book's title and author. Saving a highlight or annotation first helps identify the current book.
- Line breaks in multi-paragraph selections have not been fully verified on a device. Reopening an old highlight and sharing it through the text selection menu may produce incomplete text.
- Complete on-device sharing tests and compatibility checks across firmware versions are still pending. Testing so far has covered Kobo Libra Colour and Kobo Clara 2E. Back up relevant data before use.

Released under the [MIT License](LICENSE). See [Local Testing](docs/LOCAL_TEST.md) for development and testing, and [Third-party Components](docs/THIRD_PARTY.md) for bundled component details.
