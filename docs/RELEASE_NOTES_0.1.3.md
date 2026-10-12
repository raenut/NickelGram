# NickelGram 0.1.3

## 新增与改进

- 保存高亮分享时，尽量保留 Kobo 记录的段落分隔；段首缩进会转换为换行。
- 整书 Markdown 导出优先使用 Kobo 目录标题；缺少标题时显示 “Untitled Chapter”，不再用内部文件名代替。

## 已知问题

- 直接分享选区时，Kobo/NickelMenu 可能不提供段落分隔。需要保留分段时，请先保存高亮；只有当前书中存在唯一匹配时，才会采用 Kobo 保存的全文和批注。
- Kobo Clara 2E 上拖动选区有时会卡顿，原因及其与书籍格式的关系尚未确认。
- 不同书籍格式和段落标记尚未覆盖验证；其他设备与固件也未验证。

---

# NickelGram 0.1.3

## Improvements

- Sharing a saved highlight now preserves Kobo's paragraph breaks when available, converting paragraph indentation to line breaks.
- Whole-book Markdown exports prefer Kobo table-of-contents titles. Missing titles use “Untitled Chapter” instead of internal filenames.

## Known issues

- Direct selections may not include paragraph boundaries. Save the highlight when paragraph breaks matter; NickelGram uses the saved text and annotation only when it uniquely matches within the current book.
- Selection dragging can lag on Kobo Clara 2E. The cause and any relationship to book format remain unconfirmed.
- Book-format coverage is incomplete. Other devices and firmware versions have not been verified.
