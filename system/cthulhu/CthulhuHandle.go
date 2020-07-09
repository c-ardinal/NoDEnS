package cthulhu

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/c-ardinal/Nodens/core"
)

// DiceResultLogOfCthulhu ダイスロール実行ログ型
type DiceResultLogOfCthulhu struct {
	Player  core.NaID
	Time    string
	Command string
	Result  string
}

// DiceStatsticsOfCthulhu ダイスロール統計型
type DiceStatsticsOfCthulhu struct {
	Player   core.NaID
	Critical []string
	Special  []string
	Success  []string
	Fail     []string
	Fumble   []string
}

// DiceResultLogOfCthulhus ダイスロール実行ログ格納変数
var DiceResultLogOfCthulhus = []DiceResultLogOfCthulhu{}

// CmdRestoreSessionOfCthulhu クトゥルフのセッションを復元する
func CmdRestoreSessionOfCthulhu(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	err := core.RestoreSession(md.ChannelID)
	if err != nil {
		handlerResult.Normal.Content = "Session load failed."
		handlerResult.Error = err
	} else {
		ses := core.GetSessionByID(md.ChannelID)

		// PC情報を一度JSONに戻してからクトゥルフ用PC構造体に変換する
		pcsRawData, _ := json.Marshal((*ses).Pc)
		var pcsMap = map[string]*CharacterOfCthulhu{}
		json.Unmarshal(pcsRawData, &pcsMap)

		// NPC情報を一度JSONに戻してからクトゥルフ用NPC構造体に変換する
		npcsRawData, _ := json.Marshal((*ses).Npc)
		var npcsMap = map[string]*CharacterOfCthulhu{}
		json.Unmarshal(npcsRawData, &npcsMap)

		// PC情報を格納
		for _, pcData := range pcsMap {
			(*ses).Pc[pcData.Player.ID] = pcData
		}

		// NPC情報を格納
		for _, npcData := range npcsMap {
			(*ses).Npc[npcData.Player.ID] = npcData
		}

		handlerResult.Normal.Content = "Session store successfully."
	}

	return handlerResult
}

// CmdRegistryCharacter キャラシ連携ハンドラ
func CmdRegistryCharacter(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var cd *CharacterOfCthulhu
	if len(opt) == 0 {
		handlerResult.Normal.Content = "Invalid arguments."
		return handlerResult
	}
	if core.CheckExistSession(md.ChannelID) == true {
		if core.CheckExistCharacter(md.ChannelID, md.AuthorID) == true {
			handlerResult.Normal.Content = "Character already exists."
			return handlerResult
		}
		cas, err := GetCharSheetFromURL(opt[0])
		if err != nil {
			handlerResult.Normal.Content = "Registry failed."
			handlerResult.Error = err
			return handlerResult
		}
		cd = GetCharDataFromCharSheet(cas, md.AuthorName, md.AuthorID)
		(*cs).Pc[md.AuthorID] = cd
	} else if core.GetParentIDFromChildID(md.ChannelID) != "" {
		if core.CheckExistNPCharacter(core.GetParentIDFromChildID(md.ChannelID), md.AuthorID) == true {
			handlerResult.Normal.Content = "Character already exists."
			return handlerResult
		}
		cas, err := GetCharSheetFromURL(opt[0])
		if err != nil {
			handlerResult.Normal.Content = "Registry failed."
			handlerResult.Error = err
			return handlerResult
		}
		cd = GetCharDataFromCharSheet(cas, md.AuthorName, md.AuthorID)
		(*cs).Npc[md.AuthorID] = cd
	} else {
		handlerResult.Normal.Content = "Session not found."
		return handlerResult
	}

	handlerResult.Normal.Content = "\r\n====================\r\n"
	handlerResult.Normal.Content += "**[名 前]** " + cd.Personal.Name + "\r\n"
	handlerResult.Normal.Content += "**[年 齢]** " + strconv.Itoa(cd.Personal.Age) + "歳\r\n"
	handlerResult.Normal.Content += "**[性 別]** " + cd.Personal.Sex + "\r\n"
	handlerResult.Normal.Content += "**[職 業]** " + cd.Personal.Job + "\r\n"
	for _, cdan := range CdAbilityNameList {
		a := cd.Ability[cdan]
		if a.Now == a.Init {
			handlerResult.Normal.Content += "**[ " + a.Name + " ]** " + strconv.Itoa(a.Now) + "\r\n"
		} else {
			handlerResult.Normal.Content += "**[ " + a.Name + " ]** " + strconv.Itoa(a.Now) + " (Init: " + strconv.Itoa(a.Init) + ")\r\n"
		}
	}
	handlerResult.Normal.Content += "**[メ モ]** \r\n" + cd.Memo + "\r\n"
	handlerResult.Normal.Content += "====================\r\n"

	return handlerResult
}

// CmdCharaNumCheck 能力値確認ハンドラ
func CmdCharaNumCheck(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var chara *CharacterOfCthulhu
	var exist bool
	if len(opt) == 0 {
		handlerResult.Normal.Content = "Invalid arguments."
		return handlerResult
	}
	if cs == nil {
		handlerResult.Normal.Content = "Character not registried."
		return handlerResult
	}
	if core.GetParentIDFromChildID(md.ChannelID) != "" {
		chara, exist = (*cs).Npc[md.AuthorID].(*CharacterOfCthulhu)
	} else {
		chara, exist = (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
	}
	if exist == false {
		handlerResult.Normal.Content = "Character not found."
		return handlerResult
	}
	initNum := GetSkillNum(chara, opt[0], "init")
	if initNum == "-1" {
		handlerResult.Normal.Content = "Skill not found."
		return handlerResult
	}
	startNum := GetSkillNum(chara, opt[0], "sum")
	nowNum := GetSkillNum(chara, opt[0], "now")

	handlerResult.Normal.Content = "[" + opt[0] + "] Init( " + initNum + " ), Start( " + startNum + "), Now( " + nowNum + " )"

	return handlerResult
}

// CmdCharaNumControl 能力値操作ハンドラ
func CmdCharaNumControl(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var chara *CharacterOfCthulhu
	var exist bool
	if len(opt) < 2 {
		handlerResult.Normal.Content = "Invalid arguments."
		return handlerResult
	}
	if cs == nil {
		handlerResult.Normal.Content = "Character not registried."
		return handlerResult
	}
	if core.GetParentIDFromChildID(md.ChannelID) != "" {
		chara, exist = (*cs).Npc[md.AuthorID].(*CharacterOfCthulhu)
	} else {
		chara, exist = (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
	}
	if exist == false {
		handlerResult.Normal.Content = "Character not found."
		return handlerResult
	}
	oldNum := GetSkillNum(chara, opt[0], "now")
	if oldNum == "-1" {
		handlerResult.Normal.Content = "Skill not found."
		return handlerResult
	}
	diffRegex := regexp.MustCompile("^[+-]?[0-9]+$")
	var diffCmd string = opt[1]
	if diffRegex.MatchString(diffCmd) == false {
		minusFlag := false
		if strings.Contains(diffCmd, "-") {
			diffCmd = strings.ReplaceAll(diffCmd, "-", "")
			minusFlag = true
		}
		rollResult, err := core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, diffCmd)
		if err != nil {
			handlerResult.Normal.Content = "Invalid diff num."
			handlerResult.Error = err
			return handlerResult
		}
		var sum int
		for _, r := range rollResult.Dices {
			sum += r.Value
		}

		if minusFlag {
			diffCmd = "-" + strconv.Itoa(sum)
		} else {
			diffCmd = strconv.Itoa(sum)
		}

	}
	newNum := AddSkillNum(chara, opt[0], diffCmd)

	handlerResult.Normal.Content = "[" + opt[0] + "] " + oldNum + " => " + newNum + " (Diff: " + diffCmd + ")"

	return handlerResult
}

// CmdLinkRoll キャラシ連携ダイスロールハンドラ
func CmdLinkRoll(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	if len(opt) == 0 {
		handlerResult.Normal.Content = "Invalid arguments."
		return handlerResult
	}
	if cs == nil {
		handlerResult.Normal.Content = "PC not registried."
		return handlerResult
	}
	pc, exist := (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
	if exist == false {
		handlerResult.Normal.Content = "PC not found."
		return handlerResult
	}
	diceCmd := "CCB<=" + opt[0]
	exRegex := regexp.MustCompile("[^\\+\\-\\*\\/ 　]+")
	ignoreRegex := regexp.MustCompile("^[0-9]+$")
	for _, ex := range exRegex.FindAllString(opt[0], -1) {
		if ignoreRegex.MatchString(ex) == false {
			exNum := GetSkillNum(pc, ex, "now")
			if exNum == "-1" {
				handlerResult.Normal.Content = "Skill not found."
				return handlerResult
			}
			diceCmd = strings.Replace(diceCmd, ex, exNum, -1)
		}
	}
	rollResult, err := core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, diceCmd)
	handlerResult.Normal.Content = rollResult.Result

	if err == nil {
		handlerResult.Error = err

		//const format = "2006/01/02_15:04:05"
		//parsedTime, _ := mes.Timestamp.Parse()
		var cthulhuDiceResultLog DiceResultLogOfCthulhu

		cthulhuDiceResultLog.Player.ID = md.AuthorID
		cthulhuDiceResultLog.Player.Name = md.AuthorName
		//cthulhuDiceResultLog.Time = parsedTime.Format(format)
		cthulhuDiceResultLog.Command = md.MessageString
		cthulhuDiceResultLog.Result = rollResult.Result
		DiceResultLogOfCthulhus = append(DiceResultLogOfCthulhus, cthulhuDiceResultLog)
	}

	return handlerResult
}

// CmdSecretLinkRoll キャラシ連携Secretダイスロールハンドラ
func CmdSecretLinkRoll(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	if len(opt) == 0 {
		handlerResult.Normal.Content = "Invalid arguments."
		return handlerResult
	}
	if cs == nil {
		handlerResult.Normal.Content = "NPC not registried."
		return handlerResult
	}
	pc, exist := (*cs).Npc[md.AuthorID].(*CharacterOfCthulhu)
	if exist == false {
		handlerResult.Normal.Content = "NPC not found."
		return handlerResult
	}
	diceCmd := "SCCB<=" + opt[0]
	exRegex := regexp.MustCompile("[^\\+\\-\\*\\/ 　]+")
	ignoreRegex := regexp.MustCompile("^[0-9]+$")
	for _, ex := range exRegex.FindAllString(opt[0], -1) {
		if ignoreRegex.MatchString(ex) == false {
			exNum := GetSkillNum(pc, ex, "now")
			if exNum == "-1" {
				handlerResult.Normal.Content = "Skill not found."
				return handlerResult
			}
			diceCmd = strings.Replace(diceCmd, ex, exNum, -1)
		}
	}
	rollResult, err := core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, diceCmd)
	handlerResult.Normal.Content = rollResult.Result
	if rollResult.Secret == true {
		handlerResult.Secret.Content = "**SECRET DICE**"
	}
	handlerResult.Error = err

	return handlerResult
}

// CmdSanCheckRoll SAN値チェック処理ハンドラ
func CmdSanCheckRoll(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var sucSub string
	var failSub string

	if len(opt) < 2 {
		handlerResult.Normal.Content = "Invalid arguments."
		return handlerResult
	}
	if cs == nil {
		handlerResult.Normal.Content = "PC not registried."
		return handlerResult
	}
	pc, exist := (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
	if exist == false {
		handlerResult.Normal.Content = "PC not found."
		return handlerResult
	}

	orgSanNum := GetSkillNum(pc, "san", "now")
	sanRollCmd := "SCCB<=" + orgSanNum
	sanRollResult, err := core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, sanRollCmd)

	if err != nil {
		handlerResult.Normal.Content = "Server error."
		handlerResult.Error = err
	} else {
		if strings.Contains(sanRollResult.Result, "成功") || strings.Contains(sanRollResult.Result, "スペシャル") {
			if strings.Contains(opt[0], "d") {
				sucRollResult, _ := core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, opt[0])
				sucSub = "-" + core.CalcDicesSum(sucRollResult.Dices)
			} else {
				sucSub = "-" + opt[0]
			}
			newNum := AddSkillNum(pc, "san", sucSub)
			handlerResult.Normal.Content = "sanc > [ " + sanRollResult.Result + " ] >> SAN: " + orgSanNum + " -> " + newNum + " ( " + sucSub + " )"
		} else {
			if strings.Contains(opt[1], "d") {
				failRollResult, _ := core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, opt[1])
				failSub = "-" + core.CalcDicesSum(failRollResult.Dices)
			} else {
				failSub = "-" + opt[1]
			}
			newNum := AddSkillNum(pc, "san", failSub)
			handlerResult.Normal.Content = "sanc >> [ " + sanRollResult.Result + " ] >> SAN: " + orgSanNum + " -> " + newNum + " ( " + failSub + " )"
		}
	}

	return handlerResult
}

// CmdShowStatistics ダイスロール統計表示処理
func CmdShowStatistics(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var diceResultLogs = core.GetDiceResultLogs()
	var diceResultStatstics = map[string]DiceStatsticsOfCthulhu{}

	// 共通ダイスの集計
	for _, drl := range diceResultLogs {
		drs, isExist := diceResultStatstics[drl.Player.ID]

		if isExist == false {
			drs = DiceStatsticsOfCthulhu{}
		}

		if drs.Player.ID == "" {
			drs.Player.ID = drl.Player.ID
			drs.Player.Name = drl.Player.Name
		}
		if strings.Contains(drl.Result, "決定的成功") {
			drs.Critical = append(drs.Critical, drl.Command)
		} else if strings.Contains(drl.Result, "致命的失敗") {
			drs.Fumble = append(drs.Fumble, drl.Command)
		} else {

		}

		diceResultStatstics[drl.Player.ID] = drs
	}

	// クトゥルフダイスの集計
	for _, drl := range DiceResultLogOfCthulhus {
		drs, isExist := diceResultStatstics[drl.Player.ID]

		if isExist == false {
			drs = DiceStatsticsOfCthulhu{}
		}

		if diceResultStatstics[drl.Player.ID].Player.ID == "" {
			drs.Player.ID = drl.Player.ID
			drs.Player.Name = drl.Player.Name
		}
		if strings.Contains(drl.Result, "決定的成功") {
			drs.Critical = append(drs.Critical, drl.Command)
		} else if strings.Contains(drl.Result, "致命的失敗") {
			drs.Fumble = append(drs.Fumble, drl.Command)
		} else {

		}

		diceResultStatstics[drl.Player.ID] = drs
	}

	// 集計結果の構築

	if 0 < len(diceResultStatstics) {
		handlerResult.Normal.Content = "\r\n===================="
		for _, drs := range diceResultStatstics {
			handlerResult.Normal.Content += "\r\n【" + drs.Player.Name + "】\r\n"
			if len(drs.Critical) > 0 {
				handlerResult.Normal.Content += "●決定的成功：\r\n"
				handlerResult.Normal.Content += strings.Join(drs.Critical, ", ")
				handlerResult.Normal.Content += "\r\n"
			}
			if len(drs.Fumble) > 0 {
				handlerResult.Normal.Content += "●致命的失敗：\r\n"
				handlerResult.Normal.Content += strings.Join(drs.Fumble, ", ")
				handlerResult.Normal.Content += "\r\n"
			}
		}
		handlerResult.Normal.Content += "====================\r\n"
	} else {
		handlerResult.Normal.Content += "No data."
	}

	return handlerResult
}
