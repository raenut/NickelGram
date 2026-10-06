# NickelGram

[English](README.md) | 中文

NickelGram 是一个为 Kobo 阅读器与 NickelMenu 设计的轻量工具，可以将选中的文字、最新高亮或批注发送到 Telegram，也可以将一本书中的高亮与批注导出为 Markdown 文件并发送到 Telegram。

项目目前仍处于早期阶段，功能、兼容性和使用方式可能在后续版本中继续调整。

## 功能

### v0.1.0

- **选中分享**：在文字选区菜单中，将当前选中的内容作为普通 Telegram 消息发送，并附带书名、作者以及可选的 `footer`。

- **最新高亮或批注**：从 Reader 菜单发送数据库中识别到的当前书籍最新一条可见高亮或批注。

- **整书导出**：将当前书籍中识别到的可见高亮与批注导出为 `.md` 文件，并作为 Telegram 文件发送。导出的 Markdown 包含 YAML front matter、<hr> 分隔线、可配置标签以及独立的 Markdown `footer`。

## 安装

需要一台已安装 NickelMenu 的 Kobo 阅读器。

当前发布包针对 32 位 ARM Linux 构建，尚未验证其他架构。安装前建议备份 `.kobo/KoboReader.sqlite`、现有的 `.adds/nickelgram/config.json` 以及 NickelMenu 配置。

1. 从 [v0.1.0 Release](https://github.com/raenut/NickelGram/releases/tag/v0.1.0) 下载 `NickelGram-0.1.0-kobo-arm.zip`。

2. 解压压缩包，将其中的 `.adds` 目录合并到 Kobo 磁盘根目录。不要删除原有 `.adds` 中的其他内容。包内的 `.adds/nm/nickelgram` 会添加三个 NickelMenu 菜单项。

3. 将 `.adds/nickelgram/config.example.json` 复制为同目录下的 `config.json`，然后填入自己的 Telegram Bot Token 和目标 Chat ID。真实配置仅保存在设备上，不包含在 Release 包中。

4. 安全弹出 Kobo 并重新启动，确认 NickelGram 菜单项已经出现。Telegram Bot 还需要具有向目标聊天或频道发送消息和文件的权限。

## 配置

配置示例：

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

以上内容仅为示例，占位值不能直接使用。

`footer` 控制普通 Telegram 消息末尾的 `footer`；`md_footer` 独立控制导出的 Markdown 文件末尾。将相应选项设为 `false` 即可关闭。

`md_footer_text` 可以包含 Markdown 格式，`md_tags` 可以设置为例如 `["书摘", "文学"]`。

旧版的 `"footer": "文字"` 写法仍可用于普通消息。

## 使用

| 入口 | 操作 |
| --- | --- |
| 选中文字后的菜单 | **Send to Telegram** — 发送当前选中的文字 |
| 阅读时右上角 Reader 菜单 | **Send Latest Highlight** — 发送最新高亮或批注 |
| 阅读时右上角 Reader 菜单 | **Export Highlights to Telegram** — 导出并发送整书 `.md` |

导出的 Markdown 文件保存在：

```text
.adds/nickelgram/exports/
```

文件名由书名以及基于 Kobo `ContentID` 生成的短标识组成。重复导出同一本书时会使用相同的文件名。

电脑端测试和重复测试命令见[本地测试说明](docs/LOCAL_TEST.md)。

Release 同时提供源码包 `NickelGram-0.1.0-source.zip` 与 `SHA256SUMS`。

## 已知问题与限制

- **刚打开新书时的书名与作者**：Kobo 的 `DateLastRead` 有时仍会指向上一本书。如果新高亮或批注已经写入数据库，NickelGram 会优先使用对应 Bookmark 的 `VolumeID` 来识别书籍。

  但如果只是选中文字、尚未保存为高亮或批注，则消息仍有可能显示上一本书的书名和作者。Reader 菜单中的“当前书籍”上下文也无法仅通过本地数据库快照完全验证。

- **选区换行**：NickelGram 会保留 NickelMenu 实际传入的换行，但 Kobo 的选区菜单是否会完整传递跨段落选择中的换行，目前尚未完成实机确认。

- **再次分享旧高亮**：目前不将“重新点击旧高亮并直接分享”作为受支持的使用方式。如果菜单只传入句子中的一部分文字，最终发送的内容可能被截断。

  当前正式支持的三个入口是：普通选中分享、保存高亮或批注后的 Reader 菜单分享，以及整书导出。

- **v0.1.0 实机回归测试**：本地数据库查询、格式化逻辑以及 ARM 构建已经验证，但 [v0.1.0 Release](https://github.com/raenut/NickelGram/releases/tag/v0.1.0) 包在真实设备上的三个完整菜单流程、Telegram 实际发送以及不同固件版本之间的兼容性仍需要进一步验证。

## 风险声明

NickelGram v0.1.0 是一个早期版本，可能仍包含 Bug，并可能受到 Kobo 数据库状态变化、固件或 NickelMenu 兼容性、网络或 Telegram 发送失败、书籍识别错误，以及安装或运行过程中意外数据丢失等问题影响。

安装前请自行备份设备数据与相关配置，并根据自己的情况判断是否使用。本项目按照 [LICENSE](LICENSE) 中的条款以“现状”提供，不保证适用于任何特定用途。

## 已测试环境

作者目前的实机测试环境仅包括：

- Kobo Libra Colour
- Kobo Clara 2E
- 测试时对应设备的最新固件
- NickelMenu 0.6.0

其他 Kobo 设备、固件版本以及 NickelMenu 版本尚未验证。v0.1.0 Release 包也尚未在上述设备上完成完整的端到端回归测试。

## 第三方组件与许可证

第三方组件及其许可证信息见 [THIRD_PARTY.md](docs/THIRD_PARTY.md)。
