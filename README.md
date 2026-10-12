<p align="center">
  <img src="assets/nickelgram-icon.svg" alt="NickelGram icon" width="144" height="144">
</p>

# NickelGram

[简体中文](README.zh-CN.md) ·
[Download v0.1.4](https://github.com/raenut/NickelGram/releases/tag/v0.1.4)

NickelGram adds Telegram sharing to Kobo readers through
[NickelMenu](https://github.com/pgaskin/NickelMenu).

Send selected text or the latest highlight. Export a book's highlights and
annotations as Markdown.

## Features

| Where | Menu entry | Result |
| --- | --- | --- |
| Selection | Send&nbsp;to&nbsp;Telegram | Send selected text with book details |
| Reader | Send&nbsp;Latest&nbsp;Highlight | Send the latest highlight or annotation |
| Reader | Export&nbsp;All&nbsp;Highlights | Export and send all highlights as Markdown |

Markdown files are saved in `Highlights/` at the root of Kobo storage.

## Install

- Install NickelMenu.
- Create a Telegram bot and get the target chat's Chat ID.
- Back up `.kobo/KoboReader.sqlite` and existing NickelGram/NickelMenu settings.

1. Download `NickelGram-0.1.4-kobo-arm.zip` from the
   [v0.1.4 release](https://github.com/raenut/NickelGram/releases).
2. Extract it and merge `.adds/` into the root of Kobo storage.
3. Copy `config.example.json` to `config.json` in `.adds/nickelgram/`.
4. Add your Bot Token and Chat ID, then safely eject and restart the Kobo.

The package targets 32-bit ARM Linux Kobo readers.

## Configuration

Edit `.adds/nickelgram/config.json`:

| Parameter | Purpose |
| --- | --- |
| `telegram_bot_token` | Telegram Bot Token |
| `telegram_chat_id` | Target chat |
| `footer` | Add a Telegram footer |
| `footer_text` | Footer text |
| `footer_url` | Optional footer link; previews stay off |
| `footer_link_text` | Clickable footer text; required with a URL |
| `header_format` | Set title/author order and punctuation |
| `quote_label` | Excerpt label; excerpt uses monospace |
| `annotation_label` | Annotation label |
| `md_footer` | Add a Markdown footer |
| `md_footer_text` | Markdown footer text |
| `md_tags` | Exported file tags |

`header_format` supports only `{title}` and `{author}` and must include
`{title}`. They use Kobo's book title and author; the author may be empty, so
literal punctuation remains. Templates and labels cannot contain newlines.
Older configurations keep working with default values. Keep `config.json`
private; it contains your Token and Chat ID. See
`.adds/nickelgram/config.example.json` for defaults.

## Limits

- Direct selections may lose paragraph breaks. Save a highlight to use Kobo's
  stored text when NickelGram finds one unique match in the current book.
- Markdown chapter names use Kobo's table of contents. Missing titles appear
  as “Untitled Chapter”.
- v0.1.4 was installed and launch-checked on Kobo Clara 2E. Telegram delivery
  and all reading flows have not been fully verified. Other devices and firmware
  remain unverified.

## Release notes

- [v0.1.4](docs/RELEASE_NOTES_0.1.4.md)
- [v0.1.3](docs/RELEASE_NOTES_0.1.3.md)
- [v0.1.2](docs/RELEASE_NOTES_0.1.2.md)
- [v0.1.1](docs/RELEASE_NOTES_0.1.1.md)
- [v0.1.0](docs/RELEASE_NOTES_0.1.0.md)

See [local testing](docs/TESTING.en.md) and
[third-party licenses](docs/THIRD_PARTY_LICENSES.md).

## License

MIT License. Copyright © 2026 NickelGram contributors.
