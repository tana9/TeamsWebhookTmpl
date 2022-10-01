package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"strings"
)

// MessageCard TeamsのWebhookに送信するJSON形式
type MessageCard struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// TeamsのWebhookにメッセージ送信
func postWebhook(webhookUrl string, body []byte) error {
	res, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	fmt.Println(res.Status)
	return nil
}

// 読み込んだメッセージファイルをWebhook用のJSONファイルに変換する
func makeJson(lines []string) ([]byte, error) {
	// 1行目のタイトルを取得
	var title string
	if 1 <= len(lines) {
		title = strings.TrimPrefix(lines[0], "# ")
	}

	// 2行目以降の本文を取得
	var text string
	if 2 <= len(lines) {
		text = strings.Join(lines[1:], "\n")
	}
	jsonStr, err := json.Marshal(MessageCard{Title: title, Text: text})
	if err != nil {
		return nil, err
	}
	return jsonStr, nil
}

// メッセージファイル読み込み
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

// メッセージのファイルパスを取得
func inputPath() string {
	// ファイルパスをドラッグされた場合
	if 2 <= len(os.Args) {
		return os.Args[1]
	}

	path := prompter.Prompt("送信するメッセージファイルを指定", "")
	return path
}

// 実行
func run() error {
	log.SetOutput(os.Stdout)
	var err error
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// メッセージファイルのパスを取得
	path := inputPath()
	if _, err := os.Stat(path); err != nil {
		return errors.Errorf("指定されたファイルが存在しません: %s", path)
	}

	// メッセージファイルの行読み込み
	lines, err := readLines(path)
	if err != nil {
		return err
	}

	// メッセージファイルをJSON形式に変換
	jsonBody, err := makeJson(lines)
	if err != nil {
		return err
	}

	// 投稿確認
	confirm := prompter.YN("こちらの内容で投稿しますか?", true)
	if !confirm {
		return nil
	}

	// TeamsのWebhookに送信
	if err := postWebhook(config.WebHookUrl, jsonBody); err != nil {
		return err
	}

	return nil
}

// main
func main() {
	if err := run(); err != nil {
		log.Printf("error: %+v", err)
		prompter.YN("エラーを確認してください", true)
		os.Exit(1)
	}
}
