package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"Nodens/config"
	"Nodens/core"
	"Nodens/discordDriver"
)

// configFile デフォルト設定ファイルパス
var configFile = "SystemConfig.json"

// main メイン関数
func main() {
	// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
	// Discord非依存処理
	// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
	// 引数の読み込み
	if len(os.Args) != 1 {
		configFile = os.Args[1]
		_, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Panicf("[Error]: Cannot open file '%v': '%v'", configFile, err)
		}
	}

	// 設定ファイルの読み込み
	core.LoadConfig(configFile)

	// テキストコマンドハンドラ登録
	for _, handle := range config.CmdHandleFuncTable {
		core.AddCmdHandler(handle.System, handle.Command, handle.Function)
	}

	// キャラクターデータ取得関数登録
	for _, cdFunc := range config.CharacterDataGetFuncTable {
		core.AddCharacterDataGetFunc(cdFunc.System, cdFunc.DataName, cdFunc.Function)
	}

	// セッション復元関数登録
	core.SetRestoreFunc(config.SessionRestoreFuncTable)

	// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
	// Discord依存処理
	// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
	// Discordのインスタンス生成
	discord, err := discordDriver.JobNewDiscordSession()

	// Discordのメッセージハンドラ登録
	discord.AddHandler(discordDriver.OnMessageCreate)
	discord.AddHandler(discordDriver.OnInteractionCreate)

	// スラッシュコマンドハンドラ登録
	for _, handle := range config.SlashCmdHandleFuncTable {
		core.AddSlashCmdHandler(handle.System, handle.Command, handle.Function)
		discordDriver.AddSlashCmdData(handle.System, handle.SlashCommandData)
	}

	// セッション開始
	err = discord.Open()
	if err != nil {
		log.Panicln(err)
	}

	// 共通コマンドをスラッシュコマンド(グローバル)として登録
	discordDriver.JobRegistriesGlobalAppCommands("General")

	// イベントをリッスン
	log.Println("[Event]: Listening...")
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, os.Interrupt)
	<-stopBot

	// スラッシュコマンド(グローバル)を全て削除
	discordDriver.JobDeleteGlobalAppCommands()
	discordDriver.JobDeleteLocalAppCommands()
	// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
	return
}
