package core

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

// DiceResultLog ダイスロール実行ログ型
type DiceResultLog struct {
	Player  NaID
	Time    string
	Command string
	Result  string
}

// diceResultLogs ダイスロール実行ログ格納変数
var diceResultLogs = []DiceResultLog{}

// MessageData ハンドラに渡すメッセージのデータ
type MessageData struct {
	ChannelID     string
	AuthorID      string
	AuthorName    string
	MessageString string
}

// HandlerResult ハンドラの戻りオブジェクト
type HandlerResult struct {
	Normal MessageTemplate
	Secret MessageTemplate
	Error  error
}

// MessageTemplate ユーザに返すメッセージの共通型
type MessageTemplate struct {
	EnableType int
	Content    string
	Embed      discordgo.MessageEmbed
}

const (
	//EnContent 文字によるメッセージ返却を行う(デフォルト)
	EnContent int = 0
	//EnEmbed Embedによるメッセージ返却を行う
	EnEmbed int = 1
)

// CmdHandleFunc コマンドハンドラ型
type CmdHandleFunc func(opts []string, cs *Session, md MessageData) (handlerResult HandlerResult)

// ExecuteCmd CmdHandleFunc実行処理
func (f CmdHandleFunc) ExecuteCmd(opts []string, cs *Session, md MessageData) (handlerResult HandlerResult) {
	return f(opts, cs, md)
}

// NodensVersion nodensパッケージ&cthulhuパッケージのバージョン情報
const NodensVersion string = "0.2.3"

// cmdHandleMap コマンドハンドル群登録用マップ
var cmdHandleMap = map[string]map[string]CmdHandleFunc{}

// AddCmdHandler コマンドハンドラ登録処理
func AddCmdHandler(system string, cmd string, handler CmdHandleFunc) {
	if cmdHandleMap[system] == nil {
		cmdHandleMap[system] = make(map[string]CmdHandleFunc)
	}
	cmdHandleMap[system][cmd] = handler
}

// ExecuteCmdHandler コマンドハンドラ実行処理
func ExecuteCmdHandler(md MessageData) (handlerResult HandlerResult) {
	handlerResult.Normal.Content = md.MessageString

	commandRegex := regexp.MustCompile("[^ ]+")
	if commandRegex.MatchString(md.MessageString) {
		var targetID string
		/* DiscordのチャネルIDからセッション情報を取得 */
		if GetParentIDFromChildID(md.ChannelID) != "" {
			targetID = GetParentIDFromChildID(md.ChannelID)
		} else {
			targetID = md.ChannelID
		}
		cs := GetSessionByID(targetID)

		/* コマンド部と引数部を分解 */
		cmd := commandRegex.FindAllString(md.MessageString, -1)[0]
		opts := commandRegex.FindAllString(md.MessageString, -1)[1:]

		/* コマンド部をもとにコール対象の関数ポインタを取得 */
		var fs CmdHandleFunc
		var isSysExist bool
		if cs != nil {
			system := cs.Scenario.System
			fs, isSysExist = cmdHandleMap[system][cmd]
		}
		fg, isGenExist := cmdHandleMap["General"][cmd]

		if isGenExist == true {
			/* 共通コマンドコール処理 */
			handlerResult = fg.ExecuteCmd(opts, cs, md)
		} else if isSysExist == true {
			/* 各システム用コマンドコール処理 */
			handlerResult = fs.ExecuteCmd(opts, cs, md)
		} else {
			/* セッションが生成されている場合のみダイスロールを実行 */
			if CheckExistSession(targetID) {
				var rollResult BCDiceRollResult
				rollResult, handlerResult.Error = ExecuteDiceRollAndCalc(GetConfig().EndPoint, (*cs).Scenario.System, md.MessageString)
				if handlerResult.Error != nil {
					handlerResult.Normal.Content = "Error: " + handlerResult.Error.Error()
				} else {
					handlerResult.Normal.Content = rollResult.Result
				}
				if rollResult.Secret == true {
					/* シークレットダイスが振られた旨のメッセージ */
					handlerResult.Secret.Content = "**SECRET DICE**"
				} else {
					/* シークレットダイス以外の実行結果を記録 */
					//const format = "2006/01/02_15:04:05"
					//parsedTime, _ := mes.Timestamp.Parse()
					var diceResultLog DiceResultLog

					diceResultLog.Player.ID = md.AuthorID
					diceResultLog.Player.Name = md.AuthorName
					//diceResultLog.Time = parsedTime.Format(format)
					diceResultLog.Command = md.MessageString
					diceResultLog.Result = handlerResult.Normal.Content
					diceResultLogs = append(diceResultLogs, diceResultLog)
				}
			}
		}
	}

	return handlerResult
}

// GetVersion バージョン情報取得処理
func GetVersion() string {
	return NodensVersion
}

// GetDiceResultLogs 代スクロール実行ログ取得処理
func GetDiceResultLogs() []DiceResultLog {
	return diceResultLogs
}
