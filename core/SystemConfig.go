package core

import (
	"encoding/json"
	"io/ioutil"
)

// SystemConfig 設定情報格納構造体
type SystemConfig struct {
	BotToken string `json:"discord-token"`
	BotID    string `json:"discord-botid"`
	EndPoint string `json:"bcdice-endpoint"`
}

// myConfig 設定情報格納変数
var myConfig SystemConfig

// LoadConfig 設定情報読み込み処理
func LoadConfig(path string) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &myConfig); err != nil {
		return err
	}
	return nil
}

// GetConfig 設定情報取得処理
func GetConfig() SystemConfig {
	return myConfig
}
