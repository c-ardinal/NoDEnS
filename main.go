package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/c-ardinal/Nodens/core"
	"github.com/c-ardinal/Nodens/system/cthulhu"

	"github.com/bwmarrin/discordgo"
)

// mainVersion mainパッケージバージョン
const mainVersion string = "0.0.1"

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
	{"Cthulhu", "showstat", cthulhu.CmdShowStatistics}, // ダイスロール統計表示処理
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
	channel, err := session.State.Channel(message.ChannelID)
	authorID := message.Author.ID
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%20s %20s > %s\n", message.ChannelID, message.Author.Username, message.Content)

	if authorID != core.GetConfig().BotID {
		exeResult, secResult, err := core.ExecuteCmdHandler(channel, message)
		if err != nil {
			log.Println(err)
			return
		}
		if exeResult != "" && exeResult != message.Content {
			sendReplyMessage(session, channel.ID, authorID, exeResult)

			if secResult != "" {
				sendReplyMessage(session, core.GetParentIDFromChildID(channel.ID), authorID, secResult)
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
