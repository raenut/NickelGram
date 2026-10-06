# Third-Party Components & Licenses

- **Original NickelGram source code**

  Licensed under the MIT License. See [`LICENSE`](../LICENSE) in the project root.

- **Go standard library and runtime**

  Built with Go 1.26.8 and distributed under the BSD 3-Clause License. The release package includes `LICENSE-Go`. NickelGram currently has no additional Go module dependencies.

- **Mozilla CA certificate data**

  Licensed under MPL 2.0. The release package includes `LICENSE-MPL-2.0` and `cacert.pem`.

- **SQLite / `sqlite3` / `libsqlite3.so.0`**

  The bundled `sqlite3` and `libsqlite3.so.0` are taken from previously verified components compatible with the Kobo ARM environment. The SQLite core is in the public domain.

- **NickelMenu**

  NickelMenu must be installed separately by the user. The NickelGram release package only includes a NickelMenu configuration file for integrating NickelGram and does not redistribute NickelMenu itself.

## Release Package Restrictions

Release packages must not contain a real `config.json`, Telegram Bot Token, Chat ID, user highlight or excerpt text, or runtime state data.

---

# 第三方组件与许可证

- **NickelGram 原创源代码**

  使用 MIT License，详见项目根目录中的 [`LICENSE`](../LICENSE)。

- **Go 标准库与运行时**

  构建使用 Go 1.26.8，采用 BSD 3-Clause License。发布包中包含 `LICENSE-Go`。NickelGram 当前没有额外的 Go module 依赖。

- **Mozilla CA 证书数据**

  使用 MPL 2.0。发布包中包含 `LICENSE-MPL-2.0` 与 `cacert.pem`。

- **SQLite / `sqlite3` / `libsqlite3.so.0`**

  发布包附带的 `sqlite3` 与 `libsqlite3.so.0` 来自此前已验证可用于 Kobo ARM 环境的兼容组件。SQLite 核心代码属于 public domain。

- **NickelMenu**

  NickelMenu 需要由用户另行安装。NickelGram 发布包仅提供用于集成 NickelGram 的 NickelMenu 配置文件，不包含 NickelMenu 本体。

## 发布包内容限制

发布包不得包含真实的 `config.json`、Telegram Bot Token、Chat ID、用户书摘文本或运行状态数据。
