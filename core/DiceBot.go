package core

import (
	"log"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
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

// ダイスロール実行ログ格納変数
var diceResultLogs = []DiceResultLog{}

// テキストコマンドハンドル群登録用マップ
var cmdHandleMap = map[string]map[string]CmdHandleFunc{}

// キャラクターデータ取得関数群登録用マップ
var cdGetFuncMap = map[string]map[string]CharacterDataGetFunc{}

// セッション復元関数群登録用マップ
var sessionRestoreFuncTable = map[string]SessionRestoreFunc{}

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/

// AddCmdHandler テキストコマンドハンドラ登録処理
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
		handlerResult = ExecuteCmdHandlerGeneral(md, cmdHandleMap)
	}
	return handlerResult
}

// ExecuteCmdHandlerGeneral コマンドハンドラ共通処理
func ExecuteCmdHandlerGeneral(md MessageData, handleMap map[string]map[string]CmdHandleFunc) (handlerResult HandlerResult) {
	/* チャネルIDからセッション情報を取得 */
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

	if fg, exist := handleMap["General"][md.Command]; exist {
		/* 共通コマンドコール処理 */
		log.Printf("[Event]: Command call '%v'", md.Command)
		handlerResult = fg.ExecuteCmd(cs, md)
	} else if fs, exist := handleMap[system][md.Command]; exist {
		/* 各システム用コマンドコール処理 */
		log.Printf("[Event]: Command call '%v'", md.Command)
		handlerResult = fs.ExecuteCmd(cs, md)
	} else {
		/* セッションが生成されている場合のみダイスロールを実行 */
		if CheckExistParentSession(targetID) {
			var rollResult BCDiceRollResult
			log.Printf("[Event]: Execute dice roll '%v'", md.MessageString)
			rollResult, handlerResult.Error = ExecuteDiceRollAndCalc(GetConfig().EndPoint, (*cs).Scenario.System, md.MessageString)
			if handlerResult.Error != nil {
				/* 有効にするメッセージタイプ */
				handlerResult.Normal.EnableType = EnEmbed
				/* テキストメッセージ */
				handlerResult.Normal.Content = "Error: Dice roll failed > " + handlerResult.Error.Error()
				/* Embedメッセージ */
				handlerResult.Normal.Embed = &discordgo.MessageEmbed{
					Description: handlerResult.Normal.Content,
					Color:       EnColorRed,
				}
			} else {
				if rollResult.Result != "" {
					/* 有効にするメッセージタイプ */
					handlerResult.Normal.EnableType = EnEmbed
					/* テキストメッセージ */
					handlerResult.Normal.Content = rollResult.Result
					/* Embedメッセージ */
					handlerResult.Normal.Embed = &discordgo.MessageEmbed{
						Description: handlerResult.Normal.Content,
						Color:       EnColorGreen,
					}
				}
			}

			if rollResult.Secret {
				/* シークレットダイスが振られた旨のメッセージ */
				/* 有効にするメッセージタイプ */
				handlerResult.Secret.EnableType = EnEmbed
				/* テキストメッセージ */
				handlerResult.Secret.Content = "**SECRET DICE**"
				/* Embedメッセージ */
				handlerResult.Secret.Embed = &discordgo.MessageEmbed{
					Title: "SECRET DICE",
					Color: EnColorYellow,
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
		} else {
			log.Printf("[Debug]: No event. System: '%v', MessageData: '%v'", system, md)
		}
	}
	return handlerResult
}

// キャラクターデータ取得関数登録処理
func AddCharacterDataGetFunc(system string, dataName string, getFunc CharacterDataGetFunc) {
	if cdGetFuncMap[system] == nil {
		cdGetFuncMap[system] = make(map[string]CharacterDataGetFunc)
	}
	cdGetFuncMap[system][dataName] = getFunc
}

// セッション復元関数群登録処理
func SetRestoreFunc(funcmap map[string]SessionRestoreFunc) {
	sessionRestoreFuncTable = funcmap
}

// セッション復元関数群公開処理
func GetRestoreFunc() map[string]SessionRestoreFunc {
	return sessionRestoreFuncTable
}

// キャラクターデータ取得関数登録処理
func GetCharacterDataGetFunc(system string, dataName string) CharacterDataGetFunc {
	var result CharacterDataGetFunc = nil
	if _, exist := cdGetFuncMap[system]; exist {
		if _, exist := cdGetFuncMap[system][dataName]; exist {
			result = cdGetFuncMap[system][dataName]
		}
	}
	return result
}

// ダイスロール実行ログ取得処理
func GetDiceResultLogs() []DiceResultLog {
	return diceResultLogs
}
