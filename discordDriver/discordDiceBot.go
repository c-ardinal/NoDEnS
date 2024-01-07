package discordDriver

import (
	"log"
	"Nodens/core"
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

// スラッシュコマンドハンドル群登録用マップ
var slashCmdHandleMap = map[string]map[string]core.CmdHandleFunc{}

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/

// スラッシュコマンドハンドラ登録処理
func AddSlashCmdHandler(system string, cmd string, handler core.CmdHandleFunc) {
	if slashCmdHandleMap[system] == nil {
		slashCmdHandleMap[system] = make(map[string]core.CmdHandleFunc)
	}
	slashCmdHandleMap[system][cmd] = handler
}

// ExecuteCmdHandler スラッシュコマンドハンドラ実行処理
func ExecuteSlashCmdHandler(md core.MessageData) (handlerResult core.HandlerResult) {
	log.Printf("[Event]: Execute slash command handler '%v'", md)
	handlerResult = core.ExecuteCmdHandlerGeneral(md, slashCmdHandleMap)
	return handlerResult
}

// TRPGセッション生成処理
func CmdCreateSession(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var system string
	for _, opt := range md.Options {
		if opt.Name == "system" {
			system = opt.Value
		}
	}
	handlerResult = core.CmdCreateSession(cs, md)
	if handlerResult.Error == nil {
		JobRegistriesAppCommands(system, md.GuildID)
		registeredGuildIds = append(registeredGuildIds, md.GuildID)
	}
	return handlerResult
}

// TRPGセッション復元処理
func CmdRestoreSession(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	handlerResult = core.CmdRestoreSession(cs, md)
	if handlerResult.Error == nil {
		JobRegistriesAppCommands(core.GetSessionByID(md.ChannelID).Scenario.System, md.GuildID)
		registeredGuildIds = append(registeredGuildIds, md.GuildID)
	}
	return handlerResult
}