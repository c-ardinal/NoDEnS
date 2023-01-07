package config

import (
	"Nodens/core"
	"Nodens/discordDriver"
	"Nodens/system/cthulhu"

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
/* 外部公開型定義                                                           */
/****************************************************************************/

// コマンドハンドラテーブル用構造体
type CmdHandleFuncStruct struct {
	System           string
	Command          string
	Function         core.CmdHandleFunc
	SlashCommandData discordgo.ApplicationCommand
}

// キャラクターデータ取得関数用構造体
type CharacterDataGetFuncStruct struct {
	System   string
	DataName string
	Function core.CharacterDataGetFunc
}

/****************************************************************************/
/* 外部公開定数定義                                                         */
/****************************************************************************/

// 共通コマンド定数定義
const STR_CMD_VERSION string = "version"
const STR_CMD_CREATE_SESSION string = "create-session"
const STR_CMD_CONNECT_SESSION string = "connect-session"
const STR_CMD_STORE_SESSION string = "store-session"
const STR_CMD_RESTORE_SESSION string = "restore-session"

// スラッシュコマンド設定用定数定義
var BOL_DAT_DM_PERMISSION_ALLOW bool = true                                              // DMPermissionにはconstを設定出来ないためvarで定義
var BOL_DAT_DM_PERMISSION_DENY bool = false                                              // DMPermissionにはconstを設定出来ないためvarで定義
var INT_DAT_MEMBER_PERMISSION_MANAGE_CHANNELS int64 = discordgo.PermissionManageChannels // DefaultMemberPermissionにはconstを設定出来ないためvarで定義

// スラッシュコマンドハンドラテーブル
var SlashCmdHandleFuncTable = []CmdHandleFuncStruct{
	{"General", STR_CMD_VERSION, core.CmdShowVersion, // バージョン情報表示処理
		discordgo.ApplicationCommand{
			Name:         STR_CMD_VERSION,
			Description:  "ダイスボットのバージョンを表示します。",
			DMPermission: &BOL_DAT_DM_PERMISSION_ALLOW,
		},
	},
	{"General", STR_CMD_CREATE_SESSION, discordDriver.CmdCreateSession, // 親セッション生成処理
		discordgo.ApplicationCommand{
			Name:         STR_CMD_CREATE_SESSION,
			Description:  "TRPGセッションを生成し、ダイスボットを有効化します。",
			DMPermission: &BOL_DAT_DM_PERMISSION_DENY,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "system",
					Description: "有効化するダイスシステムを指定します。",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "クトゥルフ神話TRPG",
							Value: "Cthulhu",
						},
						{
							Name:  "ただのダイスボット",
							Value: "DiceBot",
						},
						{
							Name:  "その他(オプション:other-systemを指定してください)",
							Value: "OtherSystem",
						},
					},
				},
				{
					Name:        "forced",
					Description: "既にTRPGセッションが生成されている場合に、セッションを再生成するか否かを指定します。",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "再生成する",
							Value: "--forced ",
						},
						{
							Name:  "再生成しない",
							Value: "",
						},
					},
				},
				{
					Name:        "other-system",
					Description: "オプション:system一覧に無いシステムを使用した場合に指定して下さい。",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	},
	{"General", STR_CMD_CONNECT_SESSION, core.CmdConnectSession, // 親セッション連携処理
		discordgo.ApplicationCommand{
			Name:         STR_CMD_CONNECT_SESSION,
			Description:  "生成済のTRPGセッションに接続します(シークレットダイスが振れるようになります)。",
			DMPermission: &BOL_DAT_DM_PERMISSION_ALLOW,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "channel-id",
					Description: "接続先のDiscordチャネルID(create-session実行時に表示されるID)を指定します。",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	},
	// TODO: あとで実装 {"General", "ExitSession", core.cmdExitSession},		// 親セッション連携解除処理
	// TODO: あとで実装 {"General", "RemoveSession", core.cmdRemoveSession},	// 親セッション削除処理
	{"General", STR_CMD_STORE_SESSION, core.CmdStoreSession, // セッション保存処理
		discordgo.ApplicationCommand{
			Name:                     STR_CMD_STORE_SESSION,
			Description:              "TRPGセッションを保存します。",
			DefaultMemberPermissions: &INT_DAT_MEMBER_PERMISSION_MANAGE_CHANNELS,
			DMPermission:             &BOL_DAT_DM_PERMISSION_DENY,
		},
	},
	{"General", STR_CMD_RESTORE_SESSION, discordDriver.CmdRestoreSession, // セッション復帰処理
		discordgo.ApplicationCommand{
			Name:                     STR_CMD_RESTORE_SESSION,
			Description:              "TRPGセッションを復元します。",
			DefaultMemberPermissions: &INT_DAT_MEMBER_PERMISSION_MANAGE_CHANNELS,
			DMPermission:             &BOL_DAT_DM_PERMISSION_DENY,
		},
	},
	{"Cthulhu", "sec-dice", cthulhu.CmdSecretDiceRoll, // シークレットダイスロール
		discordgo.ApplicationCommand{
			Name:         "sec-dice",
			Description:  "シークレットダイスロールを実施します。",
			DMPermission: &BOL_DAT_DM_PERMISSION_DENY,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "command",
					Description: "技能を指定します。",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	},
	{"Cthulhu", "sec-skill", cthulhu.CmdSecretLinkRoll, // シークレット技能ロール
		discordgo.ApplicationCommand{
			Name:         "sec-skill",
			Description:  "シークレット技能ロールを実施します。",
			DMPermission: &BOL_DAT_DM_PERMISSION_DENY,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "command",
					Description: "ダイスを指定します。",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	},
}

// テキストコマンドハンドラテーブル
var CmdHandleFuncTable = []CmdHandleFuncStruct{
	{"Cthulhu", "regchara", cthulhu.CmdRegistryCharacter, // キャラクターシート連携処理
		discordgo.ApplicationCommand{},
	},
	{"Cthulhu", "check", cthulhu.CmdCharaNumCheck, // 能力値確認処理
		discordgo.ApplicationCommand{},
	},
	{"Cthulhu", "ctrl", cthulhu.CmdCharaNumControl, // 能力値操作処理
		discordgo.ApplicationCommand{},
	},
	{"Cthulhu", "roll", cthulhu.CmdLinkRoll, // 能力値ダイスロール処理
		discordgo.ApplicationCommand{},
	},
	{"Cthulhu", "Sroll", cthulhu.CmdSecretLinkRoll, // 能力値シークレットダイスロール処理
		discordgo.ApplicationCommand{},
	},
	{"Cthulhu", "sanc", cthulhu.CmdSanCheckRoll, // SAN値チェック処理
		discordgo.ApplicationCommand{},
	},
	// TODO: あとで実装 {"Cthulhu", "grow", cmdGrowRoll},        // 成長ロール処理
	// TODO: あとで実装 {"Cthulhu", "resist", cmdResistRoll},    // 対抗ロール処理
	// TODO: 実装中 {"Cthulhu", "showstat", cthulhu.CmdShowStatistics}, // ダイスロール統計表示処理
}

// キャラクターデータ取得関数テーブル
var CharacterDataGetFuncTable = []CharacterDataGetFuncStruct{
	{"Cthulhu", "CharacterName", cthulhu.GetCharacterName},
	{"Cthulhu", "CSheetUrl", cthulhu.GetCharacterSheetUrl},
}

// core→各システムセッション復元関数コンフィグ
var SessionRestoreFuncTable = map[string]core.SessionRestoreFunc{
	"Cthulhu": cthulhu.JobRestoreSession,
}

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/
