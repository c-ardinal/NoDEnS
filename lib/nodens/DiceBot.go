package nodens

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/c-ardinal/Nodens/lib/cthulhu"
)

// CmdHandleFunc コマンドハンドラ型
type CmdHandleFunc func(opts []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error)

// ExecuteCmd CmdHandleFunc実行処理
func (f CmdHandleFunc) ExecuteCmd(opts []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	return f(opts, cs, ch, mes)
}

// NodensVersion nodensパッケージ&cthulhuパッケージのバージョン情報
const NodensVersion string = "0.0.1"

// cmdHandleMap コマンドハンドル群登録用マップ
var cmdHandleMap = map[string]CmdHandleFunc{}

// AddCmdHandler コマンドハンドラ登録処理
func AddCmdHandler(cmd string, handler CmdHandleFunc) {
	cmdHandleMap[cmd] = handler
}

// ExecuteCmdHandler コマンドハンドラ実行処理
func ExecuteCmdHandler(ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string = mes.Content
	var secretMes string
	var err error

	commandRegex := regexp.MustCompile("[^ ]+")
	if commandRegex.MatchString(mes.Content) {
		var targetID string
		if cthulhu.GetParentIDFromChildID(ch.ID) != "" {
			targetID = cthulhu.GetParentIDFromChildID(ch.ID)
		} else {
			targetID = ch.ID
		}
		cs := cthulhu.GetSessionByID(targetID)

		cmd := commandRegex.FindAllString(mes.Content, -1)[0]
		opts := commandRegex.FindAllString(mes.Content, -1)[1:]

		f, isExist := cmdHandleMap[cmd]
		if isExist == true {
			returnMes, secretMes, err = f.ExecuteCmd(opts, cs, ch, mes)
		} else {
			if cthulhu.CheckDuplicateSession(targetID) {
				rollResult, _ := ExecuteDiceRoll(GetConfig().EndPoint, (*cs).Scenario.System, mes.Content)
				returnMes, err = rollResult.Result, nil
				if rollResult.Secret == true {
					secretMes = "**SECRET DICE**"
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
