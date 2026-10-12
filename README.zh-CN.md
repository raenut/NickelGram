<p align="center">
  <img src="assets/nickelgram-icon.svg" alt="NickelGram 图标" width="144" height="144">
</p>

# NickelGram

[English](README.md) · [下载 v0.1.4](https://github.com/raenut/NickelGram/releases/tag/v0.1.4)

NickelGram 通过 [NickelMenu](https://github.com/pgaskin/NickelMenu) 为 Kobo
阅读器添加 Telegram 分享功能。

发送选中文字或最近一条书摘，也可将全书高亮和批注导出为 Markdown。

## 功能

| 入口 | 菜单项 | 作用 |
| --- | --- | --- |
| 文字选区 | Send&nbsp;to&nbsp;Telegram | 发送选中文字及书籍信息 |
| 阅读菜单 | Send&nbsp;Latest&nbsp;Highlight | 发送最近一条高亮或批注 |
| 阅读菜单 | Export&nbsp;All&nbsp;Highlights | 导出并发送全书高亮（Markdown） |

Markdown 文件保存在 Kobo 存储根目录的 `Highlights/` 文件夹。

## 安装

- 安装 NickelMenu。
- 创建 Telegram Bot，并准备目标聊天的 Chat ID。
- 备份 `.kobo/KoboReader.sqlite` 和已有的 NickelGram/NickelMenu 配置。

1. 从 [v0.1.4 发布页](https://github.com/raenut/NickelGram/releases)
   下载 `NickelGram-0.1.4-kobo-arm.zip`。
2. 解压后将 `.adds/` 合并到 Kobo 存储根目录。
3. 将 `.adds/nickelgram/config.example.json` 复制为同目录下的 `config.json`。
4. 填入 Bot Token 和 Chat ID，安全弹出并重启 Kobo。

当前发布包适用于 32 位 ARM Linux Kobo 阅读器。

## 配置

编辑 `.adds/nickelgram/config.json`：

| 参数 | 作用 |
| --- | --- |
| `telegram_bot_token` | Telegram Bot Token |
| `telegram_chat_id` | 目标聊天 |
| `footer` | 是否添加 Telegram 页脚 |
| `footer_text` | 页脚文字 |
| `footer_url` | 可选页脚链接；不显示网页预览 |
| `footer_link_text` | 可点击文字；设置链接时必填 |
| `header_format` | 设置书名、作者顺序和标点 |
| `quote_label` | 书摘标签；正文使用等宽字体 |
| `annotation_label` | 批注标签 |
| `md_footer` | 是否添加 Markdown 页脚 |
| `md_footer_text` | Markdown 页脚文字 |
| `md_tags` | 导出文件标签 |

`header_format` 仅支持 `{title}` 和 `{author}`，且必须包含 `{title}`。
两者分别使用 Kobo 提供的书名和作者。作者可能为空，此时模板中的标点会
保留。标题模板和标签不能包含换行。旧配置仍可使用；缺少新字段时
使用默认值。

`config.json` 包含 Token 和 Chat ID，请勿公开分享。默认值见
`.adds/nickelgram/config.example.json`。

## 限制

- 直接选区可能缺少段落换行。先保存高亮；当前书中唯一匹配时，
  NickelGram 会使用 Kobo 保存的全文。
- Markdown 章节标题取自 Kobo 目录；缺少标题时显示 “Untitled Chapter”。
- v0.1.4 已安装到 Kobo Clara 2E 并通过启动检查；Telegram 发送和完整
  阅读流程尚未全部验证。其他设备和固件也未验证。

## 发布说明

- [v0.1.4](docs/RELEASE_NOTES_0.1.4.md)
- [v0.1.3](docs/RELEASE_NOTES_0.1.3.md)
- [v0.1.2](docs/RELEASE_NOTES_0.1.2.md)
- [v0.1.1](docs/RELEASE_NOTES_0.1.1.md)
- [v0.1.0](docs/RELEASE_NOTES_0.1.0.md)

另见[本地测试说明](docs/TESTING.md)和[第三方许可证](docs/THIRD_PARTY_LICENSES.md)。

## 许可证

MIT 许可证。版权所有 © 2026 NickelGram contributors。
