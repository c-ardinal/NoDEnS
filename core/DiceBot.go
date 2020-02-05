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

// CmdHandleFunc コマンドハンドラ型
type CmdHandleFunc func(opts []string, cs *Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error)

// ExecuteCmd CmdHandleFunc実行処理
func (f CmdHandleFunc) ExecuteCmd(opts []string, cs *Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	return f(opts, cs, ch, mes)
}

// NodensVersion nodensパッケージ&cthulhuパッケージのバージョン情報
const NodensVersion string = "0.0.2"

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
func ExecuteCmdHandler(ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string
	var secretMes string
	var err error

	returnMes = mes.Content

	commandRegex := regexp.MustCompile("[^ ]+")
	if commandRegex.MatchString(mes.Content) {
		var targetID string
		/* DiscordのチャネルIDからセッション情報を取得 */
		if GetParentIDFromChildID(ch.ID) != "" {
			targetID = GetParentIDFromChildID(ch.ID)
		} else {
			targetID = ch.ID
		}
		cs := GetSessionByID(targetID)

		/* コマンド部と引数部を分解 */
		cmd := commandRegex.FindAllString(mes.Content, -1)[0]
		opts := commandRegex.FindAllString(mes.Content, -1)[1:]

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
			returnMes, secretMes, err = fg.ExecuteCmd(opts, cs, ch, mes)
		} else if isSysExist == true {
			/* 各システム用コマンドコール処理 */
			returnMes, secretMes, err = fs.ExecuteCmd(opts, cs, ch, mes)
		} else {
			/* セッションが生成されている場合のみダイスロールを実行 */
			if CheckExistSession(targetID) {
				rollResult, _ := ExecuteDiceRoll(GetConfig().EndPoint, (*cs).Scenario.System, mes.Content)
				returnMes, err = rollResult.Result, nil
				if rollResult.Secret == true {
					/* シークレットダイスが振られた旨のメッセージ */
					secretMes = "**SECRET DICE**"
				} else {
					/* シークレットダイス以外の実行結果を記録 */
					const format = "2006/01/02_15:04:05"
					parsedTime, _ := mes.Timestamp.Parse()
					var diceResultLog DiceResultLog

					diceResultLog.Player.ID = mes.Author.ID
					diceResultLog.Player.Name = mes.Author.Username
					diceResultLog.Time = parsedTime.Format(format)
					diceResultLog.Command = mes.Content
					diceResultLog.Result = returnMes
					diceResultLogs = append(diceResultLogs, diceResultLog)
				}
			}
		}
	}

	return returnMes, secretMes, err
}

// GetVersion バージョン情報取得処理
func GetVersion() string {
	return NodensVersion
}

// GetDiceResultLogs 代スクロール実行ログ取得処理
func GetDiceResultLogs() []DiceResultLog {
	return diceResultLogs
}
