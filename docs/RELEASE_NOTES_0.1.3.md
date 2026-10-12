# NickelGram 0.1.3

## 改进

- 保存高亮后分享，可保留 Kobo 记录的段落换行。
- Markdown 章节标题优先使用 Kobo 目录标题；缺失时显示
  “Untitled Chapter”。

## 已知限制

- 直接选区可能不含段落分隔。保存高亮后，只有当前书中唯一匹配时，
  才会使用 Kobo 保存的全文和批注。
- Clara 2E 拖动选区偶有卡顿，原因尚未确认。
- 其他书籍格式、设备和固件未全面验证。

---

# NickelGram 0.1.3

## Improvements

- Saved highlights preserve Kobo's paragraph breaks when available.
- Markdown chapter headings use Kobo's table of contents, with
  “Untitled Chapter” as a fallback.

## Known limits

- Direct selections may omit paragraph breaks. A saved highlight is used only
  when it uniquely matches within the current book.
- Selection dragging can lag on Clara 2E; the cause is unknown.
- Other book formats, devices, and firmware are not fully verified.
