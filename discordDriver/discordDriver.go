package discordDriver

import (
	"Nodens/core"
	"log"

	"github.com/bwmarrin/discordgo"
)

/****************************************************************************/
/* 内部型定義                                                               */
/****************************************************************************/

// cmdHandleFuncStruct コマンドハンドラテーブル用構造体
type cmdHandleFuncStruct struct {
	System           string
	SlashCommandData discordgo.ApplicationCommand
}

/****************************************************************************/
/* 内部定数定義                                                             */
/****************************************************************************/

/****************************************************************************/
/* 内部変数定義                                                             */
/****************************************************************************/

// Discordセッションインスタンス
var session *discordgo.Session

// 登録済みスラッシュコマンドリスト
var registeredGuildIds []string

// スラッシュコマンド定義リスト
var slashCmdHandleFuncList []cmdHandleFuncStruct

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/

// Discordのセッションを生成する
func JobNewDiscordSession() (ses *discordgo.Session, err error) {
	ses, err = discordgo.New(core.GetConfig().BotToken)
	if err != nil {
		log.Panicf("[Error]: Cannot create discord instance : '%v'", err)
	}
	session = ses
	return ses, err
}

// Discordのセッションを返す
func GetDiscordSession() *discordgo.Session {
	return session
}

// スラッシュコマンドの情報を設定する
func AddSlashCmdData(system string, appCommand discordgo.ApplicationCommand) {
	newData := cmdHandleFuncStruct{
		System:           system,
		SlashCommandData: appCommand,
	}
	slashCmdHandleFuncList = append(slashCmdHandleFuncList, newData)
}

// メッセージ受信時処理
func OnMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	var md core.MessageData = core.MessageData{
		ChannelID:     message.ChannelID,
		GuildID:       message.GuildID,
		MessageID:     message.ID,
		AuthorID:      message.Author.ID,
		AuthorName:    message.Author.Username,
		MessageString: message.Content,
	}

	log.Printf("[Event]: Message received. ChannelId:%20s Author:%20s > %s\n", md.ChannelID, md.AuthorName, md.MessageString)

	if md.AuthorID != core.GetConfig().BotID {
		handlerResult := core.ExecuteCmdHandler(md)
		if handlerResult.Error != nil {
			log.Printf("[Error]: %v", handlerResult.Error)
		}

		ref := discordgo.MessageReference{
			MessageID: md.MessageID,
			ChannelID: md.ChannelID,
		}

		// キャラクター名取得
		characterName := core.GetCharacterData(md.ChannelID, md.AuthorID, "CharacterName")
		cSheetUrl := core.GetCharacterData(md.ChannelID, md.AuthorID, "CSheetUrl")
		if characterName != "" {
			characterName = "【" + characterName + "】 "
		}

		/* 通常メッセージの送信 */
		if handlerResult.Normal.EnableType == core.EnContent {
			if handlerResult.Normal.Content != "" && handlerResult.Normal.Content != md.MessageString {
				handlerResult.Normal.Content = characterName + handlerResult.Normal.Content
				_, err := session.ChannelMessageSendReply(md.ChannelID, handlerResult.Normal.Content, &ref)
				if err != nil {
					log.Printf("[Warning]: Send failed. %v", err.Error())
				}
			}
		} else if handlerResult.Normal.EnableType == core.EnEmbed {
			embedAuthor := &discordgo.MessageEmbedAuthor{
				Name: characterName,
				URL:  cSheetUrl,
			}
			handlerResult.Normal.Embed.Author = embedAuthor
			_, err := session.ChannelMessageSendEmbedReply(md.ChannelID, handlerResult.Normal.Embed, &ref)
			if err != nil {
				log.Printf("[Warning]: Send failed. %v", err.Error())
			}
		}

		/* シークレットメッセージの送信 */
		if handlerResult.Error == nil {
			if handlerResult.Secret.EnableType == core.EnContent {
				if handlerResult.Secret.Content != "" {
					handlerResult.Secret.Content = "<@" + md.AuthorID + ">" + handlerResult.Secret.Content
					_, err := session.ChannelMessageSend(core.GetParentIDFromChildID(md.ChannelID), handlerResult.Secret.Content)
					if err != nil {
						log.Printf("[Warning]: Send failed. %v", err.Error())
					}
				}
			} else if handlerResult.Secret.EnableType == core.EnEmbed {
				messageSend := &discordgo.MessageSend{
					Content: "<@" + md.AuthorID + ">",
					Embed:   handlerResult.Secret.Embed,
				}
				_, err := session.ChannelMessageSendComplex(core.GetParentIDFromChildID(md.ChannelID), messageSend)
				if err != nil {
					log.Printf("[Warning]: Send failed. %v", err.Error())
				}
			}
		}
	}
}

// インタラクション受信時処理
func OnInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Message == nil {
		jobInteractionMessage(session, interaction)
	} else {
		jobInteractionButton(session, interaction)
	}
}

// ボタン処理
func jobInteractionButton(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.MessageComponentData().CustomID {
	case "is-secret-open":
		{
			// キャラクター名取得
			characterName := core.GetCharacterData(interaction.ChannelID, interaction.Member.User.ID, "CharacterName")
			cSheetUrl := core.GetCharacterData(interaction.ChannelID, interaction.Member.User.ID, "CSheetUrl")
			if characterName != "" {
				characterName = "【" + characterName + "】 "
			}
			embedAuthor := &discordgo.MessageEmbedAuthor{
				Name: characterName,
				URL:  cSheetUrl,
			}
			// 応答
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "<@" + interaction.Member.User.ID + ">",
					Embeds: []*discordgo.MessageEmbed{
						{
							Author:      embedAuthor,
							Title:       "シークレットロール結果公開",
							Description: interaction.Message.Embeds[0].Title + "\n" + interaction.Message.Embeds[0].Description,
							Color:       core.EnColorGreen,
						},
					},
				},
			})
			if err != nil {
				log.Printf("[Warning]: Send failed. %v", err.Error())
			}
		}
	default:
		/* Non process */
	}
}

// スラッシュコマンド処理
func jobInteractionMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// インタラクション送信者情報取得
	var interactionUser *discordgo.User
	if interaction.User != nil {
		interactionUser = interaction.User
	} else if interaction.Member != nil {
		interactionUser = interaction.Member.User
	}

	// スラッシュコマンド情報取得
	var options []core.CommandOption
	appCommandData := interaction.Interaction.ApplicationCommandData()
	for _, opt := range appCommandData.Options {
		options = append(options, core.CommandOption{
			Name:  opt.Name,
			Value: opt.StringValue(),
		})
	}

	// 受信メッセージデータを構築
	var md core.MessageData = core.MessageData{
		ChannelID:     interaction.ChannelID,
		GuildID:       interaction.GuildID,
		MessageID:     interaction.ID,
		AuthorID:      interactionUser.ID,
		AuthorName:    interactionUser.Username,
		MessageString: "",
		Command:       appCommandData.Name,
		Options:       options,
	}

	// コマンド処理実行
	handlerResult := ExecuteSlashCmdHandler(md)
	if handlerResult.Error != nil {
		log.Printf("[Error]: %v", handlerResult.Error)
	}

	// キャラクター名取得
	characterName := core.GetCharacterData(md.ChannelID, md.AuthorID, "CharacterName")
	cSheetUrl := core.GetCharacterData(md.ChannelID, md.AuthorID, "CSheetUrl")
	if characterName != "" {
		characterName = "【" + characterName + "】 "
	}

	// シークレットメッセージを含む場合、Ephemeralフラグを立て、シークレットメッセージ公開用ボタンを設定する。
	var flags discordgo.MessageFlags = 0
	var components = []discordgo.MessageComponent{}
	if handlerResult.Secret.EnableType == core.EnContent || handlerResult.Secret.EnableType == core.EnEmbed {
		flags = discordgo.MessageFlagsEphemeral
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "is-secret-open",
						Label:    "結果を公開する",
						Style:    discordgo.PrimaryButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "👀",
						},
					},
				},
			},
		}

	}

	/* 通常メッセージの送信 */
	if handlerResult.Normal.EnableType == core.EnContent {
		if handlerResult.Normal.Content != "" && handlerResult.Normal.Content != md.MessageString {
			handlerResult.Normal.Content = characterName + handlerResult.Normal.Content
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:      flags,
					Content:    handlerResult.Normal.Content,
					Components: components,
				},
			})
			if err != nil {
				log.Printf("[Warning]: Send failed. %v", err.Error())
			}
		}
	} else if handlerResult.Normal.EnableType == core.EnEmbed {
		embedAuthor := &discordgo.MessageEmbedAuthor{
			Name: characterName,
			URL:  cSheetUrl,
		}
		handlerResult.Normal.Embed.Author = embedAuthor
		err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:      flags,
				Embeds:     []*discordgo.MessageEmbed{handlerResult.Normal.Embed},
				Components: components,
			},
		})
		if err != nil {
			log.Printf("[Warning]: Send failed. %v", err.Error())
		}
	}
	/* シークレットメッセージの送信 */
	if handlerResult.Error == nil {
		if handlerResult.Secret.EnableType == core.EnContent {
			if handlerResult.Secret.Content != "" {
				handlerResult.Secret.Content = "<@" + md.AuthorID + ">" + handlerResult.Secret.Content
				_, err := session.ChannelMessageSend(md.ChannelID, handlerResult.Secret.Content)
				if err != nil {
					log.Printf("[Warning]: Send failed. %v", err.Error())
				}
			}
		} else if handlerResult.Secret.EnableType == core.EnEmbed {
			messageSend := &discordgo.MessageSend{
				Content: "<@" + md.AuthorID + ">",
				Embed:   handlerResult.Secret.Embed,
			}
			_, err := session.ChannelMessageSendComplex(md.ChannelID, messageSend)
			if err != nil {
				log.Printf("[Warning]: Send failed. %v", err.Error())
			}
		}
	}
}

// スラッシュコマンド(グローバル)登録
func JobRegistriesGlobalAppCommands(targetSystem string) {
	JobRegistriesAppCommands(targetSystem, "")
}

// スラッシュコマンド登録
func JobRegistriesAppCommands(targetSystem string, guildId string) {
	for _, slashCmd := range slashCmdHandleFuncList {
		if slashCmd.System == targetSystem {
			cmd, err := session.ApplicationCommandCreate(session.State.User.ID, guildId, &slashCmd.SlashCommandData)
			if err != nil {
				log.Panicf("[Error]: Cannot register '%v' command: %v", cmd, err)
			}
			log.Printf("[Event]: Command registered '%v'", cmd.Name)
		}
	}
}

// スラッシュコマンド(グローバル)全削除
func JobDeleteGlobalAppCommands() {
	JobDeleteAppCommands("")
}

// スラッシュコマンド(ローカル)全削除
func JobDeleteLocalAppCommands() {
	for _, guildId := range registeredGuildIds {
		JobDeleteAppCommands(guildId)
	}
}

// スラッシュコマンド削除
func JobDeleteAppCommands(guildId string) {
	registeredCommandList, err := session.ApplicationCommands(session.State.User.ID, guildId)
	if err == nil {
		for _, appCommand := range registeredCommandList {
			err := session.ApplicationCommandDelete(session.State.User.ID, guildId, appCommand.ID)
			if err != nil {
				log.Panicf("[Error]: Cannot delete '%v' command: %v", appCommand.Name, err)
			}
			log.Printf("[Event]: Command deleted '%v'", appCommand.Name)
		}
	}
}
