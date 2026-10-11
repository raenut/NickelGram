<p align="center">
  <img src="assets/nickelgram-icon.svg" alt="NickelGram 图标" width="144" height="144">
</p>

# NickelGram

[English](README.md) · [下载 v0.1.2](https://github.com/raenut/NickelGram/releases/tag/v0.1.2)

[![最新版本](https://img.shields.io/github/v/release/raenut/NickelGram?label=release&color=788c96)](https://github.com/raenut/NickelGram/releases/latest) [![许可证](https://img.shields.io/github/license/raenut/NickelGram?color=788c96)](LICENSE)

NickelGram 通过 [NickelMenu](https://github.com/pgaskin/NickelMenu) 为 Kobo 阅读器添加 Telegram 分享功能。

可以发送选中文字或最近一条书摘，也可以把整本书的高亮和批注导出为 Markdown 文件。

## 功能

| 入口 | 菜单项 | 作用 |
| --- | --- | --- |
| 文字选区 | Send&nbsp;to&nbsp;Telegram | 发送选中文字及书名、作者 |
| 阅读菜单 | Send&nbsp;Latest&nbsp;Highlight | 发送最近一条高亮或批注 |
| 阅读菜单 | Export&nbsp;All&nbsp;Highlights | 导出并发送全书高亮与批注（Markdown） |

导出的文件保存在 Kobo 存储根目录的 `Highlights/` 中。

## 适用设备

当前发布包适用于 32 位 ARM Linux Kobo 阅读器。

## 安装

- 在 Kobo 阅读器上安装 NickelMenu。
- 准备 Telegram Bot Token 和目标聊天的 Chat ID。
- 备份 `.kobo/KoboReader.sqlite`、已有的 NickelGram 配置和 NickelMenu 配置。

1. 从 [v0.1.2 发布页](https://github.com/raenut/NickelGram/releases/tag/v0.1.2) 下载 `NickelGram-0.1.2-kobo-arm.zip` 并解压。
2. 将解压得到的 `.adds` 文件夹合并到 Kobo 存储根目录，保留设备上已有的其他文件。
3. 将 `.adds/nickelgram/config.example.json` 复制为同目录下的 `config.json`。
4. 参照下方配置表修改 `config.json`。
5. 安全弹出并重启 Kobo。

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

- 直接选中跨段文字时，Kobo/NickelMenu 传入的文本可能丢失换行，或在段首带入空格。缺失的分段无法可靠地从选区文本还原；匹配到已保存高亮时使用 Kobo 保存的全文。
- Markdown 的章节名依赖 Kobo 目录和可读取的书籍文件。锚点无法匹配时，标题会退回文件名或“未知章节”。章节定位只读取相关书籍片段，不用于重建选区文本。
- 已在 Kobo Libra Colour、Kobo Clara 2E 上测试部分流程；v0.1.2 最终包尚未完成整套实机复测，其他设备和固件尚未验证。

## 版本

[v0.1.2](https://github.com/raenut/NickelGram/tree/v0.1.2) · [v0.1.1](https://github.com/raenut/NickelGram/tree/v0.1.1) · [v0.1.0](https://github.com/raenut/NickelGram/tree/v0.1.0)

开发与本地测试说明见[测试说明](docs/TESTING.md)。

第三方组件信息见[第三方许可证](docs/THIRD_PARTY_LICENSES.md)。

## 许可证

本项目采用 [MIT 许可证](LICENSE)。版权所有 © 2026 NickelGram contributors。

---

Made with 🤍 by [raenut](https://github.com/raenut) 👩🏻‍💻
