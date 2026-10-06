# NickelGram

[English](README.md) · [下载最新版本](https://github.com/raenut/NickelGram/releases/latest) · [历次发布](https://github.com/raenut/NickelGram/releases)

NickelGram 通过 NickelMenu 为 Kobo 阅读器添加 Telegram 分享功能。你可以把选中文字或最新书摘作为消息发送，也可以把整本书的高亮与批注导出为 Markdown 文件。

## 功能与入口

| 位置 | 菜单项 | 作用 |
| --- | --- | --- |
| 选中文字后的菜单 | **Send to Telegram** | 连同书名、作者发送选中的文字 |
| 阅读时的 Reader 菜单 | **Send Latest Highlight** | 发送当前书籍最新一条可见高亮或批注 |
| 阅读时的 Reader 菜单 | **Export All Highlights** | 将当前书籍的可见高亮与批注保存为 `.md`，并将文件发送到 Telegram |

Markdown 书摘按书籍内容顺序、再按内容中的位置排列。同一章节中连续的书摘共用一个 Markdown 标题，在支持标题折叠的编辑器中可以折叠；每条书摘用 `<hr>` 分隔。若章节名称是 HTML 路径，只显示文件名，例如 `part0043.xhtml`。文件还包含 YAML front matter，其中 `exported_at` 精确到秒、不带时区；标签和 Markdown 页脚可选。缺少位置数据的书摘排在有位置数据的书摘之后。

导出文件保存在 Kobo 磁盘根目录的 `Highlights/` 文件夹。文件名由书名和基于 Kobo `ContentID` 生成的短标识组成；重复导出同一本书会使用相同的文件名。

## 安装

需要一台已安装 [NickelMenu](https://github.com/pgaskin/NickelMenu) 的 Kobo 阅读器。发布包针对 32 位 ARM Linux 构建，其他架构尚未验证。安装前请备份 `.kobo/KoboReader.sqlite`、现有的 `.adds/nickelgram/config.json` 和 NickelMenu 配置。

1. 从[最新发布版本](https://github.com/raenut/NickelGram/releases/latest)下载 Kobo ARM ZIP。发布页也提供源码包和 `SHA256SUMS`。
2. 解压 ZIP，将其中的 `.adds` 目录合并到 Kobo 磁盘根目录。保留原有 `.adds` 中的其他内容。包内的 `.adds/nm/nickelgram` 会添加上表中的三个菜单项。
3. 将 `.adds/nickelgram/config.example.json` 复制为同目录下的 `config.json`，填入 Telegram Bot Token 和目标 Chat ID。真实配置仅保存在设备上，不包含在发布包中。
4. 安全弹出并重启 Kobo，确认菜单项出现，并确认 Bot 有权限向目标聊天或频道发送消息和文件。

## 配置

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

将占位的 Token 和 Chat ID 换成自己的。`footer` 控制 Telegram 消息的页脚，`md_footer` 独立控制导出 Markdown 的页脚；将对应选项设为 `false` 即可关闭。`md_footer_text` 支持 Markdown 格式，`md_tags` 可设置为 `["书摘", "文学"]` 等列表。旧版的 `"footer": "文字"` 写法仍可用于普通消息。

## 已知限制

- **书籍识别：** 刚打开另一本书时，Kobo 的 `DateLastRead` 可能仍指向上一本书。已保存的高亮或批注可让 NickelGram 通过 Bookmark 的 `VolumeID` 识别当前书籍；仅分享选中文字时，书名和作者仍可能来自上一本书。
- **选区文字：** NickelGram 会保留 NickelMenu 传入的换行，但跨段落选区的换行尚未完成实机验证。重新打开旧高亮，再通过选区菜单分享不属于受支持的用法；传入的文字可能不完整。
- **实机验证：** 数据库查询、格式化逻辑和 ARM 构建已在本地检查；三个菜单流程及 Telegram 实际发送尚未完成完整实机测试，不同固件版本的兼容性也有待验证。

目前的测试设备仅包括 Kobo Libra Colour 和 Kobo Clara 2E，使用测试时可用的固件及 NickelMenu 0.6.0。其他设备和版本尚未验证。NickelGram 仍处于早期阶段，使用前请备份设备数据和配置。本项目按 [LICENSE](LICENSE) 条款提供，不作保证。

电脑端重复测试见[本地测试说明](docs/LOCAL_TEST.md)。第三方组件及其许可证见[第三方组件说明](docs/THIRD_PARTY.md)。
