# NickelGram 0.1.4

## 新增与改进

- Telegram 页脚支持链接；默认只链接 `NickelGram`，不显示网页预览。
- 书摘正文使用等宽字体；书摘和批注标签可自定义。
- 书名与作者的顺序和标点可通过 `header_format` 自定义。
- 选区开头与当前书中唯一高亮匹配至少 5 个字符时，发送完整高亮。
- 旧配置继续可用；新增字段缺省时采用默认值。

## 配置限制

- `header_format` 只支持 `{title}` 和 `{author}`，必须包含 `{title}`。
- 书名取自 Kobo 的 `content.Title`，作者取自 `content.Attribution`。
- 作者可能为空；标签和标题模板均不能包含换行。

## 验证与限制

- 已安装到 Kobo Clara 2E，并通过程序启动检查。
- Telegram 实际发送和完整阅读流程尚未全部验证。
- 直接选区缺少段落标记时，无法还原原书分段。
- Clara 2E 拖动选区偶有卡顿，原因尚未确认。

---

# NickelGram 0.1.4

## Changes

- Telegram footers support links. Only `NickelGram` is linked by default;
  previews stay off.
- Excerpts use monospace. Excerpt and annotation labels are configurable.
- Customize title/author order and punctuation with `header_format`.
- A unique match of at least five opening characters in the current book
  sends the full saved highlight.
- Older configurations remain compatible; new fields use defaults when omitted.

## Configuration limits

- `header_format` supports only `{title}` and `{author}` and requires `{title}`.
- Kobo provides the title from `content.Title` and author from
  `content.Attribution`.
- The author may be empty. Labels and header templates cannot contain newlines.

## Verification and limits

- Installed and launch-checked on Kobo Clara 2E.
- Telegram delivery and all reading flows have not been fully verified.
- Direct selections without paragraph markers cannot restore book paragraphs.
- Selection dragging can lag on Clara 2E; the cause is unknown.
