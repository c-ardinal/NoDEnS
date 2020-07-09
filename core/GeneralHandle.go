package core

// CmdShowVersion バージョン情報確認ハンドラ
func CmdShowVersion(opt []string, cs *Session, md MessageData) (handlerResult HandlerResult) {
	/* BCDiceAPIのバージョン情報取得 */
	verResult, err := ExecuteVersionCheck(GetConfig().EndPoint)
	if err != nil {
		handlerResult.Normal.Content = "Reference error."
		handlerResult.Error = err
		return handlerResult
	}
	/* 返却メッセージの構築 */
	handlerResult.Normal.Content = "\r\n[Nodens] " + GetVersion()
	handlerResult.Normal.Content += "\r\n[API] " + verResult.API
	handlerResult.Normal.Content += "\r\n[BCDice] " + verResult.BCDice
	return handlerResult
}

// CmdCreateSession 親セッション生成ハンドラ
func CmdCreateSession(opt []string, cs *Session, md MessageData) (handlerResult HandlerResult) {
	option := string(opt[0])
	if option == "--forced" {
		sys := string(opt[1])
		isContains := CheckContainsSystem(GetConfig().EndPoint, sys)
		if sys != "" && isContains == true {
			if CheckExistSession(md.ChannelID) == true {
				/* セッションの強制再生成実行 */
				RemoveSession(md.ChannelID)
				cs = NewSession(md.ChannelID, sys, md.AuthorName, md.AuthorID)
				handlerResult.Normal.Content = "Session recreate successfully. (System: " + cs.Scenario.System + ", ID: " + md.ChannelID + ")"
			} else {
				/* セッションを生成 */
				cs = NewSession(md.ChannelID, sys, md.AuthorName, md.AuthorID)
				handlerResult.Normal.Content = "Session create successfully. (System: " + cs.Scenario.System + ", ID: " + md.ChannelID + ")"
			}
		} else {
			/* システムの指定が無いの場合はセッションの強制再生成しない */
			handlerResult.Normal.Content = "Session create failed."
		}
	} else {
		sys := string(opt[0])
		isContains := CheckContainsSystem(GetConfig().EndPoint, sys)
		if sys != "" && isContains == true {
			if CheckExistSession(md.ChannelID) == true {
				/* セッションが生成済みなら生成しない */
				handlerResult.Normal.Content = "Session already exists."
			} else {
				/* セッションを生成 */
				cs = NewSession(md.ChannelID, sys, md.AuthorName, md.AuthorID)
				handlerResult.Normal.Content = "Session create successfully. (System: " + cs.Scenario.System + ", ID: " + md.ChannelID + ")"
			}
		} else {
			/* システムの指定が無いの場合はセッションを生成しない */
			handlerResult.Normal.Content = "Session create failed."
		}
	}
	return handlerResult
}

// CmdConnectSession 親セッション接続ハンドラ
func CmdConnectSession(opt []string, cs *Session, md MessageData) (handlerResult HandlerResult) {
	/* 親セッションの存在有無確認 */
	parentID := string(opt[0])
	if CheckExistSession(parentID) == true {
		if parentID != md.ChannelID {
			/* 自セッションと親セッションが異なるセッションなら接続 */
			pcs := GetSessionByID(parentID)
			(*pcs).Discord.Child = append((*pcs).Discord.Child, NaID{ID: md.ChannelID})
			handlerResult.Normal.Content = "Parent session connect successfully. (Parent: " + (*pcs).Discord.Parent.ID + ", Child: " + md.ChannelID + ")"
		} else {
			/* 親・自セッションが同一の場合は接続しない */
			handlerResult.Normal.Content = "Invalid session id."
		}
	} else {
		/* 親セッションが存在しない場合は接続しない */
		handlerResult.Normal.Content = "Parent session not found."
	}
	return handlerResult
}

// CmdStoreSession セッション保存ハンドラ
func CmdStoreSession(opt []string, cs *Session, md MessageData) (handlerResult HandlerResult) {
	if CheckExistSession(md.ChannelID) == true {
		_, err := StoreSession(md.ChannelID)
		if err != nil {
			handlerResult.Normal.Content = "Session store failed."
			handlerResult.Error = err
		} else {
			handlerResult.Normal.Content = "Session store successfully."
		}
	} else {
		handlerResult.Normal.Content = "Session not created."
	}

	return handlerResult
}

// CmdLoadSession セッション復帰ハンドラ
//coreパッケージからPC・NPC情報の構造復元が出来ないため，システム側で実装する
