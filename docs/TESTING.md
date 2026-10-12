# 本地测试

版本变更与限制见[发布说明](RELEASE_NOTES_0.1.4.md)。本地工具使用
`KoboReader.sqlite` 的只读副本和电脑上的 `sqlite3`，运行正式 Go 查询与
格式化代码。Debug 命令不会连接 Telegram。数据库、配置和输出均留在
Git 忽略的 `local-test/`。

## 准备

电脑需要 `sqlite3` 和 Go 1.23 或更新版本。从 Kobo 备份中取得
`KoboReader.sqlite`，在项目根目录运行：

```sh
KOBO_DB_SOURCE=/path/to/KoboReader.sqlite ./test.sh setup
./test.sh list-books
```

`setup` 会复制数据库，不修改原文件。源数据库更新后重新运行。工具不在
`PATH` 时，可设置 `SQLITE3=/path/to/sqlite3` 或 `GO=/path/to/go`。
Kobo 安装包使用 ARM 版 `sqlite3`；本地测试使用电脑版本，SQL 查询相同。

## 常用命令

```sh
./test.sh list-books
./test.sh list-highlights
./test.sh current-book
./test.sh latest-highlight '<书籍标识>'
./test.sh render-message '<BookmarkID>'
./test.sh export-md '<书籍标识>'
```

书籍标识可使用 `current`。`render-message` 写入
`local-test/output/message.txt`；`export-md` 输出按书中位置排序的 Markdown，
包含章节标题和 `exported_at`。Kobo 上的正式导出写入根目录 `Highlights/`。

检查跨段选区可运行：

```sh
./test.sh render-selection current '第一段

第二段'
```

输出中的 `RAW_INPUT_RAW`、`RAW_TEXT_RAW` 和 `FINAL_MESSAGE_RAW` 会将换行显示为
`\n`。编号的 `LINE` 行会显示文本行和空行。

使用自定义配置预览时，设置 `NICKELGRAM_TEST_CONFIG=/path/to/config.json`。
工具会读取配置并预览消息，但不会输出 Token 或发送 Telegram 请求。

## 本地测试范围

本地测试可检查数据库结构、书名、作者、高亮、批注、时间、消息格式、
Markdown 和文件名。数据库副本可能与设备现状不同；跨段效果也取决于
书籍和样本。

发布前检查书名/作者格式、缺失作者、书摘等宽样式、自定义标签，以及
有链接和无链接的页脚。还需在 Kobo 上检查普通选区、已保存高亮、
最近高亮和整书导出。

本地测试不能确认 NickelMenu 传入的 `{1}` 文本、阅读界面行为、Kobo ARM 运行或
Telegram 实际发送。没有段落标记时，程序无法推断原书的分段；数据库中
没有的批注也无法读取。
