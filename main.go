package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/c-ardinal/Nodens/core"
	"github.com/c-ardinal/Nodens/system/cthulhu"

	"github.com/bwmarrin/discordgo"
)

// configFile デフォルト設定ファイルパス
var configFile = "SystemConfig.json"

// handleFuncStruct コマンドハンドラテーブル用構造体
type handleFuncStruct struct {
	system   string
	command  string
	function core.CmdHandleFunc
}

// handleFuncTable コマンドハンドラテーブル
var handleFuncTable = []handleFuncStruct{
	{"General", "ShowVersion", core.CmdShowVersion},       // バージョン情報表示処理
	{"General", "CreateSession", core.CmdCreateSession},   // 親セッション生成処理
	{"General", "ConnectSession", core.CmdConnectSession}, // 親セッション連携処理
	// TODO: あとで実装 {"General", "ExitSession", core.cmdExitSession},		// 親セッション連携処理
	// TODO: あとで実装 {"General", "RemoveSession", core.cmdRemoveSession},	// セッション削除処理
	{"General", "StoreSession", core.CmdStoreSession},                        // セッション保存処理
	{"General", "RestoreCthulhuSession", cthulhu.CmdRestoreSessionOfCthulhu}, // セッション復帰処理
	{"Cthulhu", "regchara", cthulhu.CmdRegistryCharacter},                    // キャラシート連携処理
	{"Cthulhu", "check", cthulhu.CmdCharaNumCheck},                           // 能力値確認処理
	{"Cthulhu", "ctrl", cthulhu.CmdCharaNumControl},                          // 能力値操作処理
	{"Cthulhu", "roll", cthulhu.CmdLinkRoll},                                 // 能力値ダイスロール処理
	{"Cthulhu", "Sroll", cthulhu.CmdSecretLinkRoll},                          // 能力値シークレットダイスロール処理
	{"Cthulhu", "sanc", cthulhu.CmdSanCheckRoll},                             // SAN値チェック処理
	// TODO: あとで実装 {"Cthulhu", "grow", cmdGrowRoll},    // 成長ロール処理
	// TODO: 実装中 {"Cthulhu", "showstat", cthulhu.CmdShowStatistics}, // ダイスロール統計表示処理
}

// main メイン関数
func main() {
	// 引数の読み込み
	if len(os.Args) != 1 {
		configFile = os.Args[1]
		_, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Println("Error => " + err.Error())
			os.Exit(1)
		}
	}

	// 設定ファイルの読み込み
	core.LoadConfig(configFile)

	// Discordのインスタンス生成
	discord, err := discordgo.New()
	discord.Token = core.GetConfig().BotToken
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// Discordのメッセージハンドラ登録
	discord.AddHandler(onMessageCreate)

	// コマンドハンドラ登録
	for _, handle := range handleFuncTable {
		core.AddCmdHandler(handle.system, handle.command, handle.function)
	}

	// セッション開始
	err = discord.Open()
	if err != nil {
		panic(err)
	}
	log.Println("Listening...")
	stopBot := make(chan bool)
	<-stopBot
	return
}

// onMessageCreate メッセージ受信時処理
func onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	var md core.MessageData
	md.ChannelID = message.ChannelID
	md.AuthorID = message.Author.ID
	md.AuthorName = message.Author.Username
	md.MessageString = message.Content

	var replyObject discordgo.MessageSend

	log.Printf("%20s %20s > %s\n", md.ChannelID, md.AuthorName, md.MessageString)

	if md.AuthorID != core.GetConfig().BotID {
		handlerResult := core.ExecuteCmdHandler(md)
		if handlerResult.Error != nil {
			log.Println(handlerResult.Error)
		}

		/* 通常メッセージの送信 */
		if handlerResult.Normal.EnableType == core.EnContent {
			if handlerResult.Normal.Content != "" && handlerResult.Normal.Content != md.MessageString {
				replyObject.Content = handlerResult.Normal.Content
				sendReplyMessage(session, md.ChannelID, md.AuthorID, replyObject)
			}
		} else if handlerResult.Normal.EnableType == core.EnEmbed {
			replyObject.Embed = &handlerResult.Normal.Embed
			sendReplyMessage(session, md.ChannelID, md.AuthorID, replyObject)
		} else {
			/* Non proccess */
		}

		/* シークレットメッセージの送信 */
		if handlerResult.Secret.EnableType == core.EnContent {
			if handlerResult.Secret.Content != "" {
				replyObject.Content = handlerResult.Secret.Content
				sendReplyMessage(session, core.GetParentIDFromChildID(md.ChannelID), md.AuthorID, replyObject)
			}
		} else if handlerResult.Secret.EnableType == core.EnEmbed {
			replyObject.Embed = &handlerResult.Secret.Embed
			sendReplyMessage(session, core.GetParentIDFromChildID(md.ChannelID), md.AuthorID, replyObject)
		} else {
			/* Non proccess */
		}

	}
}

// sendReplyMessage メッセージ返信処理
func sendReplyMessage(session *discordgo.Session, chID string, to string, message discordgo.MessageSend) {
	if to == "" {
		_, err := session.ChannelMessageSendComplex(chID, &message)
		if err != nil {
			log.Println(err)
		}
	} else {
		message.Content = "<@" + to + "> " + message.Content
		_, err := session.ChannelMessageSendComplex(chID, &message)
		if err != nil {
			log.Println(err)
		}
	}
}

// sendMessage メッセージ送信処理
func sendMessage(session *discordgo.Session, chID string, message discordgo.MessageSend) {
	sendReplyMessage(session, chID, "", message)
}
