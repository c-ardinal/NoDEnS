package main

import (
	"io/ioutil"
	"log"
	"os"

	"Nodens/config"
	"Nodens/core"

	"github.com/bwmarrin/discordgo"
)

// configFile デフォルト設定ファイルパス
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
	core.LoadConfig(configFile)

	// Discordのインスタンス生成
	discord, err := discordgo.New(core.GetConfig().BotToken)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// Discordのメッセージハンドラ登録
	discord.AddHandler(onMessageCreate)

	// コマンドハンドラ登録
	for _, handle := range config.CmdHandleFuncTable {
		core.AddCmdHandler(handle.System, handle.Command, handle.Function)
	}

	// キャラデータ取得関数登録
	for _, cdFunc := range config.CharacterDataGetFuncTable {
		core.AddCharacterDataGetFunc(cdFunc.System, cdFunc.DataName, cdFunc.Function)
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
			/* Non process */
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
			/* Non process */
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
		characterName := core.GetCharacterName(chID, to)
		if characterName != "" {
			characterName = "【" + characterName + "】"
		}
		message.Content = "<@" + to + "> " + characterName + " " + message.Content
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
