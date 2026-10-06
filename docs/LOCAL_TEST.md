# 本地测试

本地入口使用你自己的 `KoboReader.sqlite` 只读快照和电脑上的 `sqlite3`，调用 NickelGram 正式 Go 程序中的查询与格式化代码。`debug` 命令不会建立 Telegram 客户端或调用 Telegram API。真实数据库、配置和生成结果都留在被 Git 忽略的 `local-test/`。

## 第一次运行

电脑需要 `sqlite3` 和 Go 1.23 或更新版本。从 Kobo 备份中取得 `KoboReader.sqlite`，在项目根目录执行：

```sh
KOBO_DB_SOURCE=/path/to/KoboReader.sqlite ./test.sh setup
./test.sh list-books
```

`setup` 会复制一个本地快照；原始数据库不被修改。数据库有变化时，再运行一次 `setup`。如果 `sqlite3` 或 Go 不在 `PATH`，可分别设置 `SQLITE3=/path/to/sqlite3`、`GO=/path/to/go`。Kobo 安装包继续使用其中的 ARM `sqlite3`；电脑测试使用电脑平台的 `sqlite3`，查询 SQL 保持一致。

```text
local-test/
├── KoboReader.sqlite       # 私人的数据库快照，不提交
├── footer.txt              # 未指定配置时的测试 footer
├── bin/nickelgram-debug    # 电脑端构建产物
├── cache/                  # Go 构建缓存
└── output/                 # 最终消息、Markdown 和排查结果
```

如需测试自己的配置，可在命令前加 `NICKELGRAM_TEST_CONFIG=/path/to/config.json`。程序会解析配置，但只使用 footer 与标签生成预览；不会发送 Telegram 请求，也不会输出 Token。未指定时使用 `local-test/footer.txt` 的内容。

## 常用命令

```sh
./test.sh list-books
./test.sh list-highlights
./test.sh latest-highlight '<list-books 输出的书籍标识>'
./test.sh render-message '<list-highlights 输出的 BookmarkID>'
./test.sh export-md '<list-books 输出的书籍标识>'
```

`latest-highlight` 和 `export-md` 的书籍标识也可填 `current`。`./test.sh current-book` 显示正式代码当前推断的书籍。`render-message` 写入 `output/message.txt`；`export-md` 写入带书名与短标识的 `.md`。v0.1.1 的 `export-md` 预览按书中位置排序，包含章节标题和不带时区的 `exported_at`。电脑端预览仍保存在 `local-test/output/`；Kobo 上的正式导出保存到磁盘根目录的 `Highlights/`。

排查跨段换行可运行：

```sh
./test.sh render-selection current '第一段

第二段'
```

输出中的 `RAW_INPUT_RAW`、`RAW_TEXT_RAW`、`FINAL_MESSAGE_RAW` 用 `\n` 标出换行；编号的 `LINE` 行显示每一行和空行。

## 本地环境的边界

本地可检查这份快照的 SQLite schema、书名、作者、Text、Annotation、时间、footer、普通消息、Markdown 和文件名。快照中的内容与设备实时状态可能不同；是否有多段高亮取决于你提供的数据库样本。

本地不能确认 NickelMenu `{1}` 实际传入的文本、selection UI 的行为、新书刚打开时 Bookmark 何时写入、Reader 菜单触发状态、Kobo ARM 环境和 Telegram 实际发送。这些需要最后的实机验证。尚未保存高亮的纯选区仍可能因 `DateLastRead` 滞后而得到上一书的元数据。
