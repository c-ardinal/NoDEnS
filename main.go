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

	// セッション復元関数登録
	core.SetRestoreFunc(config.SessionRestoreFuncTable)

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
	var md core.MessageData = core.MessageData{
		ChannelID:     message.ChannelID,
		MessageID:     message.ID,
		AuthorID:      message.Author.ID,
		AuthorName:    message.Author.Username,
		MessageString: message.Content,
	}

	//var replyObject discordgo.MessageSend

	log.Printf("%20s %20s > %s\n", md.ChannelID, md.AuthorName, md.MessageString)

	if md.AuthorID != core.GetConfig().BotID {
		handlerResult := core.ExecuteCmdHandler(md)
		if handlerResult.Error != nil {
			log.Println(handlerResult.Error)
		}

		ref := discordgo.MessageReference{
			MessageID: md.MessageID,
			ChannelID: md.ChannelID,
		}

		characterName := core.GetCharacterData(md.ChannelID, md.AuthorID, "CharacterName")
		cSheetUrl := core.GetCharacterData(md.ChannelID, md.AuthorID, "CSheetUrl")
		if characterName != "" {
			characterName = "【" + characterName + "】 "
		}

		/* 通常メッセージの送信 */
		if handlerResult.Normal.EnableType == core.EnContent {
			if handlerResult.Normal.Content != "" && handlerResult.Normal.Content != md.MessageString {
				handlerResult.Normal.Content = characterName + handlerResult.Normal.Content
				session.ChannelMessageSendReply(md.ChannelID, handlerResult.Normal.Content, &ref)
			}
		} else if handlerResult.Normal.EnableType == core.EnEmbed {
			embedAuthor := &discordgo.MessageEmbedAuthor{
				Name: characterName,
				URL:  cSheetUrl,
			}
			handlerResult.Normal.Embed.Author = embedAuthor
			session.ChannelMessageSendEmbedReply(md.ChannelID, handlerResult.Normal.Embed, &ref)
		} else {
			/* Non process */
		}

		/* シークレットメッセージの送信 */
		if handlerResult.Secret.EnableType == core.EnContent {
			if handlerResult.Secret.Content != "" {
				handlerResult.Secret.Content = "<@" + md.AuthorID + ">" + handlerResult.Secret.Content
				session.ChannelMessageSend(core.GetParentIDFromChildID(md.ChannelID), handlerResult.Secret.Content)
			}
		} else if handlerResult.Secret.EnableType == core.EnEmbed {
			messageSend := &discordgo.MessageSend{
				Content: "<@" + md.AuthorID + ">",
				Embed:   handlerResult.Secret.Embed,
			}
			session.ChannelMessageSendComplex(core.GetParentIDFromChildID(md.ChannelID), messageSend)
		} else {
			/* Non process */
		}
	}
}
