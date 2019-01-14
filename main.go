package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/c-ardinal/Nodens/lib/cthulhu"
	"github.com/c-ardinal/Nodens/lib/nodens"
)

// mainVersion mainパッケージバージョン
const mainVersion string = "0.0.1"

// configFile 設定ファイルパス
var configFile = "SystemConfig.json"

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
	nodens.LoadConfig(configFile)

	// Discordのインスタンス生成
	discord, err := discordgo.New()
	discord.Token = nodens.GetConfig().BotToken
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// ハンドラ登録
	discord.AddHandler(onMessageCreate)
	nodens.AddCmdHandler("ShowVersion", nodens.CmdHandleFunc(cmdShowVersion)) // バージョン情報表示処理
	// TODO: あとで実装 nodens.AddCmdHandler("ShowStatistics", nodens.CmdHandleFunc(cmdShowStatistics))	// ダイスロール統計情報表示処理
	nodens.AddCmdHandler("CreateSession", nodens.CmdHandleFunc(cmdCreateSession))   // 親セッション生成処理
	nodens.AddCmdHandler("ConnectSession", nodens.CmdHandleFunc(cmdConnectSession)) // 親セッション連携処理
	// TODO: あとで実装 nodens.AddCmdHandler("ExitSession", nodens.CmdHandleFunc(cmdExitSession))		// 親セッション連携処理
	// TODO: あとで実装 nodens.AddCmdHandler("RemoveSession", nodens.CmdHandleFunc(cmdRemoveSession))	// セッション削除処理
	// TODO: あとで実装 nodens.AddCmdHandler("StoreSession", nodens.CmdHandleFunc(cmdStoreSession))	// セッション保存処理
	// TODO: あとで実装 nodens.AddCmdHandler("LoadSession", nodens.CmdHandleFunc(cmdLoadSession))		// セッション復帰処理
	nodens.AddCmdHandler("regchara", nodens.CmdHandleFunc(cmdRegistryCharacter)) // キャラシート連携処理
	// TODO: あとで実装 nodens.AddCmdHandler("check", nodens.CmdHandleFunc(cmdCharaNumCheck))        // 能力値確認処理
	nodens.AddCmdHandler("ctrl", nodens.CmdHandleFunc(cmdCharaNumControl)) // 能力値操作処理
	nodens.AddCmdHandler("roll", nodens.CmdHandleFunc(cmdLinkRoll))        // 能力値ダイスロール処理
	nodens.AddCmdHandler("Sroll", nodens.CmdHandleFunc(cmdSecretLinkRoll)) // 能力値シークレットダイスロール処理
	// TODO: あとで実装 nodens.AddCmdHandler("sanc", nodens.CmdHandleFunc(cmdSanRoll))				// SAN値チェック処理
	// TODO: あとで実装 nodens.AddCmdHandler("grow", nodens.CmdHandleFunc(cmdGrowRoll))				// 成長ロール処理

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
	channel, err := session.State.Channel(message.ChannelID)
	authorID := message.Author.ID
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%20s %20s > %s\n", message.ChannelID, message.Author.Username, message.Content)

	if authorID != nodens.GetConfig().BotID {
		exeResult, secResult, err := nodens.ExecuteCmdHandler(channel, message)
		if err != nil {
			log.Println(err)
			return
		}
		if exeResult != "" && exeResult != message.Content {
			sendReplyMessage(session, channel.ID, authorID, exeResult)

			if secResult != "" {
				sendReplyMessage(session, cthulhu.GetParentIDFromChildID(channel.ID), authorID, secResult)
			}
		}
	}
}

// sendReplyMessage メッセージ返信処理
func sendReplyMessage(session *discordgo.Session, chID string, to string, text string) {
	if to == "" {
		_, err := session.ChannelMessageSend(chID, text)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := session.ChannelMessageSend(chID, "<@"+to+"> "+text)
		if err != nil {
			log.Println(err)
		}
	}
}

// sendMessage メッセージ送信処理
func sendMessage(session *discordgo.Session, chID string, text string) {
	sendReplyMessage(session, chID, "", text)
}
