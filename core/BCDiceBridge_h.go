package core

/****************************************************************************/
/* 外部公開型定義                                                           */
/****************************************************************************/

// BCDiceバージョン情報格納構造体
type BCDiceVersionResult struct {
	API    string `json:"api"`
	BCDice string `json:"bcdice"`
}

// BCDiceシステム一覧取得結果格納構造体
type BCDiceSystemsResult struct {
	Systems []BCDiceSystem `json:"game_system"`
}

// BCDiceシステム情報格納構造体
type BCDiceSystem struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	SortKey string `json:"sort_key"`
}

// ダイスロール実行結果格納構造体
type BCDiceRollResult struct {
	Ok       bool   `json:"ok"`
	Result   string `json:"text"`
	Secret   bool   `json:"secret"`
	Success  bool   `json:"success"`
	Failure  bool   `json:"failure"`
	Critical bool   `json:"critical"`
	Fumble   bool   `json:"fumble"`
	Dices    []Dice `json:"rands"`
}

// ダイス情報格納構造体
type Dice struct {
	Kind  string `json:"kind"`
	Faces int    `json:"sides"`
	Value int    `json:"value"`
}

/****************************************************************************/
/* 外部公開定数定義                                                         */
/****************************************************************************/
