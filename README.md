<p align="center">
  <img src="assets/nickelgram-icon.svg" alt="NickelGram 图标" width="144" height="144">
</p>

# NickelGram

[English](README.en.md) · [下载最新版本](https://github.com/raenut/NickelGram/releases/latest)

[![最新版本](https://img.shields.io/github/v/release/raenut/NickelGram?label=release&color=788c96)](https://github.com/raenut/NickelGram/releases/latest) [![许可证](https://img.shields.io/github/license/raenut/NickelGram?color=788c96)](LICENSE)

NickelGram 通过 [NickelMenu](https://github.com/pgaskin/NickelMenu) 为 Kobo 阅读器添加 Telegram 分享功能：发送选中文字或最近一条书摘，也可以把整本书的高亮和批注导出为 Markdown 文件。

## 🚀 快速使用

使用前需要一台已安装 NickelMenu 的 Kobo 阅读器，以及可用的 Telegram Bot Token 和目标聊天的 Chat ID。发布包适用于 32 位 ARM Linux Kobo 阅读器。安装前请备份 `.kobo/KoboReader.sqlite`、已有的 NickelGram 配置和 NickelMenu 配置。

1. 从[最新发布版本](https://github.com/raenut/NickelGram/releases/latest)下载 Kobo ARM ZIP。
2. 解压后，将其中的 `.adds` 文件夹合并到 Kobo 存储根目录，保留设备上已有的其他文件。
3. 将 `.adds/nickelgram/config.example.json` 复制为同目录下的 `config.json`，填入 `telegram_bot_token` 和 `telegram_chat_id`。
4. 安全弹出并重启 Kobo。从选中文字后的菜单或阅读器菜单中选择 NickelGram 的分享命令。

## 📖 功能

| 入口 | 菜单项 | 作用 |
| --- | --- | --- |
| 选中文字后的菜单 | **Send to Telegram** | 连同书名、作者发送选中文字 |
| 阅读器菜单 | **Send Latest Highlight** | 发送当前书籍最近一条可见高亮或批注 |
| 阅读器菜单 | **Export All Highlights** | 导出当前书籍的高亮与批注为 Markdown，并将文件发送到 Telegram |

导出的文件保存在 Kobo 存储根目录的 `Highlights/` 中；再次导出同一本书会使用相同的文件名。

## ⚙️ 配置

```json
{
  "telegram_bot_token": "YOUR_BOT_TOKEN",
  "telegram_chat_id": "YOUR_CHAT_ID",
  "footer": true,
  "footer_text": "📖 来自 Kobo",
  "md_footer": true,
  "md_footer_text": "*📖 摘自 Kobo*",
  "md_tags": []
}
```

`footer` 控制 Telegram 消息的页脚，`md_footer` 控制导出文件的 Markdown 页脚；将对应值设为 `false` 即可关闭。`md_tags` 可填写标签，例如 `["书摘", "文学"]`。Bot Token 和 Chat ID 只需写入设备上的 `config.json`，不要公开分享该文件。

## ⚠️ 使用限制

- 刚切换书籍后，仅分享选中文字时，书名和作者偶尔可能来自上一本书。先保存一条高亮或批注有助于正确识别当前书籍。
- 跨段落选中文字的换行效果尚未在设备上充分验证；重新打开旧高亮后再通过文字选区分享，可能得到不完整内容。
- 完整的设备端分享流程和不同固件版本的兼容性仍待进一步验证。目前测试设备为 Kobo Libra Colour 和 Kobo Clara 2E。使用前请备份相关数据。

项目按 [MIT 许可证](LICENSE) 提供。开发与本地测试说明见[本地测试文档](docs/LOCAL_TEST.md)，第三方组件信息见[第三方组件文档](docs/THIRD_PARTY.md)。
