package cthulhu

import "Nodens/core"

/****************************************************************************/
/* 外部公開型定義                                                           */
/****************************************************************************/

// キャラクター情報構造体
type CharacterOfCthulhu struct {
	Player     core.NaID           `json:"player"`
	URL        string              `json:"url"`
	Md         string              `json:"md"`
	ID         string              `json:"id"`
	Personal   Personal            `json:"personal-data"`
	Ability    map[string]*Ability `json:"ability"`
	Skill      map[string]*Skill   `json:"skill"`
	Battle     Battle              `json:"battle"`
	Belongings Belongings          `json:"belongings"`
	Memo       string              `json:"memo"`
}

// キャラクター個人情報構造体
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

// キャラクター能力情報構造体
type Ability struct {
	Name  string `json:"name"`
	Init  int    `json:"init"`
	Add   int    `json:"add"`
	Temp  int    `json:"temp"`
	Sum   int    `json:"sum"`
	Start int    `json:"start"`
	Now   int    `json:"now"`
}

// キャラクター技能情報構造体
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

// キャラクター戦闘情報構造体
type Battle struct {
	Db     string              `json:"db"`
	Weapon map[string]*UnitArm `json:"weapon"`
	Armor  map[string]*UnitArm `json:"armor"`
}

// 戦闘情報構造体
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

// キャラクター所持品情報構造体
type Belongings struct {
	Pocket int              `json:"pocket-money"`
	Other  int              `json:"other-money"`
	Items  map[string]*Item `json:"items"`
}

// アイテム情報構造体
type Item struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
	Other string `json:"other"`
}

/****************************************************************************/
/* 外部公開定数定義                                                         */
/****************************************************************************/
