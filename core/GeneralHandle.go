package core

import (
	"github.com/bwmarrin/discordgo"
)

// CmdShowVersion バージョン情報確認ハンドラ
func CmdShowVersion(opt []string, cs *Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string
	/* BCDiceAPIのバージョン情報取得 */
	verResult, err := ExecuteVersionCheck(GetConfig().EndPoint)
	if err != nil {
		return "Reference error.", "", nil
	}
	/* 返却メッセージの構築 */
	returnMes = "\r\n[Nodens] " + GetVersion()
	//returnMes += "\r\n[Main] " + mainVersion
	returnMes += "\r\n[API] " + verResult.API
	returnMes += "\r\n[BCDice] " + verResult.BCDice
	return returnMes, "", nil
}

// CmdCreateSession 親セッション生成ハンドラ
func CmdCreateSession(opt []string, cs *Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string
	/* セッションの存在有無を確認 */
	if CheckExistSession(ch.ID) == true {
		opt := string(opt[0])
		sys := string(opt[1])
		if opt == "--forced" {
			if sys != "" {
				/* セッションの強制再生成実行 */
				RemoveSession(ch.ID)
				cs = NewSession(ch.ID, sys, mes.Author.Username, mes.Author.ID)
				returnMes = "Session recreate successfully. (System: " + cs.Scenario.System + ", ID: " + ch.ID + ")"
			} else {
				/* システムの指定が無いの場合はセッションの強制再生成しない */
				returnMes = "Session create failed."
			}
		} else {
			/* セッションが生成済みなら生成しない */
			returnMes = "Session already exists."
		}
	} else {
		sys := string(opt[0])
		if sys != "" {
			/* セッションを生成 */
			cs = NewSession(ch.ID, sys, mes.Author.Username, mes.Author.ID)
			returnMes = "Session create successfully. (System: " + cs.Scenario.System + ", ID: " + ch.ID + ")"
		} else {
			/* システムの指定が無いの場合はセッションを生成しない */
			returnMes = "Session create failed."
		}
	}
	return returnMes, "", nil
}

// CmdConnectSession 親セッション接続ハンドラ
func CmdConnectSession(opt []string, cs *Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string
	/* 親セッションの存在有無確認 */
	parentID := string(opt[0])
	if CheckExistSession(parentID) == true {
		if parentID != ch.ID {
			/* 自セッションと親セッションが異なるセッションなら接続 */
			pcs := GetSessionByID(parentID)
			(*pcs).Discord.Child = append((*pcs).Discord.Child, NaID{ID: ch.ID})
			returnMes = "Parent session connect successfully. (Parent: " + (*pcs).Discord.Parent.ID + ", Child: " + ch.ID + ")"
		} else {
			/* 親・自セッションが同一の場合は接続しない */
			returnMes = "Invalid session id."
		}
	} else {
		/* 親セッションが存在しない場合は接続しない */
		returnMes = "Parent session not found."
	}
	return returnMes, "", nil
}
