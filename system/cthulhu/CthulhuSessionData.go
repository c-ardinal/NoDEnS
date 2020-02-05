package cthulhu

import (
	"strconv"
	"strings"

	"github.com/c-ardinal/Nodens/core"
)

// CharacterOfCthulhu キャラクタ情報構造体
type CharacterOfCthulhu struct {
	Player     core.NaID           `json:"player"`
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
	Name  string `json:"name"`
	Init  int    `json:"init"`
	Add   int    `json:"add"`
	Temp  int    `json:"temp"`
	Sum   int    `json:"sum"`
	Start int    `json:"start"`
	Now   int    `json:"now"`
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

// GetSkillNum 能力値取得
func GetSkillNum(pc *CharacterOfCthulhu, skill string, stype string) string {
	var returnNum = -1
	ua, exist := pc.Ability[strings.ToUpper(skill)]
	if exist == true {
		switch stype {
		case "init":
			returnNum = (*ua).Init
			break
		case "sum":
			returnNum = (*ua).Sum
			break
		case "now":
			returnNum = (*ua).Now
			break
		default:
			break
		}
	} else {
		for _, s := range pc.Skill {
			if strings.Contains((*s).Name+" "+(*s).Sub, skill) == true {
				switch stype {
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
func AddSkillNum(pc *CharacterOfCthulhu, skill string, add string) string {
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
