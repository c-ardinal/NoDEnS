package core

import "github.com/bwmarrin/discordgo"

/****************************************************************************/
/* 外部公開型定義                                                           */
/****************************************************************************/

// ダイスロール実行ログ型
type DiceResultLog struct {
	Player  NaID
	Time    string
	Command string
	Result  string
}

// ハンドラに渡すメッセージ型
type MessageData struct {
	ChannelID     string
	GuildID       string
	MessageID     string
	AuthorID      string
	AuthorName    string
	MessageString string
	Command       string
	Options       []CommandOption
}

// コマンドオプション情報
type CommandOption struct {
	Name  string
	Value string
}

// ハンドラの戻りオブジェクト
type HandlerResult struct {
	Normal MessageTemplate
	Secret MessageTemplate
	Error  error
}

// ユーザに返すメッセージの共通型
type MessageTemplate struct {
	EnableType int
	Content    string
	Embed      *discordgo.MessageEmbed
}

// コマンドハンドラ型
type CmdHandleFunc func(cs *Session, md MessageData) (handlerResult HandlerResult)

// CmdHandleFunc実行処理
func (f CmdHandleFunc) ExecuteCmd(cs *Session, md MessageData) (handlerResult HandlerResult) {
	return f(cs, md)
}

// キャラクターデータ取得関数型
type CharacterDataGetFunc func(cd interface{}) string

// CharacterDataGetFunc実行処理
func (f CharacterDataGetFunc) ExecuteCharacterDataGet(cd interface{}) string {
	return f(cd)
}

// セッション復元関数型
type SessionRestoreFunc func(ses *Session) bool

// SessionRestoreFunc実行処理
func (f SessionRestoreFunc) ExecuteSessionRestore(ses *Session) bool {
	return f(ses)
}

/****************************************************************************/
/* 外部公開定数定義                                                         */
/****************************************************************************/

// メッセージ応答フラグ
const (
	// メッセージを応答しない(デフォルト)
	EnNoMessage int = 0
	// 文字によるメッセージ応答を行う
	EnContent int = 1
	// Embedによるメッセージ応答を行う
	EnEmbed int = 2
)

// カラーコード定義
const (
	EnColorRed    int = 0xFF0000
	EnColorYellow int = 0xFFFF00
	EnColorGreen  int = 0x00FF00
)
