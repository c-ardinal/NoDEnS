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
	oldNum := cthulhu.GetSkillNum(chara, opt[0])
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
			exNum := cthulhu.GetSkillNum(pc, ex)
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
			exNum := cthulhu.GetSkillNum(pc, ex)
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
