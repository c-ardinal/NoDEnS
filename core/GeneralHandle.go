package core

import (
	"errors"

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

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/

// バージョン情報確認ハンドラ
func CmdShowVersion(cs *Session, md MessageData) (handlerResult HandlerResult) {
	var returnMes string
	var returnMesColor int

	/* BCDiceAPIのバージョン情報取得 */
	verResult, err := ExecuteVersionCheck(GetConfig().EndPoint)
	if err != nil {
		returnMes = "Reference error."
		handlerResult.Error = err
		returnMesColor = EnColorRed
	} else {
		/* テキストメッセージ */
		returnMes = "\n[NoDEnS]: " + GetVersion()
		returnMes += "\n[BCDice-API]: " + verResult.API
		returnMes += "\n[BCDice]: " + verResult.BCDice
		returnMesColor = EnColorGreen
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

// 親セッション生成ハンドラ
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
	if CheckExistChildSession(md.ChannelID) {
		/* 子セッションが既に有る場合はセッションを生成しない */
		returnMes = "child session already exists"
		handlerResult.Error = errors.New(returnMes)
		returnMesColor = EnColorRed
	} else {
		if forced != "" {
			isContains := CheckContainsSystem(GetConfig().EndPoint, system)
			if system != "" {
				if isContains {
					if CheckExistParentSession(md.ChannelID) {
						/* セッションの強制再生成実行 */
						RemoveSession(md.ChannelID)
						cs = NewSession(md.ChannelID, system, md.AuthorName, md.AuthorID)
						returnMes = "session recreate successfully. \n[System]: " + cs.Scenario.System + " \n[ChannelID]: " + md.ChannelID
						returnMesColor = EnColorGreen
					} else {
						/* セッションを生成 */
						cs = NewSession(md.ChannelID, system, md.AuthorName, md.AuthorID)
						returnMes = "session create successfully. \n[System]: " + cs.Scenario.System + " \n[ChannelID]: " + md.ChannelID
						returnMesColor = EnColorGreen
					}
				} else {
					/* 指定されたシステムが見つからない場合はセッションの強制再生成をしない */
					returnMes = "system not found"
					handlerResult.Error = errors.New(returnMes)
					returnMesColor = EnColorRed
				}
			} else {
				/* システムの指定が無い場合はセッションの強制再生成をしない */
				returnMes = "session create failed"
				handlerResult.Error = errors.New(returnMes)
				returnMesColor = EnColorRed
			}
		} else {
			isContains := CheckContainsSystem(GetConfig().EndPoint, system)
			if system != "" {
				if isContains {
					if CheckExistParentSession(md.ChannelID) {
						/* セッションが生成済みなら生成しない */
						returnMes = "Session already exists."
						returnMesColor = EnColorYellow
					} else {
						/* セッションを生成 */
						cs = NewSession(md.ChannelID, system, md.AuthorName, md.AuthorID)
						returnMes = "session create successfully. \n[System]: " + cs.Scenario.System + " \n[ChannelID]: " + md.ChannelID
						returnMesColor = EnColorGreen
					}
				} else {
					/* 指定されたシステムが見つからない場合はセッションの強制再生成をしない */
					returnMes = "system not found"
					handlerResult.Error = errors.New(returnMes)
					returnMesColor = EnColorRed
				}
			} else {
				/* システムの指定が無い場合はセッションを生成しない */
				returnMes = "session create failed"
				handlerResult.Error = errors.New(returnMes)
				returnMesColor = EnColorRed
			}
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

// 親セッション接続ハンドラ
func CmdConnectSession(cs *Session, md MessageData) (handlerResult HandlerResult) {
	var returnMes string
	var returnMesColor int

	/* 親セッションの存在有無確認 */
	parentID := md.Options[0].Value
	if CheckExistParentSession(parentID) {
		if parentID != md.ChannelID {
			/* 自セッションと親セッションが異なるセッションなら接続 */
			pcs := GetSessionByID(parentID)
			(*pcs).Chat.Child = append((*pcs).Chat.Child, NaID{ID: md.ChannelID})
			returnMes = "parent session connect successfully.\n[Parent]: " + (*pcs).Chat.Parent.ID + "\n[Child]: " + md.ChannelID + ")"
			returnMesColor = EnColorGreen
		} else {
			/* 親・自セッションが同一の場合は接続しない */
			returnMes = "invalid session id"
			handlerResult.Error = errors.New(returnMes)
			returnMesColor = EnColorRed
		}
	} else {
		/* 親セッションが存在しない場合は接続しない */
		returnMes = "parent session not found"
		handlerResult.Error = errors.New(returnMes)
		returnMesColor = EnColorRed
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

// セッション保存ハンドラ
func CmdStoreSession(cs *Session, md MessageData) (handlerResult HandlerResult) {
	var returnMes string
	var returnMesColor int

	if CheckExistParentSession(md.ChannelID) {
		_, err := StoreSession(md.ChannelID)
		if err != nil {
			returnMes = "session store failed"
			handlerResult.Error = err
			returnMesColor = EnColorRed
		} else {
			returnMes = "session store successfully"
			returnMesColor = EnColorGreen
		}
	} else {
		returnMes = "session not created"
		handlerResult.Error = errors.New(returnMes)
		returnMesColor = EnColorRed
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

// セッション復元ハンドラ
func CmdRestoreSession(cs *Session, md MessageData) (handlerResult HandlerResult) {
	var returnMes string
	var returnMesColor int

	// 各システム共通情報の復元
	err := RestoreSession(md.ChannelID)
	if err != nil {
		returnMes = "session load failed"
		returnMesColor = EnColorRed
		handlerResult.Error = err
	} else {
		// 各システム固有情報の復元
		ses := GetSessionByID(md.ChannelID)
		systemRestoreFunc, isExist := GetRestoreFunc()[(*ses).Scenario.System]
		if isExist {
			systemRestoreFunc.ExecuteSessionRestore(ses)
		}
		returnMes = "session restore successfully"
		returnMesColor = EnColorGreen
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
