package cthulhu

import (
	"encoding/json"
	"strconv"
	"strings"

	"Nodens/core"
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

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/

// 能力値取得
func GetSkillNum(pc *CharacterOfCthulhu, skill string, stype string) string {
	var returnNum = -1
	ua, exist := pc.Ability[strings.ToUpper(skill)]
	if exist {
		switch stype {
		case "init":
			returnNum = (*ua).Init
		case "sum":
			returnNum = (*ua).Sum
		case "now":
			returnNum = (*ua).Now
		default:
		}
	} else {
		for _, s := range pc.Skill {
			if strings.Contains((*s).Name+" "+(*s).Sub, skill) {
				switch stype {
				case "init":
					returnNum = (*s).Init
				case "sum":
					returnNum = (*s).Sum
				case "now":
					returnNum = (*s).Now
				default:
				}
				break
			}
		}
	}
	return strconv.Itoa(returnNum)
}

// 能力値操作
func AddSkillNum(pc *CharacterOfCthulhu, skill string, add string) string {
	var returnNum = -1

	addNum, _ := strconv.Atoi(add)
	ua, exist := pc.Ability[strings.ToUpper(skill)]
	if exist {
		(*ua).Now += addNum
		returnNum = (*ua).Now
	} else {
		for _, sk := range pc.Skill {
			if strings.Contains((*sk).Name+" "+(*sk).Sub, skill) {
				(*sk).Now += addNum
				returnNum = (*sk).Now
				break
			}
		}
	}
	return strconv.Itoa(returnNum)
}

// キャラクター名取得
func GetCharacterName(pc interface{}) string {
	return pc.(*CharacterOfCthulhu).Personal.Name
}

// キャラクター名取得
func GetCharacterSheetUrl(pc interface{}) string {
	return pc.(*CharacterOfCthulhu).URL
}

// セッション復元処理(システム固有部)
func JobRestoreSession(ses *core.Session) bool {
	/* PC情報を一度JSONに戻してからクトゥルフ用PC構造体に変換する */
	pcsRawData, _ := json.Marshal((*ses).Pc)
	var pcsMap = map[string]*CharacterOfCthulhu{}
	json.Unmarshal(pcsRawData, &pcsMap)

	/* NPC情報を一度JSONに戻してからクトゥルフ用NPC構造体に変換する */
	npcsRawData, _ := json.Marshal((*ses).Npc)
	var npcsMap = map[string]*CharacterOfCthulhu{}
	json.Unmarshal(npcsRawData, &npcsMap)

	/* PC情報を格納 */
	for _, pcData := range pcsMap {
		(*ses).Pc[pcData.Player.ID] = pcData
	}
	/* NPC情報を格納 */
	for _, npcData := range npcsMap {
		(*ses).Npc[npcData.Player.ID] = npcData
	}
	return true
}
