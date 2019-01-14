package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/c-ardinal/Nodens/lib/cthulhu"
	"github.com/c-ardinal/Nodens/lib/nodens"
)

// cmdShowVersion バージョン情報確認ハンドラ
func cmdShowVersion(opt []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string
	verResult, err := nodens.ExecuteVersionCheck(nodens.GetConfig().EndPoint)
	if err != nil {
		return "Reference error.", "", nil
	}
	returnMes = "\r\n[Nodens] " + nodens.GetVersion()
	returnMes += "\r\n[Main] " + mainVersion
	returnMes += "\r\n[API] " + verResult.API
	returnMes += "\r\n[BCDice] " + verResult.BCDice
	return returnMes, "", nil
}

// cmdCreateSession 親セッション生成ハンドラ
func cmdCreateSession(opt []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string
	if cthulhu.CheckDuplicateSession(ch.ID) == true {
		if opt[0] == "--forced" {
			if opt[1] != "" {
				cthulhu.RemoveSession(ch.ID)
				cs = cthulhu.NewSession(ch.ID, opt[1], mes.Author.Username, mes.Author.ID)
				returnMes = "Session recreate successfully. (System: " + cs.Scenario.System + ", ID: " + ch.ID + ")"
			} else {
				returnMes = "Session create failed."
			}
		} else {
			returnMes = "Session already exists."
		}
	} else {
		if opt[0] != "" {
			cs = cthulhu.NewSession(ch.ID, opt[0], mes.Author.Username, mes.Author.ID)
			returnMes = "Session create successfully. (System: " + cs.Scenario.System + ", ID: " + ch.ID + ")"
		} else {
			returnMes = "Session create failed."
		}
	}
	return returnMes, "", nil
}

// cmdConnectSession 親セッション接続ハンドラ
func cmdConnectSession(opt []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string
	if cthulhu.CheckDuplicateSession(opt[0]) == true {
		if opt[0] != ch.ID {
			pcs := cthulhu.GetSessionByID(opt[0])
			(*pcs).Discord.Child = append((*pcs).Discord.Child, cthulhu.NaID{ID: ch.ID})
			returnMes = "Parent session connect successfully. (Parent: " + (*pcs).Discord.Parent.ID + ", Child: " + ch.ID + ")"
		} else {
			returnMes = "Invalid session id."
		}
	} else {
		returnMes = "Parent session not found."
	}
	return returnMes, "", nil
}
