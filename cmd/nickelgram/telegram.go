package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	BotToken string   `json:"telegram_bot_token"`
	ChatID   string   `json:"telegram_chat_id"`
	Footer   string   `json:"footer"`
	MDFooter string   `json:"md_footer"`
	MDTags   []string `json:"md_tags"`
}

var tokenPattern = regexp.MustCompile(`^[0-9]+:[A-Za-z0-9_-]+$`)
var chatPattern = regexp.MustCompile(`^(-?[0-9]+|@[A-Za-z0-9_]+)$`)

func readConfig(root string) (Config, error) {
	return readConfigFile(filepath.Join(root, "config.json"))
}

func readConfigFile(path string) (Config, error) {
	var raw struct {
		Token        string          `json:"telegram_bot_token"`
		Chat         string          `json:"telegram_chat_id"`
		LegacyToken  string          `json:"bot_token"`
		LegacyChat   string          `json:"chat_id"`
		Footer       json.RawMessage `json:"footer"`
		FooterText   string          `json:"footer_text"`
		MDFooter     *bool           `json:"md_footer"`
		MDFooterText string          `json:"md_footer_text"`
		MDTags       []string        `json:"md_tags"`
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, errors.New("未找到 config.json；请填写配置后重试。")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if d.Decode(&raw) != nil || d.Decode(new(any)) != io.EOF {
		return Config{}, errors.New("配置格式有误；请检查 Telegram 配置。")
	}
	if raw.Token != "" || raw.Chat != "" {
		if raw.LegacyToken != "" || raw.LegacyChat != "" {
			return Config{}, errors.New("配置中不能混用新旧字段。")
		}
	} else {
		raw.Token, raw.Chat = raw.LegacyToken, raw.LegacyChat
	}
	c := Config{BotToken: raw.Token, ChatID: raw.Chat}
	if !tokenPattern.MatchString(c.BotToken) || !chatPattern.MatchString(c.ChatID) {
		return Config{}, errors.New("配置格式有误；请检查 Bot Token 和 Chat ID。")
	}
	if len(raw.Footer) != 0 {
		var enabled bool
		if json.Unmarshal(raw.Footer, &enabled) == nil {
			if enabled {
				c.Footer = raw.FooterText
			}
		} else if json.Unmarshal(raw.Footer, &c.Footer) != nil {
			return Config{}, errors.New("footer 必须是 true、false 或旧版字符串。")
		}
	}
	if raw.MDFooter != nil && *raw.MDFooter {
		c.MDFooter = raw.MDFooterText
	}
	c.MDTags = raw.MDTags
	return c, nil
}

type Outcome struct {
	OK      bool
	Message string
}
type Telegram struct {
	Client  *http.Client
	BaseURL string
}

func newTelegram(root string) (*Telegram, error) {
	pem, err := os.ReadFile(filepath.Join(root, "cacert.pem"))
	pool := x509.NewCertPool()
	if err != nil || !pool.AppendCertsFromPEM(pem) {
		return nil, errors.New("CA 证书缺失或损坏；请重新安装 NickelGram。")
	}
	return &Telegram{Client: &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }, Transport: &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: 3 * time.Second}).DialContext, TLSClientConfig: &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}, TLSHandshakeTimeout: 3 * time.Second, ResponseHeaderTimeout: 4 * time.Second, DisableKeepAlives: true}}, BaseURL: "https://api.telegram.org"}, nil
}
func (t *Telegram) SendMessage(c Config, document string) Outcome {
	payload := map[string]any{
		"chat_id":              c.ChatID,
		"text":                 document,
		"link_preview_options": map[string]bool{"is_disabled": true},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return Outcome{Message: "无法编码 Telegram 请求。"}
	}
	req, err := http.NewRequest("POST", t.BaseURL+"/bot"+c.BotToken+"/sendMessage", bytes.NewReader(body))
	if err != nil {
		return Outcome{Message: "无法建立 Telegram 请求。"}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.Client.Do(req)
	if err != nil {
		return Outcome{Message: "Telegram 网络请求失败；结果可能未知。"}
	}
	defer resp.Body.Close()
	var api struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
		ErrorCode   int    `json:"error_code"`
		Result      struct {
			MessageID int64 `json:"message_id"`
		} `json:"result"`
	}
	data, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if readErr == nil && json.Unmarshal(data, &api) == nil && resp.StatusCode == 200 && api.OK && api.Result.MessageID > 0 {
		return Outcome{OK: true}
	}
	if readErr == nil {
		description := strings.ReplaceAll(api.Description, c.BotToken, "[已隐藏]")
		if description != "" {
			return Outcome{Message: fmt.Sprintf("Telegram API 错误 %d：%s", api.ErrorCode, description)}
		}
	}
	return Outcome{Message: fmt.Sprintf("Telegram 未确认发送成功（HTTP %d）。", resp.StatusCode)}
}

func (t *Telegram) SendDocument(c Config, path, name string) Outcome {
	file, err := os.Open(path)
	if err != nil {
		return Outcome{Message: "无法读取 Markdown 文件。"}
	}
	defer file.Close()
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	if err = form.WriteField("chat_id", c.ChatID); err != nil {
		return Outcome{Message: "无法编码 Telegram 文件请求。"}
	}
	part, err := form.CreateFormFile("document", name)
	if err != nil {
		return Outcome{Message: "无法编码 Telegram 文件请求。"}
	}
	if _, err = io.Copy(part, file); err != nil {
		return Outcome{Message: "无法读取 Markdown 文件。"}
	}
	if err = form.Close(); err != nil {
		return Outcome{Message: "无法编码 Telegram 文件请求。"}
	}
	req, err := http.NewRequest("POST", t.BaseURL+"/bot"+c.BotToken+"/sendDocument", &body)
	if err != nil {
		return Outcome{Message: "无法建立 Telegram 文件请求。"}
	}
	req.Header.Set("Content-Type", form.FormDataContentType())
	client := *t.Client
	client.Timeout = 30 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return Outcome{Message: "Telegram 文件请求失败；结果可能未知。"}
	}
	defer resp.Body.Close()
	var api struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
		ErrorCode   int    `json:"error_code"`
		Result      struct {
			MessageID int64 `json:"message_id"`
		} `json:"result"`
	}
	data, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if readErr == nil && json.Unmarshal(data, &api) == nil && resp.StatusCode == 200 && api.OK && api.Result.MessageID > 0 {
		return Outcome{OK: true}
	}
	if readErr == nil && api.Description != "" {
		return Outcome{Message: fmt.Sprintf("Telegram API 错误 %d：%s", api.ErrorCode, strings.ReplaceAll(api.Description, c.BotToken, "[已隐藏]"))}
	}
	return Outcome{Message: fmt.Sprintf("Telegram 未确认文件发送成功（HTTP %d）。", resp.StatusCode)}
}

func runAction(root, selection string) error {
	lock, err := acquireLock()
	if err != nil {
		return errors.New("另一个 NickelGram 操作正在运行。")
	}
	defer func() { _ = lock.Close(); _ = os.Remove(lockPath) }()
	selection, err = validateSelection(selection)
	if err != nil {
		return err
	}
	c, err := readConfig(root)
	if err != nil {
		return err
	}
	tg, err := newTelegram(root)
	if err != nil {
		return err
	}
	if err = verifySQLiteRuntime(root); err != nil {
		return err
	}
	binary, lib := filepath.Join(root, "sqlite3"), filepath.Join(root, "lib")
	book, err := readCurrentBook(binary, lib, koboDB)
	if err != nil {
		return err
	}
	if full := readMatchingBookmarkText(binary, lib, koboDB, book, selection); full != "" {
		selection = full
	}
	outcome := tg.SendMessage(c, selectionMessage(book, selection, c.Footer))
	if !outcome.OK {
		return errors.New(outcome.Message)
	}
	fmt.Println("发送成功")
	return nil
}

func runLatestHighlight(root string) error {
	lock, err := acquireLock()
	if err != nil {
		return errors.New("另一个 NickelGram 操作正在运行。")
	}
	defer func() { _ = lock.Close(); _ = os.Remove(lockPath) }()
	c, err := readConfig(root)
	if err != nil {
		return err
	}
	tg, err := newTelegram(root)
	if err != nil {
		return err
	}
	if err = verifySQLiteRuntime(root); err != nil {
		return err
	}
	binary, lib := filepath.Join(root, "sqlite3"), filepath.Join(root, "lib")
	book, err := readCurrentBook(binary, lib, koboDB)
	if err != nil {
		return err
	}
	mark, err := readLatestHighlight(binary, lib, koboDB, book)
	if err != nil {
		return err
	}
	outcome := tg.SendMessage(c, highlightMessage(book, mark, c.Footer))
	if !outcome.OK {
		return errors.New(outcome.Message)
	}
	fmt.Println("发送成功")
	return nil
}
