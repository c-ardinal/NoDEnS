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

// MessageData ハンドラに渡すメッセージのデータ
type MessageData struct {
	ChannelID     string
	MessageID     string
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
	Embed      *discordgo.MessageEmbed
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

// CharacterDataGetFunc コマンドハンドラ型
type CharacterDataGetFunc func(cd interface{}) string

// ExecuteCharacterDataGet CharacterDataGetFunc実行処理
func (f CharacterDataGetFunc) ExecuteCharacterDataGet(cd interface{}) string {
	return f(cd)
}

// diceResultLogs ダイスロール実行ログ格納変数
var diceResultLogs = []DiceResultLog{}

// cmdHandleMap コマンドハンドル群登録用マップ
var cmdHandleMap = map[string]map[string]CmdHandleFunc{}

// cdGetFuncMap キャラデータ取得関数群登録用マップ
var cdGetFuncMap = map[string]map[string]CharacterDataGetFunc{}

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
				/* 有効にするメッセージタイプ */
				handlerResult.Normal.EnableType = EnEmbed

				var rollResult BCDiceRollResult
				rollResult, handlerResult.Error = ExecuteDiceRollAndCalc(GetConfig().EndPoint, (*cs).Scenario.System, md.MessageString)
				if handlerResult.Error != nil {
					/* テキストメッセージ */
					handlerResult.Normal.Content = "Error: " + handlerResult.Error.Error()
					/* Embedメッセージ */
					handlerResult.Normal.Embed = &discordgo.MessageEmbed{
						Description: handlerResult.Normal.Content,
						Color:       0xff0000, // Red
					}
				} else {
					/* テキストメッセージ */
					handlerResult.Normal.Content = rollResult.Result
					/* Embedメッセージ */
					handlerResult.Normal.Embed = &discordgo.MessageEmbed{
						Description: handlerResult.Normal.Content,
						Color:       0x00ff00, // Green
					}
				}

				if rollResult.Secret == true {
					/* シークレットダイスが振られた旨のメッセージ */
					/* テキストメッセージ */
					handlerResult.Secret.Content = "**SECRET DICE**"
					/* Embedメッセージ */
					handlerResult.Normal.Embed = &discordgo.MessageEmbed{
						Title: "SECRET DICE",
						Color: 0xffff00, // Yellow
					}
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

// AddCharacterDataGetFunc キャラデータ取得関数登録処理
func AddCharacterDataGetFunc(system string, dataName string, getFunc CharacterDataGetFunc) {
	if cdGetFuncMap[system] == nil {
		cdGetFuncMap[system] = make(map[string]CharacterDataGetFunc)
	}
	cdGetFuncMap[system][dataName] = getFunc
}

// GetCharacterDataGetFuncキャラデータ取得関数登録処理
func GetCharacterDataGetFunc(system string, dataName string) CharacterDataGetFunc {
	var result CharacterDataGetFunc = nil
	_, systemExist := cdGetFuncMap[system]
	if systemExist == true {
		_, dnExist := cdGetFuncMap[system][dataName]
		if dnExist == true {
			result = cdGetFuncMap[system][dataName]
		}
	}
	return result
}

// GetDiceResultLogs 代スクロール実行ログ取得処理
func GetDiceResultLogs() []DiceResultLog {
	return diceResultLogs
}
