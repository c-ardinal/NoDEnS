package core

import (
	"log"
	"regexp"
	"strconv"

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
type CmdHandleFunc func(cs *Session, md MessageData) (handlerResult HandlerResult)

// ExecuteCmd CmdHandleFunc実行処理
func (f CmdHandleFunc) ExecuteCmd(cs *Session, md MessageData) (handlerResult HandlerResult) {
	return f(cs, md)
}

// CharacterDataGetFunc キャラデータ取得関数型
type CharacterDataGetFunc func(cd interface{}) string

// ExecuteCharacterDataGet CharacterDataGetFunc実行処理
func (f CharacterDataGetFunc) ExecuteCharacterDataGet(cd interface{}) string {
	return f(cd)
}

// SessionRestoreFunc セッション復元関数型
type SessionRestoreFunc func(ses *Session) bool

// ExecuteSessionRestore SessionRestoreFunc実行処理
func (f SessionRestoreFunc) ExecuteSessionRestore(ses *Session) bool {
	return f(ses)
}

// diceResultLogs ダイスロール実行ログ格納変数
var diceResultLogs = []DiceResultLog{}

// cmdHandleMap コマンドハンドル群登録用マップ
var cmdHandleMap = map[string]map[string]CmdHandleFunc{}

// cdGetFuncMap キャラデータ取得関数群登録用マップ
var cdGetFuncMap = map[string]map[string]CharacterDataGetFunc{}

// SessionRestoreFuncTable セッション復元関数群登録用マップ
var SessionRestoreFuncTable = map[string]SessionRestoreFunc{}

// AddCmdHandler コマンドハンドラ登録処理
func AddCmdHandler(system string, cmd string, handler CmdHandleFunc) {
	if cmdHandleMap[system] == nil {
		cmdHandleMap[system] = make(map[string]CmdHandleFunc)
	}
	cmdHandleMap[system][cmd] = handler
}

// ExecuteCmdHandler コマンドハンドラ実行処理
func ExecuteCmdHandler(md MessageData) (handlerResult HandlerResult) {
	log.Printf("[Event]: Execute text command handler '%v'", md)
	commandRegex := regexp.MustCompile("[^ ]+")
	if commandRegex.MatchString(md.MessageString) {
		/* コマンド部と引数部を分解 */
		md.Command = commandRegex.FindAllString(md.MessageString, -1)[0]
		for i, str := range commandRegex.FindAllString(md.MessageString, -1)[1:] {
			md.Options = append(md.Options, CommandOption{
				Name:  strconv.Itoa(i),
				Value: str,
			})
		}
		handlerResult = executeCmdHandlerGeneral(md, cmdHandleMap)
	}
	return handlerResult
}

// executeCmdHandlerGeneral コマンドハンドラ共通処理
func executeCmdHandlerGeneral(md MessageData, handleMap map[string]map[string]CmdHandleFunc) (handlerResult HandlerResult) {
		/* DiscordのチャネルIDからセッション情報を取得 */
	var targetID string = ""
		if GetParentIDFromChildID(md.ChannelID) != "" {
			targetID = GetParentIDFromChildID(md.ChannelID)
		} else {
			targetID = md.ChannelID
		}
		cs := GetSessionByID(targetID)

		/* コマンド部をもとにコール対象の関数ポインタを取得 */
	var system string = ""
		if cs != nil {
		system = cs.Scenario.System
		}

	if fg, exist := handleMap["General"][md.Command]; exist == true {
			/* 共通コマンドコール処理 */
		log.Printf("[Event]: Command call '%v'", md.Command)
		handlerResult = fg.ExecuteCmd(cs, md)
	} else if fs, exist := handleMap[system][md.Command]; exist == true {
			/* 各システム用コマンドコール処理 */
		log.Printf("[Event]: Command call '%v'", md.Command)
		handlerResult = fs.ExecuteCmd(cs, md)
		} else {
			/* セッションが生成されている場合のみダイスロールを実行 */
			if CheckExistSession(targetID) {
				/* 有効にするメッセージタイプ */
				handlerResult.Normal.EnableType = EnEmbed
				handlerResult.Secret.EnableType = EnEmbed

				var rollResult BCDiceRollResult
			log.Printf("[Event]: Execute dice roll '%v'", md.MessageString)
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
					handlerResult.Secret.Embed = &discordgo.MessageEmbed{
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
	return handlerResult
}

// AddCharacterDataGetFunc キャラデータ取得関数登録処理
func AddCharacterDataGetFunc(system string, dataName string, getFunc CharacterDataGetFunc) {
	if cdGetFuncMap[system] == nil {
		cdGetFuncMap[system] = make(map[string]CharacterDataGetFunc)
	}
	cdGetFuncMap[system][dataName] = getFunc
}

// SetRestoreFunc セッション復元関数群登録処理
func SetRestoreFunc(funcmap map[string]SessionRestoreFunc) {
	SessionRestoreFuncTable = funcmap
}

// GetCharacterDataGetFuncキャラデータ取得関数登録処理
func GetCharacterDataGetFunc(system string, dataName string) CharacterDataGetFunc {
	var result CharacterDataGetFunc = nil
	if _, exist := cdGetFuncMap[system]; exist == true {
		if _, exist := cdGetFuncMap[system][dataName]; exist == true {
			result = cdGetFuncMap[system][dataName]
		}
	}
	return result
}

// GetDiceResultLogs 代スクロール実行ログ取得処理
func GetDiceResultLogs() []DiceResultLog {
	return diceResultLogs
}
