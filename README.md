# NickelGram

NickelGram 是供 Kobo 阅读器与 NickelMenu 使用的小工具：把选中文字、最新高亮/批注，或一本书的高亮导出为 Markdown 并发送到 Telegram。项目仍处于早期阶段。

## v0.1.0 的功能

- **选中分享**：在选区菜单将选中文本作为普通 Telegram 消息发送，附书名、作者和可选 footer。
- **最新高亮/批注**：从 Reader 菜单发送数据库中识别出的当前书最新一条可见高亮或批注。
- **整书导出**：把识别出的当前书的可见高亮与批注写入 `.md` 文件，并作为 Telegram 文件发送。文件包含 YAML front matter、`<hr>` 分隔线、可配置标签和独立的 Markdown footer。
- **本地预览**：用电脑的 `sqlite3` 和项目同一套查询与格式化代码预览普通消息、Markdown 和换行；本地测试默认不访问 Telegram。

## 安装

需要已安装 **NickelMenu** 的 Kobo 阅读器。本发布包针对 32 位 ARM Linux 构建，未验证其他架构。安装前请备份 `.kobo/KoboReader.sqlite`、现有 `.adds/nickelgram/config.json` 和 NickelMenu 配置。

1. 从 [v0.1.0 Release](https://github.com/raenut/NickelGram/releases/tag/v0.1.0) 下载 `NickelGram-0.1.0-kobo-arm.zip`。
2. 解压后把包内的 `.adds` **合并**到 Kobo 磁盘根目录；保留已有 `.adds` 的其他内容。包内 `.adds/nm/nickelgram` 提供三个菜单项。
3. 将 `.adds/nickelgram/config.example.json` 复制为同目录的 `config.json`，填入自己的 Telegram Bot Token 和目标 Chat ID。真实配置只保存在设备上，不在 Release 包内。
4. 安全弹出 Kobo 并重启；确认菜单项出现。Bot 还须具有向目标聊天或频道发送消息/文件的权限。

配置示例（占位值不可直接使用）：

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

`footer` 控制普通消息末尾；`md_footer` 独立控制 Markdown 文件末尾。把相应开关设为 `false` 即可关闭。`md_footer_text` 可包含 Markdown 格式；`md_tags` 可设为 `["书摘", "文学"]`。旧版 `"footer": "文字"` 仍可用于普通消息。

## 使用

| 入口 | 操作 |
| --- | --- |
| 选中文字后的菜单 | **Send to Telegram**：发送所选文字 |
| 阅读时右上角 Reader 菜单 | **Send Latest Highlight**：发送最新高亮/批注 |
| 阅读时右上角 Reader 菜单 | **Export Highlights to Telegram**：生成并发送整书 `.md` |

导出的 Markdown 写在 `.adds/nickelgram/exports/`。文件名由书名与基于 Kobo `ContentID` 的短标识组成，重复导出同一本书会使用相同文件名。电脑端的重复测试命令见 [本地测试说明](docs/LOCAL_TEST.md)。源码包为 `NickelGram-0.1.0-source.zip`；随 Release 提供 `SHA256SUMS`。

## Known Issues / Limitations

- **刚打开新书的书名/作者**：Kobo 的 `DateLastRead` 可能仍指向上一本书。若新高亮/批注已写入数据库，程序会优先使用该 Bookmark 的 `VolumeID`；仅选中文字且尚未保存高亮时，仍可能显示上一本书的书名和作者。Reader 菜单的当前书上下文也不能仅靠本地快照完全验证。
- **选区换行**：程序会保留 NickelMenu 传入的实际换行，但 Kobo 选区菜单是否把跨段换行完整传给脚本尚未实机确认。
- **旧高亮再次点击分享**：不作为受支持功能；若菜单只传来句中片段，消息可能被截断。普通选中分享、画线/标注后的 Reader 菜单分享和整书导出是当前的三个入口。
- **本次发行包的实机回归**：本地数据库查询、格式化及 ARM 构建已验证；v0.1.0 包在真实设备上的三个菜单流程、Telegram 实际发送和不同固件兼容性仍需验证。

## Risk Disclaimer

NickelGram v0.1.0 是早期版本，可能含有 Bug，可能遇到数据库状态变化、兼容性问题、发送失败、错误书籍归属，或安装/运行期间的数据丢失等风险。安装前请自行备份设备数据及配置；使用者应自行判断并承担使用风险。本项目按 `LICENSE` 所述按现状提供，不保证适用于任何特定用途。

作者提供的实机测试范围**仅**为 Kobo Libra Colour 与 Kobo Clara 2E 的当时最新固件，搭配 NickelMenu 0.6.0。其他设备、固件和 NickelMenu 版本未验证；本次 v0.1.0 发行包尚未在上述设备完成全流程回归。

第三方组件与许可证见 [说明](docs/THIRD_PARTY.md)。
