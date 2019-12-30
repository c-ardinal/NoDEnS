package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/c-ardinal/Nodens/lib/cthulhu"
	"github.com/c-ardinal/Nodens/lib/nodens"
)

// cmdRegistryCharacter キャラシ連携ハンドラ
func cmdRegistryCharacter(opt []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var returnMes string
	var cd *cthulhu.Character
	if cthulhu.CheckDuplicateSession(ch.ID) == true {
		if cthulhu.CheckDuplicateCharacter(ch.ID, mes.Author.ID) == true {
			return "Character already exists.", "", nil
		}
		cas, err := cthulhu.GetCharSheetFromURL(opt[0])
		if err != nil {
			return "Registry failed.", "", err
		}
		cd = cthulhu.GetCharDataFromCharSheet(cas, mes.Author.Username, mes.Author.ID)
		(*cs).Pc[mes.Author.ID] = cd
	} else if cthulhu.GetParentIDFromChildID(ch.ID) != "" {
		if cthulhu.CheckDuplicateNPCharacter(cthulhu.GetParentIDFromChildID(ch.ID), mes.Author.ID) == true {
			return "Character already exists.", "", nil
		}
		cas, err := cthulhu.GetCharSheetFromURL(opt[0])
		if err != nil {
			return "Registry failed.", "", err
		}
		cd = cthulhu.GetCharDataFromCharSheet(cas, mes.Author.Username, mes.Author.ID)
		(*cs).Npc[mes.Author.ID] = cd
	} else {
		return "Session not found.", "", nil
	}
	returnMes = "\r\n====================\r\n"
	returnMes += "**[名 前]** " + cd.Personal.Name + "\r\n"
	returnMes += "**[年 齢]** " + strconv.Itoa(cd.Personal.Age) + "歳\r\n"
	returnMes += "**[性 別]** " + cd.Personal.Sex + "\r\n"
	returnMes += "**[職 業]** " + cd.Personal.Job + "\r\n"
	for _, cdan := range cthulhu.CdAbilityNameList {
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

// cmdCharaNumCheck 能力値確認ハンドラ
func cmdCharaNumCheck(opt []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var chara *cthulhu.Character
	var exist bool
	if cs == nil {
		return "Character not registried.", "", nil
	}
	if cthulhu.GetParentIDFromChildID(ch.ID) != "" {
		chara, exist = (*cs).Npc[mes.Author.ID]
	} else {
		chara, exist = (*cs).Pc[mes.Author.ID]
	}
	if exist == false {
		return "Character not found.", "", nil
	}
	initNum := cthulhu.GetSkillNum(chara, opt[0], "init")
	if initNum == "-1" {
		return "Skill not found.", "", nil
	}
	startNum := cthulhu.GetSkillNum(chara, opt[0], "sum")
	nowNum := cthulhu.GetSkillNum(chara, opt[0], "now")

	return "[" + opt[0] + "] Init( " + initNum + " ), Start( " + startNum + "), Now( " + nowNum + " )", "", nil
}

// cmdCharaNumControl 能力値操作ハンドラ
func cmdCharaNumControl(opt []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var chara *cthulhu.Character
	var exist bool
	if cs == nil {
		return "Character not registried.", "", nil
	}
	if cthulhu.GetParentIDFromChildID(ch.ID) != "" {
		chara, exist = (*cs).Npc[mes.Author.ID]
	} else {
		chara, exist = (*cs).Pc[mes.Author.ID]
	}
	if exist == false {
		return "Character not found.", "", nil
	}
	oldNum := cthulhu.GetSkillNum(chara, opt[0], "now")
	if oldNum == "-1" {
		return "Skill not found.", "", nil
	}
	diffRegex := regexp.MustCompile("^[+-]?[0-9]+$")
	if diffRegex.MatchString(opt[1]) == false {
		return "Invalid diff num.", "", nil
	}
	newNum := cthulhu.AddSkillNum(chara, opt[0], opt[1])
	return "[" + opt[0] + "] " + oldNum + " => " + newNum + " (Diff: " + opt[1] + ")", "", nil
}

// cmdLinkRoll キャラシ連携ダイスロールハンドラ
func cmdLinkRoll(opt []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	if cs == nil {
		return "PC not registried.", "", nil
	}
	pc, exist := (*cs).Pc[mes.Author.ID]
	if exist == false {
		return "PC not found.", "", nil
	}
	diceCmd := "CCB<=" + opt[0]
	exRegex := regexp.MustCompile("[^\\+\\-\\*\\/ 　]+")
	ignoreRegex := regexp.MustCompile("^[0-9]+$")
	for _, ex := range exRegex.FindAllString(opt[0], -1) {
		if ignoreRegex.MatchString(ex) == false {
			exNum := cthulhu.GetSkillNum(pc, ex, "now")
			if exNum == "-1" {
				return "Skill not found.", "", nil
			}
			diceCmd = strings.Replace(diceCmd, ex, exNum, -1)
		}
	}
	rollResult, err := nodens.ExecuteDiceRoll(nodens.GetConfig().EndPoint, (*cs).Scenario.System, diceCmd)
	return rollResult.Result, "", err
}

// cmdSecretLinkRoll キャラシ連携Secretダイスロールハンドラ
func cmdSecretLinkRoll(opt []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	if cs == nil {
		return "NPC not registried.", "", nil
	}
	pc, exist := (*cs).Npc[mes.Author.ID]
	if exist == false {
		return "NPC not found.", "", nil
	}
	diceCmd := "SCCB<=" + opt[0]
	exRegex := regexp.MustCompile("[^\\+\\-\\*\\/ 　]+")
	ignoreRegex := regexp.MustCompile("^[0-9]+$")
	for _, ex := range exRegex.FindAllString(opt[0], -1) {
		if ignoreRegex.MatchString(ex) == false {
			exNum := cthulhu.GetSkillNum(pc, ex, "now")
			if exNum == "-1" {
				return "Skill not found.", "", nil
			}
			diceCmd = strings.Replace(diceCmd, ex, exNum, -1)
		}
	}
	rollResult, err := nodens.ExecuteDiceRoll(nodens.GetConfig().EndPoint, (*cs).Scenario.System, diceCmd)
	var secretMes string
	if rollResult.Secret == true {
		secretMes = "**SECRET DICE**"
	}
	return rollResult.Result, secretMes, err
}

// cmdSanCheckRoll SAN値チェック処理ハンドラ
func cmdSanCheckRoll(opt []string, cs *cthulhu.CthulhuSession, ch *discordgo.Channel, mes *discordgo.MessageCreate) (string, string, error) {
	var sucSub string
	var failSub string
	var resultMes string

	if cs == nil {
		return "PC not registried.", "", nil
	}
	pc, exist := (*cs).Pc[mes.Author.ID]
	if exist == false {
		return "PC not found.", "", nil
	}
	if len(opt) < 2 {
		return "Invalid arguments.", "", nil
	}

	orgSanNum := cthulhu.GetSkillNum(pc, "san", "now")
	sanRollCmd := "SCCB<=" + orgSanNum
	sanRollResult, err := nodens.ExecuteDiceRoll(nodens.GetConfig().EndPoint, (*cs).Scenario.System, sanRollCmd)

	if strings.Contains(sanRollResult.Result, "成功") {
		if strings.Contains(opt[0], "d") {
			sucRollResult, _ := nodens.ExecuteDiceRoll(nodens.GetConfig().EndPoint, (*cs).Scenario.System, opt[0])
			sucSub = "-" + nodens.CalcDicesSum(sucRollResult.Dices)
		} else {
			sucSub = "-" + opt[0]
		}
		newNum := cthulhu.AddSkillNum(pc, "san", sucSub)
		resultMes = "sanc > 成功 > SAN: " + orgSanNum + " -> " + newNum + " ( " + sucSub + " )"
	} else {
		if strings.Contains(opt[1], "d") {
			failRollResult, _ := nodens.ExecuteDiceRoll(nodens.GetConfig().EndPoint, (*cs).Scenario.System, opt[1])
			failSub = "-" + nodens.CalcDicesSum(failRollResult.Dices)
		} else {
			failSub = "-" + opt[1]
		}
		newNum := cthulhu.AddSkillNum(pc, "san", failSub)
		resultMes = "sanc > 失敗 > SAN: " + orgSanNum + " -> " + newNum + " ( " + failSub + " )"
	}

	return resultMes, "", err
}
