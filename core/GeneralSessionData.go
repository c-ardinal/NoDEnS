package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Session セッションデータ構造体
type Session struct {
	Discord  Discord                `json:"discord"`
	Scenario Scenario               `json:"scenario"`
	Pc       map[string]interface{} `json:"pc"`
	Npc      map[string]interface{} `json:"npc"`
}

// Discord Discord情報構造体
type Discord struct {
	Parent NaID   `json:"text"`
	Child  []NaID `json:"private"`
	Voice  NaID   `json:"voice"`
}

// Scenario シナリオ情報構造体
type Scenario struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Keeper  NaID   `json:"keeper"`
	Chatlog string `json:"chatlog"`
}

// NaID 名前&ID紐づけ構造体
type NaID struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// trpgSession セッション情報格納マップ
var trpgSession = map[string]*Session{}

// NewSession セッション生成処理
func NewSession(chID string, system string, kpName string, kpID string) *Session {
	var newSession Session
	newSession.Pc = map[string]interface{}{}
	newSession.Npc = map[string]interface{}{}
	newSession.Discord.Parent.ID = chID
	newSession.Scenario.System = system
	newSession.Scenario.Keeper.Name = kpName
	newSession.Scenario.Keeper.ID = kpID
	trpgSession[chID] = &newSession
	return &newSession
}

// RemoveSession セッション削除処理
func RemoveSession(chID string) bool {
	result := false
	if CheckExistSession(chID) == true {
		delete(trpgSession, chID)
		result = true
	}
	return result
}

// StoreSession セッション保存処理
func StoreSession(chID string) (*os.File, error) {
	outputJSON, _ := json.MarshalIndent(*(trpgSession[chID]), "", "\t")
	os.Mkdir("./session_data/", 0755)
	file, err := os.Create("./session_data/" + chID + ".json")
	defer file.Close()
	file.Write(outputJSON)
	return file, err
}

// RestoreSession セッション復旧処理
func RestoreSession(chID string) error {
	rawData, err := ioutil.ReadFile("./session_data/" + chID + ".json")
	if err != nil {
		return err
	}

	if CheckExistSession(chID) == false {
		var newSession Session
		trpgSession[chID] = &newSession
	}

	var ses Session
	json.Unmarshal(rawData, &ses)
	trpgSession[chID] = &ses

	return nil
}

// CheckExistSession セッション存在有無チェック処理
func CheckExistSession(chID string) bool {
	_, result := trpgSession[chID]
	return result
}

// GetSessionByID セッション取得処理
func GetSessionByID(chID string) *Session {
	return trpgSession[chID]
}

// CheckExistCharacter プレイヤーキャラ登録有無チェック処理
func CheckExistCharacter(chID string, plID string) bool {
	ts, _ := trpgSession[chID]
	_, result := (*ts).Pc[plID]
	return result
}

// CheckExistNPCharacter ノンプレイヤーキャラ登録有無チェック処理
func CheckExistNPCharacter(chID string, plID string) bool {
	ts, _ := trpgSession[chID]
	_, result := (*ts).Npc[plID]
	return result
}

// GetCharacterName キャラ名取得処理
func GetCharacterName(chID string, plID string) string {
	var result string = ""
	ts, sesExist := trpgSession[chID]

	if sesExist == true {
		cdGetFunc := GetCharacterDataGetFunc((*ts).Scenario.System, "CharacterName")
		if cdGetFunc != nil {
			_, pcExist := (*ts).Pc[plID]
			_, npcExist := (*ts).Npc[plID]
			if pcExist == true {
				result = cdGetFunc((*ts).Pc[plID])
			} else if npcExist == true {
				result = cdGetFunc((*ts).Npc[plID])
			}
		}
	}
	return result
}

// GetParentIDFromChildID 親セッションID取得処理
func GetParentIDFromChildID(childID string) string {
	for _, ps := range trpgSession {
		for _, cid := range (*ps).Discord.Child {
			if childID == cid.ID {
				return (*ps).Discord.Parent.ID
			}
		}
	}
	return ""
}
