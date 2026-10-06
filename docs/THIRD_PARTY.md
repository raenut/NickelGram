# 第三方组件与许可证

- NickelGram 原创源代码：MIT，见项目 `LICENSE`。
- Go 标准库及运行时：构建使用 Go 1.26.8，BSD 3-Clause，发布包含 `LICENSE-Go`；无额外 Go module 依赖。
- Mozilla CA 证书数据：MPL 2.0，发布包含 `LICENSE-MPL-2.0` 与 `cacert.pem`。
- 随包 `sqlite3` / `libsqlite3.so.0` 来自早期已验证的 Kobo ARM 兼容组件；SQLite 核心为 public domain。
- NickelMenu 由用户另行安装，本包只提供一份 NickelMenu 配置文件。

发布包不得包含真实 `config.json`、Bot Token、Chat ID、书摘文本或运行状态。
