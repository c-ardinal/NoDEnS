package core

import (
	"github.com/bwmarrin/discordgo"
)

// CmdShowVersion バージョン情報確認ハンドラ
func CmdShowVersion(cs *Session, md MessageData) (handlerResult HandlerResult) {
	var returnMes string
	var returnMesColor int

	/* BCDiceAPIのバージョン情報取得 */
	verResult, err := ExecuteVersionCheck(GetConfig().EndPoint)
	if err != nil {
		returnMes = "Reference error."
		handlerResult.Error = err
		returnMesColor = 0xff0000 // Red
	} else {
		/* テキストメッセージ */
		returnMes = "\n[NoDEnS]: " + GetVersion()
		returnMes += "\n[BCDice-API]: " + verResult.API
		returnMes += "\n[BCDice]: " + verResult.BCDice
		returnMesColor = 0x00ff00 // Green
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = EnEmbed

	/* テキストメッセージ */
	handlerResult.Normal.Content = returnMes

	/* Embedメッセージ */
	handlerResult.Normal.Embed = &discordgo.MessageEmbed{
		Description: returnMes,
		Color:       returnMesColor,
	}

	return handlerResult
}

// CmdCreateSession 親セッション生成ハンドラ
func CmdCreateSession(cs *Session, md MessageData) (handlerResult HandlerResult) {
	var forced string
	var system string
	var returnMes string
	var returnMesColor int

	for _, opt := range md.Options {
		if opt.Name == "forced" {
			forced = opt.Value
		} else {
			if opt.Value != "OtherSystem" {
				system = opt.Value
			}
		}
	}

	if forced != "" {
		isContains := CheckContainsSystem(GetConfig().EndPoint, system)
		if system != "" {
			if isContains == true {
				if CheckExistSession(md.ChannelID) == true {
					/* セッションの強制再生成実行 */
					RemoveSession(md.ChannelID)
					cs = NewSession(md.ChannelID, system, md.AuthorName, md.AuthorID)
					returnMes = "Session recreate successfully. \n[System]: " + cs.Scenario.System + " \n[ChannelID]: " + md.ChannelID
					returnMesColor = 0x00ff00 // Green
				} else {
					/* セッションを生成 */
					cs = NewSession(md.ChannelID, system, md.AuthorName, md.AuthorID)
					returnMes = "Session create successfully. \n[System]: " + cs.Scenario.System + " \n[ChannelID]: " + md.ChannelID
					returnMesColor = 0x00ff00 // Green
				}
			} else {
				/* 指定されたシステムが見つからない場合はセッションの強制再生成をしない */
				returnMes = "System not found."
				returnMesColor = 0xff0000 // Red
			}
		} else {
			/* システムの指定が無い場合はセッションの強制再生成をしない */
			returnMes = "Session create failed."
			returnMesColor = 0xff0000 // Red
		}
	} else {
		isContains := CheckContainsSystem(GetConfig().EndPoint, system)
		if system != "" {
			if isContains == true {
				if CheckExistSession(md.ChannelID) == true {
					/* セッションが生成済みなら生成しない */
					returnMes = "Session already exists."
					returnMesColor = 0xffff00 // Yellow
				} else {
					/* セッションを生成 */
					cs = NewSession(md.ChannelID, system, md.AuthorName, md.AuthorID)
					returnMes = "Session create successfully. \n[System]: " + cs.Scenario.System + " \n[ChannelID]: " + md.ChannelID
					returnMesColor = 0x00ff00 // Green
				}
			} else {
				/* 指定されたシステムが見つからない場合はセッションの強制再生成をしない */
				returnMes = "System not found."
				returnMesColor = 0xff0000 // Red
			}
		} else {
			/* システムの指定が無い場合はセッションを生成しない */
			returnMes = "Session create failed."
			returnMesColor = 0xff0000 // Red
		}
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = EnEmbed

	/* テキストメッセージ */
	handlerResult.Normal.Content = returnMes

	/* Embedメッセージ */
	handlerResult.Normal.Embed = &discordgo.MessageEmbed{
		Description: returnMes,
		Color:       returnMesColor,
	}

	return handlerResult
}

// CmdConnectSession 親セッション接続ハンドラ
func CmdConnectSession(cs *Session, md MessageData) (handlerResult HandlerResult) {
	var returnMes string
	var returnMesColor int

	/* 親セッションの存在有無確認 */
	parentID := md.Options[0].Value
	if CheckExistSession(parentID) == true {
		if parentID != md.ChannelID {
			/* 自セッションと親セッションが異なるセッションなら接続 */
			pcs := GetSessionByID(parentID)
			(*pcs).Discord.Child = append((*pcs).Discord.Child, NaID{ID: md.ChannelID})
			returnMes = "Parent session connect successfully.\n[Parent]: " + (*pcs).Discord.Parent.ID + "\n[Child]: " + md.ChannelID + ")"
			returnMesColor = 0x00ff00 // Green
		} else {
			/* 親・自セッションが同一の場合は接続しない */
			returnMes = "Invalid session id."
			returnMesColor = 0xff0000 // Red
		}
	} else {
		/* 親セッションが存在しない場合は接続しない */
		returnMes = "Parent session not found."
		returnMesColor = 0xff0000 // Red
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = EnEmbed

	/* テキストメッセージ */
	handlerResult.Normal.Content = returnMes

	/* Embedメッセージ */
	handlerResult.Normal.Embed = &discordgo.MessageEmbed{
		Description: returnMes,
		Color:       returnMesColor,
	}

	return handlerResult
}

// CmdStoreSession セッション保存ハンドラ
func CmdStoreSession(cs *Session, md MessageData) (handlerResult HandlerResult) {
	var returnMes string
	var returnMesColor int

	if CheckExistSession(md.ChannelID) == true {
		_, err := StoreSession(md.ChannelID)
		if err != nil {
			returnMes = "Session store failed."
			handlerResult.Error = err
			returnMesColor = 0xff0000 // Red
		} else {
			returnMes = "Session store successfully."
			returnMesColor = 0x00ff00 // Green
		}
	} else {
		returnMes = "Session not created."
		returnMesColor = 0xff0000 // Red
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = EnEmbed

	/* テキストメッセージ */
	handlerResult.Normal.Content = returnMes

	/* Embedメッセージ */
	handlerResult.Normal.Embed = &discordgo.MessageEmbed{
		Description: returnMes,
		Color:       returnMesColor,
	}

	return handlerResult
}

// CmdRestoreSession セッション復元ハンドラ
func CmdRestoreSession(cs *Session, md MessageData) (handlerResult HandlerResult) {
	var returnMes string
	var returnMesColor int

	// 各システム共通情報の復元
	err := RestoreSession(md.ChannelID)
	if err != nil {
		returnMes = "Session load failed."
		returnMesColor = 0xff0000 // Red
		handlerResult.Error = err
	} else {
		// 各システム固有情報の復元
		ses := GetSessionByID(md.ChannelID)
		systemRestoreFunc, isExist := SessionRestoreFuncTable[(*ses).Scenario.System]
		if isExist == true {
			systemRestoreFunc.ExecuteSessionRestore(ses)
		}
		returnMes = "Session restore successfully."
		returnMesColor = 0x00ff00 // Green
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = EnEmbed

	/* テキストメッセージ */
	handlerResult.Normal.Content = returnMes

	/* Embedメッセージ */
	handlerResult.Normal.Embed = &discordgo.MessageEmbed{
		Description: returnMes,
		Color:       returnMesColor,
	}

	return handlerResult
}
