package core

/****************************************************************************/
/* 外部公開型定義                                                           */
/****************************************************************************/

// セッションデータ構造体
type Session struct {
	Chat     Chat                   `json:"chat"`
	Scenario Scenario               `json:"scenario"`
	Pc       map[string]interface{} `json:"pc"`
	Npc      map[string]interface{} `json:"npc"`
}

// Chat情報構造体
type Chat struct {
	Parent NaID   `json:"text"`
	Child  []NaID `json:"private"`
	Voice  NaID   `json:"voice"`
}

// シナリオ情報構造体
type Scenario struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Keeper  NaID   `json:"keeper"`
	Chatlog string `json:"chatlog"`
}

// 名前&ID紐づけ構造体
type NaID struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

/****************************************************************************/
/* 外部公開定数定義                                                         */
/****************************************************************************/
