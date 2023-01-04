package config

import (
	"Nodens/core"
	"Nodens/system/cthulhu"
)

// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
// コマンドハンドラコンフィグ
// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
// CmdHandleFuncStruct コマンドハンドラテーブル用構造体
type CmdHandleFuncStruct struct {
	System   string
	Command  string
	Function core.CmdHandleFunc
}

// CmdHandleFuncTable コマンドハンドラテーブル
var CmdHandleFuncTable = []CmdHandleFuncStruct{
	{"General", "ShowVersion", core.CmdShowVersion},       // バージョン情報表示処理
	{"General", "CreateSession", core.CmdCreateSession},   // 親セッション生成処理
	{"General", "ConnectSession", core.CmdConnectSession}, // 親セッション連携処理
	// TODO: あとで実装 {"General", "ExitSession", core.cmdExitSession},		// 親セッション連携解除処理
	// TODO: あとで実装 {"General", "RemoveSession", core.cmdRemoveSession},	// 親セッション削除処理
	{"General", "StoreSession", core.CmdStoreSession},     // セッション保存処理
	{"General", "RestoreSession", core.CmdRestoreSession}, // セッション復帰処理
	{"Cthulhu", "regchara", cthulhu.CmdRegistryCharacter}, // キャラシート連携処理
	{"Cthulhu", "check", cthulhu.CmdCharaNumCheck},        // 能力値確認処理
	{"Cthulhu", "ctrl", cthulhu.CmdCharaNumControl},       // 能力値操作処理
	{"Cthulhu", "roll", cthulhu.CmdLinkRoll},              // 能力値ダイスロール処理
	{"Cthulhu", "Sroll", cthulhu.CmdSecretLinkRoll},       // 能力値シークレットダイスロール処理
	{"Cthulhu", "sanc", cthulhu.CmdSanCheckRoll},          // SAN値チェック処理
	// TODO: あとで実装 {"Cthulhu", "grow", cmdGrowRoll},        // 成長ロール処理
	// TODO: あとで実装 {"Cthulhu", "resist", cmdResistRoll},    // 対抗ロール処理
	// TODO: 実装中 {"Cthulhu", "showstat", cthulhu.CmdShowStatistics}, // ダイスロール統計表示処理
}

// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
// core→各システムキャラクターデータ取得関数コンフィグ
// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
// CharacterDataGetFuncStruct キャラクターデータ取得関数用構造体
type CharacterDataGetFuncStruct struct {
	System   string
	DataName string
	Function core.CharacterDataGetFunc
}

// CharacterDataGetFuncTable キャラクターデータ取得関数テーブル
var CharacterDataGetFuncTable = []CharacterDataGetFuncStruct{
	{"Cthulhu", "CharacterName", cthulhu.GetCharacterName},
	{"Cthulhu", "CSheetUrl", cthulhu.GetCharacterSheetUrl},
}

// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
// core→各システムセッション復元関数コンフィグ
// ■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□■□
var SessionRestoreFuncTable = map[string]core.SessionRestoreFunc{
	"Cthulhu": cthulhu.JobRestoreSession,
}
