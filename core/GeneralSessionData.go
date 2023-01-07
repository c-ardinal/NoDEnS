package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

/****************************************************************************/
/* 内部型定義                                                               */
/****************************************************************************/

/****************************************************************************/
/* 内部定数定義                                                             */
/****************************************************************************/

/****************************************************************************/
/* 内部変数定義                                                             */
/****************************************************************************/

// セッション情報格納マップ
var trpgSession = map[string]*Session{}

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/

// セッション生成処理
func NewSession(chID string, system string, kpName string, kpID string) *Session {
	var newSession Session
	newSession.Pc = map[string]interface{}{}
	newSession.Npc = map[string]interface{}{}
	newSession.Chat.Parent.ID = chID
	newSession.Scenario.System = system
	newSession.Scenario.Keeper.Name = kpName
	newSession.Scenario.Keeper.ID = kpID
	trpgSession[chID] = &newSession
	return &newSession
}

// セッション削除処理
func RemoveSession(chID string) bool {
	result := false
	if CheckExistSession(chID) == true {
		delete(trpgSession, chID)
		result = true
	}
	return result
}

// セッション保存処理
func StoreSession(chID string) (*os.File, error) {
	outputJSON, _ := json.MarshalIndent(*(trpgSession[chID]), "", "\t")
	os.Mkdir("./session_data/", 0755)
	file, err := os.Create("./session_data/" + chID + ".json")
	defer file.Close()
	file.Write(outputJSON)
	return file, err
}

// セッション復旧処理
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

// セッション存在有無チェック処理
func CheckExistSession(chID string) bool {
	_, result := trpgSession[chID]
	return result
}

// セッション取得処理
func GetSessionByID(chID string) *Session {
	return trpgSession[chID]
}

// プレイヤーキャラクター登録有無チェック処理
func CheckExistCharacter(chID string, plID string) bool {
	ts, _ := trpgSession[chID]
	_, result := (*ts).Pc[plID]
	return result
}

// ノンプレイヤーキャラクター登録有無チェック処理
func CheckExistNPCharacter(chID string, plID string) bool {
	ts, _ := trpgSession[chID]
	_, result := (*ts).Npc[plID]
	return result
}

// キャラクター名取得処理
func GetCharacterData(chID string, plID string, dataName string) string {
	var result string = ""
	ts, sesExist := trpgSession[chID]

	if sesExist == true {
		cdGetFunc := GetCharacterDataGetFunc((*ts).Scenario.System, dataName)
		if cdGetFunc != nil {
			_, pcExist := (*ts).Pc[plID]
			if pcExist == true {
				result = cdGetFunc((*ts).Pc[plID])
			}
		}
	} else {
		parentId := GetParentIDFromChildID(chID)
		if parentId != "" {
			tsParent, _ := trpgSession[parentId]
			cdGetFunc := GetCharacterDataGetFunc((*tsParent).Scenario.System, dataName)
			if cdGetFunc != nil {
				_, npcExist := (*tsParent).Npc[plID]
				if npcExist == true {
					result = cdGetFunc((*tsParent).Npc[plID])
				}
			}
		}
	}
	return result
}

// 親セッションID取得処理
func GetParentIDFromChildID(childID string) string {
	for _, ps := range trpgSession {
		for _, cid := range (*ps).Chat.Child {
			if childID == cid.ID {
				return (*ps).Chat.Parent.ID
			}
		}
	}
	return ""
}
