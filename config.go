package main

import (
	"github.com/Songmu/prompter"
	"github.com/spf13/viper"
)

// Config 設定情報
type Config struct {
	WebHookUrl string // TeamsのWebhookUrl
}

// LoadConfig 設定ファイル読み込み
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 設定ファイルが存在しない場合
			return CreateConfig()
		} else {
			return nil, err
		}
	}
	webhookUrl := viper.GetString("WebhookUrl")
	return &Config{WebHookUrl: webhookUrl}, nil
}

// CreateConfig 設定ファイルを作成
func CreateConfig() (*Config, error) {
	webhookUrl := prompter.Prompt("TeamsのWebhookUrlを入力", "")
	viper.Set("WebhookUrl", webhookUrl)
	if err := viper.WriteConfigAs("./config.yaml"); err != nil {
		return nil, err
	}
	return &Config{WebHookUrl: webhookUrl}, nil
}
