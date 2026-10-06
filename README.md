<p align="center">
  <img src="assets/nickelgram-icon.svg" alt="NickelGram 图标" width="144" height="144">
</p>

# NickelGram

[English](README.en.md) · [Latest Release](https://github.com/raenut/NickelGram/releases/latest)

[![最新版本](https://img.shields.io/github/v/release/raenut/NickelGram?label=release&color=788c96)](https://github.com/raenut/NickelGram/releases/latest) [![许可证](https://img.shields.io/github/license/raenut/NickelGram?color=788c96)](LICENSE)

NickelGram 通过 [NickelMenu](https://github.com/pgaskin/NickelMenu) 为 Kobo 阅读器添加 Telegram 分享功能。

可以发送选中文字或最近一条书摘，也可以把整本书的高亮和批注导出为 Markdown 文件。

## 适用设备

当前发布包适用于 32 位 ARM Linux Kobo 阅读器。

## 快速使用

### 1. 准备工作

- 在 Kobo 阅读器上安装 NickelMenu。
- 准备 Telegram Bot Token 和目标聊天的 Chat ID。
- 备份 `.kobo/KoboReader.sqlite`、已有的 NickelGram 配置和 NickelMenu 配置。


### 2. 正式安装

1. 从 [Latest Release](https://github.com/raenut/NickelGram/releases/latest) 下载文件名以 `-kobo-arm.zip` 结尾的压缩包。
2. 解压压缩包。
3. 将其中的 `.adds` 文件夹合并到 Kobo 存储根目录，保留设备上已有的其他文件。
4. 将 `.adds/nickelgram/config.example.json` 复制为同目录下的 `config.json`。
5. 在 `config.json` 中填入 `telegram_bot_token` 和 `telegram_chat_id`。
6. 安全弹出并重启 Kobo。

## 功能

| 入口 | 菜单项 | 作用 |
| --- | --- | --- |
| 选中文字后的菜单 | Send to Telegram | 连同书名、作者发送选中文字 |
| 阅读器菜单 | Send Latest Highlight | 发送当前书籍最近一条可见高亮或批注 |
| 阅读器菜单 | Export All Highlights | 导出当前书籍的高亮与批注为 Markdown，并将文件发送到 Telegram |

导出的文件保存在 Kobo 存储根目录的 `Highlights/` 中。

## 配置

在 `.adds/nickelgram/config.json` 中设置以下参数：

| 参数 | 说明 | 示例 |
| --- | --- | --- |
| `telegram_bot_token` | Telegram Bot Token | `"YOUR_BOT_TOKEN"` |
| `telegram_chat_id` | 目标聊天的 Chat ID | `"YOUR_CHAT_ID"` |
| `footer` | 是否在 Telegram 消息中添加页脚 | `true` |
| `footer_text` | Telegram 消息的页脚文字 | `"📖 来自 Kobo"` |
| `md_footer` | 是否在导出的 Markdown 中添加页脚 | `true` |
| `md_footer_text` | Markdown 页脚文字 | `"*📖 摘自 Kobo*"` |
| `md_tags` | 导出文件的标签 | `["书摘", "文学"]` |

请勿公开分享包含 Bot Token 和 Chat ID 的 `config.json`。

## 使用限制

- 部分选区或刚切换书籍后的分享结果可能不准确。
- 不同设备和固件的兼容性尚未全面验证。目前测试设备为 Kobo Libra Colour 和 Kobo Clara 2E。

项目按 [MIT 许可证](LICENSE) 提供。

开发与本地测试说明见[本地测试文档](docs/LOCAL_TEST.md)。

第三方组件信息见[第三方组件文档](docs/THIRD_PARTY.md)。
