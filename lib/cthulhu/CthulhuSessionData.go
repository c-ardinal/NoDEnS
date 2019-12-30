package cthulhu

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

// CthulhuSession セッションデータ構造体
type CthulhuSession struct {
	Discord  Discord               `json:"discord"`
	Scenario Scenario              `json:"scenario"`
	Pc       map[string]*Character `json:"pc"`
	Npc      map[string]*Character `json:"npc"`
}

// Discord Discrod情報構造体
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

// Character キャラクタ情報構造体
type Character struct {
	Player     NaID                `json:"player"`
	Md         string              `json:"md"`
	ID         string              `json:"id"`
	Personal   Personal            `json:"personal-data"`
	Ability    map[string]*Ability `json:"ability"`
	Skill      map[string]*Skill   `json:"skill"`
	Battle     Battle              `json:"battle"`
	Belongings Belongings          `json:"belongings"`
	Memo       string              `json:"memo"`
}

// Personal キャラクタ個人情報構造体
type Personal struct {
	Name    string `json:"name"`
	Job     string `json:"job"`
	Age     int    `json:"age"`
	Sex     string `json:"sex"`
	Height  string `json:"height"`
	Weight  string `json:"weight"`
	Country string `json:"country"`
	Haircl  string `json:"hair-color"`
	Eyecl   string `json:"eye-color"`
	Skincl  string `json:"skin-color"`
}

// Ability キャラクタ能力情報構造体
type Ability struct {
	Name string `json:"name"`
	Init int    `json:"init"`
	Add  int    `json:"add"`
	Temp int    `json:"temp"`
	Sum  int    `json:"sum"`
	Start int   `json:"start"`
	Now  int    `json:"now"`
}

// Skill キャラクタ技能情報構造体
type Skill struct {
	Name    string `json:"name"`
	Sub     string `json:"sub-name"`
	Init    int    `json:"init"`
	Job     int    `json:"job"`
	Hobby   int    `json:"hobby"`
	Growth  int    `json:"growth"`
	Other   int    `json:"other"`
	Growflg int    `json:"growflg"`
	Sum     int    `json:"sum"`
	Start   int    `json:"start"`
	Now     int    `json:"now"`
}

// Battle キャラクタ戦闘情報構造体
type Battle struct {
	Db     string              `json:"db"`
	Weapon map[string]*UnitArm `json:"weapon"`
	Armor  map[string]*UnitArm `json:"armor"`
}

// UnitArm 戦闘情報構造体
type UnitArm struct {
	Name       string `json:"name"`
	Accuracy   int    `json:"accuracy"`
	Damage     string `json:"damage"`
	Scope      int    `json:"scope"`
	Continuous int    `json:"continuous"`
	Bullet     int    `json:"stock"`
	Durability int    `json:"durability"`
	Other      string `json:"other"`
}

// Belongings キャラクタ所持品情報構造体
type Belongings struct {
	Pocket int              `json:"pocket-money"`
	Other  int              `json:"other-money"`
	Items  map[string]*Item `json:"items"`
}

// Item アイテム情報構造体
type Item struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
	Other string `json:"other"`
}

// NaID 名前&ID紐づけ構造体
type NaID struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// cthulhuSession セッション情報格納マップ
var cthulhuSession = map[string]*CthulhuSession{}

// NewSession セッション生成処理
func NewSession(chID string, system string, kpName string, kpID string) *CthulhuSession {
	var newSession CthulhuSession
	newSession.Pc = map[string]*Character{}
	newSession.Npc = map[string]*Character{}
	newSession.Discord.Parent.ID = chID
	newSession.Scenario.System = system
	newSession.Scenario.Keeper.Name = kpName
	newSession.Scenario.Keeper.ID = kpID
	cthulhuSession[chID] = &newSession
	return &newSession
}

// RemoveSession セッション削除処理
func RemoveSession(chID string) bool {
	result := false
	if CheckDuplicateSession(chID) == true {
		delete(cthulhuSession, chID)
		result = true
	}
	return result
}

// StoreSession セッション保存処理
func StoreSession(chID string) (*os.File, error) {
	const format = "20060102150405"
	outputJSON, _ := json.MarshalIndent(*(cthulhuSession[chID]), "", "\t")
	file, err := os.Create("./session_data/" + time.Now().Format(format) + "_" + chID + ".json")
	defer file.Close()
	file.Write(outputJSON)
	return file, err
}

// LoadSession セッション復旧処理
func LoadSession(file *os.File) {
	// TODO: そのうち実装
}

// CheckDuplicateSession セッション重複チェック処理
func CheckDuplicateSession(chID string) bool {
	_, result := cthulhuSession[chID]
	return result
}

// GetSessionByID セッション取得処理
func GetSessionByID(chID string) *CthulhuSession {
	return cthulhuSession[chID]
}

// CheckDuplicateCharacter プレイヤーキャラ登録重複チェック処理
func CheckDuplicateCharacter(chID string, plID string) bool {
	ts, _ := cthulhuSession[chID]
	_, result := (*ts).Pc[plID]
	return result
}

// CheckDuplicateNPCharacter ノンプレイヤーキャラ登録重複チェック処理
func CheckDuplicateNPCharacter(chID string, plID string) bool {
	ts, _ := cthulhuSession[chID]
	_, result := (*ts).Npc[plID]
	return result
}

// GetParentIDFromChildID 親セッションID取得処理
func GetParentIDFromChildID(childID string) string {
	for _, ps := range cthulhuSession {
		for _, cid := range (*ps).Discord.Child {
			if childID == cid.ID {
				return (*ps).Discord.Parent.ID
			}
		}
	}
	return ""
}

// GetSkillNum 能力値取得
func GetSkillNum(pc *Character, skill string, stype string) string {
	var returnNum = -1
	ua, exist := pc.Ability[strings.ToUpper(skill)]
	if exist == true {
		switch(stype) {
		case "init":
			returnNum = (*ua).Init
			break
		case "sum":
			returnNum = (*ua).Sum
			break;
		case "now":
			returnNum = (*ua).Now
			break
		default:
			break
		}
	} else {
		for _, s := range pc.Skill {
			if strings.Contains((*s).Name+" "+(*s).Sub, skill) == true {
				switch(stype) {
				case "init":
					returnNum = (*s).Init
					break
				case "sum":
					returnNum = (*s).Sum
					break
				case "now":
					returnNum = (*s).Now
					break
				default:
					break
				}
				break
			}
		}
	}
	return strconv.Itoa(returnNum)
}

// AddSkillNum 能力値操作
func AddSkillNum(pc *Character, skill string, add string) string {
	var returnNum = -1

	addNum, _ := strconv.Atoi(add)
	ua, exist := pc.Ability[strings.ToUpper(skill)]
	if exist == true {
		(*ua).Now += addNum
		returnNum = (*ua).Now
	} else {
		for _, sk := range pc.Skill {
			if strings.Contains((*sk).Name+" "+(*sk).Sub, skill) == true {
				(*sk).Now += addNum
				returnNum = (*sk).Now
				break
			}
		}
	}
	return strconv.Itoa(returnNum)
}
