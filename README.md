# NickelGram

[中文说明](README.zh-CN.md) · [Download the latest release](https://github.com/raenut/NickelGram/releases/latest) · [Release history](https://github.com/raenut/NickelGram/releases)

NickelGram adds Telegram sharing to Kobo e-readers through NickelMenu. Send selected text or the latest highlight as a message, or export a book's highlights and annotations as a Markdown file.

## What it does

| Where | Menu entry | Result |
| --- | --- | --- |
| Text selection menu | **Send to Telegram** | Sends the selected text with the book title and author |
| Reader menu | **Send Latest Highlight** | Sends the latest visible highlight or annotation for the current book |
| Reader menu | **Export All Highlights** | Saves the current book's visible highlights and annotations as `.md` and sends the file to Telegram |

The Markdown export follows book content order, then highlight position within each content section. Consecutive highlights from the same chapter share a Markdown heading that can be folded in supported editors; each highlight has an `<hr>` separator. If a chapter is named after an HTML file, only its filename is shown (for example, `part0043.xhtml`). The file also contains YAML front matter with an `exported_at` timestamp to the second without a time zone, optional tags, and an optional Markdown footer. Highlights without position data appear after those with known positions.

Exports are saved in `Highlights/` at the root of Kobo storage. The filename combines the book title with a short identifier derived from Kobo's `ContentID`; exporting the same book again uses the same filename.

## Install

You need a Kobo e-reader with [NickelMenu](https://github.com/pgaskin/NickelMenu) installed. The release package is built for 32-bit ARM Linux; other architectures have not been verified. Back up `.kobo/KoboReader.sqlite`, any existing `.adds/nickelgram/config.json`, and your NickelMenu configuration before installing.

1. Download the Kobo ARM ZIP from the [latest release](https://github.com/raenut/NickelGram/releases/latest). The release also provides a source archive and `SHA256SUMS`.
2. Extract the ZIP and merge its `.adds` directory into the root of Kobo storage. Keep the other contents of your existing `.adds` directory. The included `.adds/nm/nickelgram` adds the three menu entries above.
3. Copy `.adds/nickelgram/config.example.json` to `config.json` in the same directory. Add your Telegram Bot Token and target Chat ID. Your actual configuration stays on the device and is not included in the release package.
4. Safely eject and restart the Kobo. Check that the menu entries appear and that the bot can send messages and files to the target chat or channel.

## Configure

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

Replace the placeholder token and chat ID. `footer` controls the footer on Telegram messages; `md_footer` controls the footer in exported Markdown. Set either to `false` to turn it off. `md_footer_text` supports Markdown formatting, and `md_tags` accepts a list such as `["highlights", "literature"]`. The older `"footer": "text"` form remains supported for messages.

## Limitations

- **Book identification:** Immediately after opening a different book, Kobo's `DateLastRead` can still point to the previous one. A saved highlight or annotation helps NickelGram identify the current book through its Bookmark `VolumeID`; sharing only a text selection may show the previous book's title and author.
- **Selection text:** NickelGram preserves the line breaks NickelMenu passes to it, but paragraph breaks in multi-paragraph selections have not been fully verified on a device. Reopening an old highlight and sharing it through the selection menu is not a supported workflow; the passed text may be incomplete.
- **Device testing:** Database queries, formatting, and the ARM build have been checked locally. Complete on-device tests of all three menu workflows and Telegram delivery, as well as compatibility across firmware versions, are still pending.

Testing so far has been limited to Kobo Libra Colour and Kobo Clara 2E, with the firmware available at the time of testing and NickelMenu 0.6.0. Other devices and versions have not been verified. NickelGram is an early project; back up your device data and configuration before using it. It is provided under the [LICENSE](LICENSE) terms, without warranty.

For repeatable computer-side tests, see [Local Testing](docs/LOCAL_TEST.md). For bundled components and licenses, see [Third-party Components](docs/THIRD_PARTY.md).
