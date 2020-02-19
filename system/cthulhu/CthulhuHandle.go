package cthulhu

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
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

// CmdRegistryCharacter キャラシ連携ハンドラ
func CmdRegistryCharacter(opt []string, cs *core.Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string
	var cd *CharacterOfCthulhu
	if len(opt) == 0 {
		return "Invalid arguments.", "", nil
	}
	if core.CheckExistSession(ch.ID) == true {
		if core.CheckExistCharacter(ch.ID, mes.Author.ID) == true {
			return "Character already exists.", "", nil
		}
		cas, err := GetCharSheetFromURL(opt[0])
		if err != nil {
			return "Registry failed.", "", err
		}
		cd = GetCharDataFromCharSheet(cas, mes.Author.Username, mes.Author.ID)
		(*cs).Pc[mes.Author.ID] = cd
	} else if core.GetParentIDFromChildID(ch.ID) != "" {
		if core.CheckExistNPCharacter(core.GetParentIDFromChildID(ch.ID), mes.Author.ID) == true {
			return "Character already exists.", "", nil
		}
		cas, err := GetCharSheetFromURL(opt[0])
		if err != nil {
			return "Registry failed.", "", err
		}
		cd = GetCharDataFromCharSheet(cas, mes.Author.Username, mes.Author.ID)
		(*cs).Npc[mes.Author.ID] = cd
	} else {
		return "Session not found.", "", nil
	}
	returnMes = "\r\n====================\r\n"
	returnMes += "**[名 前]** " + cd.Personal.Name + "\r\n"
	returnMes += "**[年 齢]** " + strconv.Itoa(cd.Personal.Age) + "歳\r\n"
	returnMes += "**[性 別]** " + cd.Personal.Sex + "\r\n"
	returnMes += "**[職 業]** " + cd.Personal.Job + "\r\n"
	for _, cdan := range CdAbilityNameList {
		a := cd.Ability[cdan]
		if a.Now == a.Init {
			returnMes += "**[ " + a.Name + " ]** " + strconv.Itoa(a.Now) + "\r\n"
		} else {
			returnMes += "**[ " + a.Name + " ]** " + strconv.Itoa(a.Now) + " (Init: " + strconv.Itoa(a.Init) + ")\r\n"
		}
	}
	returnMes += "**[メ モ]** \r\n" + cd.Memo + "\r\n"
	returnMes += "====================\r\n"
	return returnMes, "", nil
}

// CmdCharaNumCheck 能力値確認ハンドラ
func CmdCharaNumCheck(opt []string, cs *core.Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var chara *CharacterOfCthulhu
	var exist bool
	if len(opt) == 0 {
		return "Invalid arguments.", "", nil
	}
	if cs == nil {
		return "Character not registried.", "", nil
	}
	if core.GetParentIDFromChildID(ch.ID) != "" {
		chara, exist = (*cs).Npc[mes.Author.ID].(*CharacterOfCthulhu)
	} else {
		chara, exist = (*cs).Pc[mes.Author.ID].(*CharacterOfCthulhu)
	}
	if exist == false {
		return "Character not found.", "", nil
	}
	initNum := GetSkillNum(chara, opt[0], "init")
	if initNum == "-1" {
		return "Skill not found.", "", nil
	}
	startNum := GetSkillNum(chara, opt[0], "sum")
	nowNum := GetSkillNum(chara, opt[0], "now")

	returnMes := "[" + opt[0] + "] Init( " + initNum + " ), Start( " + startNum + "), Now( " + nowNum + " )"

	return returnMes, "", nil
}

// CmdCharaNumControl 能力値操作ハンドラ
func CmdCharaNumControl(opt []string, cs *core.Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var chara *CharacterOfCthulhu
	var exist bool
	if len(opt) < 2 {
		return "Invalid arguments.", "", nil
	}
	if cs == nil {
		return "Character not registried.", "", nil
	}
	if core.GetParentIDFromChildID(ch.ID) != "" {
		chara, exist = (*cs).Npc[mes.Author.ID].(*CharacterOfCthulhu)
	} else {
		chara, exist = (*cs).Pc[mes.Author.ID].(*CharacterOfCthulhu)
	}
	if exist == false {
		return "Character not found.", "", nil
	}
	oldNum := GetSkillNum(chara, opt[0], "now")
	if oldNum == "-1" {
		return "Skill not found.", "", nil
	}
	diffRegex := regexp.MustCompile("^[+-]?[0-9]+$")
	if diffRegex.MatchString(opt[1]) == false {
		return "Invalid diff num.", "", nil
	}
	newNum := AddSkillNum(chara, opt[0], opt[1])

	returnMes := "[" + opt[0] + "] " + oldNum + " => " + newNum + " (Diff: " + opt[1] + ")"

	return returnMes, "", nil
}

// CmdLinkRoll キャラシ連携ダイスロールハンドラ
func CmdLinkRoll(opt []string, cs *core.Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	if len(opt) == 0 {
		return "Invalid arguments.", "", nil
	}
	if cs == nil {
		return "PC not registried.", "", nil
	}
	pc, exist := (*cs).Pc[mes.Author.ID].(*CharacterOfCthulhu)
	if exist == false {
		return "PC not found.", "", nil
	}
	diceCmd := "CCB<=" + opt[0]
	exRegex := regexp.MustCompile("[^\\+\\-\\*\\/ 　]+")
	ignoreRegex := regexp.MustCompile("^[0-9]+$")
	for _, ex := range exRegex.FindAllString(opt[0], -1) {
		if ignoreRegex.MatchString(ex) == false {
			exNum := GetSkillNum(pc, ex, "now")
			if exNum == "-1" {
				return "Skill not found.", "", nil
			}
			diceCmd = strings.Replace(diceCmd, ex, exNum, -1)
		}
	}
	rollResult, err := core.ExecuteDiceRoll(core.GetConfig().EndPoint, (*cs).Scenario.System, diceCmd)

	if err == nil {
		const format = "2006/01/02_15:04:05"
		parsedTime, _ := mes.Timestamp.Parse()
		var cthulhuDiceResultLog DiceResultLogOfCthulhu

		cthulhuDiceResultLog.Player.ID = mes.Author.ID
		cthulhuDiceResultLog.Player.Name = mes.Author.Username
		cthulhuDiceResultLog.Time = parsedTime.Format(format)
		cthulhuDiceResultLog.Command = mes.Content
		cthulhuDiceResultLog.Result = rollResult.Result
		DiceResultLogOfCthulhus = append(DiceResultLogOfCthulhus, cthulhuDiceResultLog)
	}

	return rollResult.Result, "", err
}

// CmdSecretLinkRoll キャラシ連携Secretダイスロールハンドラ
func CmdSecretLinkRoll(opt []string, cs *core.Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	if len(opt) == 0 {
		return "Invalid arguments.", "", nil
	}
	if cs == nil {
		return "NPC not registried.", "", nil
	}
	pc, exist := (*cs).Npc[mes.Author.ID].(*CharacterOfCthulhu)
	if exist == false {
		return "NPC not found.", "", nil
	}
	diceCmd := "SCCB<=" + opt[0]
	exRegex := regexp.MustCompile("[^\\+\\-\\*\\/ 　]+")
	ignoreRegex := regexp.MustCompile("^[0-9]+$")
	for _, ex := range exRegex.FindAllString(opt[0], -1) {
		if ignoreRegex.MatchString(ex) == false {
			exNum := GetSkillNum(pc, ex, "now")
			if exNum == "-1" {
				return "Skill not found.", "", nil
			}
			diceCmd = strings.Replace(diceCmd, ex, exNum, -1)
		}
	}
	rollResult, err := core.ExecuteDiceRoll(core.GetConfig().EndPoint, (*cs).Scenario.System, diceCmd)
	var secretMes string
	if rollResult.Secret == true {
		secretMes = "**SECRET DICE**"
	}
	return rollResult.Result, secretMes, err
}

// CmdSanCheckRoll SAN値チェック処理ハンドラ
func CmdSanCheckRoll(opt []string, cs *core.Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var sucSub string
	var failSub string
	var resultMes string

	if len(opt) < 2 {
		return "Invalid arguments.", "", nil
	}
	if cs == nil {
		return "PC not registried.", "", nil
	}
	pc, exist := (*cs).Pc[mes.Author.ID].(*CharacterOfCthulhu)
	if exist == false {
		return "PC not found.", "", nil
	}

	orgSanNum := GetSkillNum(pc, "san", "now")
	sanRollCmd := "SCCB<=" + orgSanNum
	sanRollResult, err := core.ExecuteDiceRoll(core.GetConfig().EndPoint, (*cs).Scenario.System, sanRollCmd)

	if strings.Contains(sanRollResult.Result, "成功") {
		if strings.Contains(opt[0], "d") {
			sucRollResult, _ := core.ExecuteDiceRoll(core.GetConfig().EndPoint, (*cs).Scenario.System, opt[0])
			sucSub = "-" + core.CalcDicesSum(sucRollResult.Dices)
		} else {
			sucSub = "-" + opt[0]
		}
		newNum := AddSkillNum(pc, "san", sucSub)
		resultMes = "sanc > [ " + sanRollResult.Result + " ] >> SAN: " + orgSanNum + " -> " + newNum + " ( " + sucSub + " )"
	} else {
		if strings.Contains(opt[1], "d") {
			failRollResult, _ := core.ExecuteDiceRoll(core.GetConfig().EndPoint, (*cs).Scenario.System, opt[1])
			failSub = "-" + core.CalcDicesSum(failRollResult.Dices)
		} else {
			failSub = "-" + opt[1]
		}
		newNum := AddSkillNum(pc, "san", failSub)
		resultMes = "sanc >> [ " + sanRollResult.Result + " ] >> SAN: " + orgSanNum + " -> " + newNum + " ( " + failSub + " )"
	}

	return resultMes, "", err
}

// CmdShowStatistics ダイスロール統計表示処理
func CmdShowStatistics(opt []string, cs *core.Session, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
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
			drs.Critical = append(drs.Critical, "["+drl.Time+"] "+drl.Command+" ➡ "+drl.Result)
		} else if strings.Contains(drl.Result, "致命的失敗") {
			drs.Fumble = append(drs.Fumble, "["+drl.Time+"] "+drl.Command+" ➡ "+drl.Result)
		} else if strings.Contains(drl.Result, "成功") {
			drs.Success = append(drs.Success, "["+drl.Time+"] "+drl.Command+" ➡ "+drl.Result)
		} else if strings.Contains(drl.Result, "失敗") {
			drs.Fail = append(drs.Fail, "["+drl.Time+"] "+drl.Command+" ➡ "+drl.Result)
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
			drs.Critical = append(drs.Critical, "["+drl.Time+"] "+drl.Command+" ➡ "+drl.Result)
		} else if strings.Contains(drl.Result, "致命的失敗") {
			drs.Fumble = append(drs.Fumble, "["+drl.Time+"] "+drl.Command+" ➡ "+drl.Result)
		} else if strings.Contains(drl.Result, "成功") {
			drs.Success = append(drs.Success, "["+drl.Time+"] "+drl.Command+" ➡ "+drl.Result)
		} else if strings.Contains(drl.Result, "失敗") {
			drs.Fail = append(drs.Fail, "["+drl.Time+"] "+drl.Command+" ➡ "+drl.Result)
		} else {

		}

		diceResultStatstics[drl.Player.ID] = drs
	}

	// 集計結果の構築
	var returnMes string
	if 0 < len(diceResultStatstics) {
		returnMes = "\r\n===================="
		for _, drs := range diceResultStatstics {
			returnMes += "\r\n【" + drs.Player.Name + "】\r\n"
			returnMes += "    ●決定的成功：\r\n"
			for _, critical := range drs.Critical {
				returnMes += "      ・" + critical + "\r\n"
			}
			returnMes += "    ●成功：\r\n"
			for _, success := range drs.Success {
				returnMes += "      ・" + success + "\r\n"
			}
			returnMes += "    ●失敗：\r\n"
			for _, fail := range drs.Fail {
				returnMes += "      ・" + fail + "\r\n"
			}
			returnMes += "    ●致命的失敗：\r\n"
			for _, fumble := range drs.Fumble {
				returnMes += "      ・" + fumble + "\r\n"
			}
		}
		returnMes += "====================\r\n"
	} else {
		returnMes += "No data."
	}

	return returnMes, "", nil
}
