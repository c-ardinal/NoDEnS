package core

import (
	"encoding/json"
	"io/ioutil"
)

/****************************************************************************/
/* 内部型定義                                                               */
/****************************************************************************/

// 設定情報格納構造体
type SystemConfig struct {
	BotToken string `json:"discord-token"`
	BotID    string `json:"discord-botid"`
	GuildId  string `json:"discord-guildid"`
	EndPoint string `json:"bcdice-endpoint"`
}

/****************************************************************************/
/* 内部定数定義                                                             */
/****************************************************************************/

/****************************************************************************/
/* 内部変数定義                                                             */
/****************************************************************************/

// 設定情報格納変数
var myConfig SystemConfig

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/

// 設定情報読み込み処理
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

// 設定情報取得処理
func GetConfig() SystemConfig {
	return myConfig
}
