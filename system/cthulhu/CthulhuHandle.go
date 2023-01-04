package cthulhu

import (
	"regexp"
	"strconv"
	"strings"

	"Nodens/core"

	"github.com/bwmarrin/discordgo"
)

// DiceResultLogOfCthulhu ダイスロール実行ログ型
type DiceResultLogOfCthulhu struct {
	Player  core.NaID
	Time    string
	Command string
	Result  string
}

// DiceStatisticsOfCthulhu ダイスロール統計型
type DiceStatisticsOfCthulhu struct {
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
func CmdRegistryCharacter(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var cd *CharacterOfCthulhu
	var returnMes string

	if len(opt) == 0 {
		returnMes = "Invalid arguments."
	} else {
		if core.CheckExistSession(md.ChannelID) == true {
			/* 親セッションでキャラ登録コマンドが来た場合，PCとして登録する */
			if core.CheckExistCharacter(md.ChannelID, md.AuthorID) == true {
				returnMes = "Character already exists."
			} else {
				cas, err := GetCharSheetFromURL(opt[0])
				if err != nil {
					returnMes = "Registry failed."
					handlerResult.Error = err
				} else {
					cd = GetCharDataFromCharSheet(cas, md.AuthorName, md.AuthorID)
					(*cd).URL = opt[0]
					(*cs).Pc[md.AuthorID] = cd
				}
			}
		} else if core.GetParentIDFromChildID(md.ChannelID) != "" {
			/* 子セッションでキャラ登録コマンドが来た場合，NPCとして登録する */
			if core.CheckExistNPCharacter(core.GetParentIDFromChildID(md.ChannelID), md.AuthorID) == true {
				returnMes = "Character already exists."
			} else {
				cas, err := GetCharSheetFromURL(opt[0])
				if err != nil {
					returnMes = "Registry failed."
					handlerResult.Error = err
				} else {
					cd = GetCharDataFromCharSheet(cas, md.AuthorName, md.AuthorID)
					(*cd).URL = opt[0]
					(*cs).Npc[md.AuthorID] = cd
				}
			}
		} else {
			returnMes = "Session not found."
		}
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = core.EnEmbed

	/* テキストメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Content = returnMes
	} else {
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
	}

	/* Embedメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Description: returnMes,
			Color:       0xff0000, // Red
		}
	} else {
		var fields []*discordgo.MessageEmbedField
		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:   "\u200B",
				Value:  "---------------------------------------------------------",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "[年 齢]",
				Value:  strconv.Itoa(cd.Personal.Age) + "歳",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "[性 別]",
				Value:  cd.Personal.Sex,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "[職 業]",
				Value:  cd.Personal.Job,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "\u200B",
				Value:  "---------------------------------------------------------",
				Inline: false,
			})

		for _, cdan := range CdAbilityNameList {
			a := cd.Ability[cdan]
			if a.Now == a.Init {
				fields = append(fields,
					&discordgo.MessageEmbedField{
						Name:   "[ " + a.Name + " ]",
						Value:  strconv.Itoa(cd.Ability[a.Name].Now),
						Inline: true,
					})
			} else {
				fields = append(fields,
					&discordgo.MessageEmbedField{
						Name:   "[ " + a.Name + " ]",
						Value:  strconv.Itoa(cd.Ability[a.Name].Now) + " (Init: " + strconv.Itoa(a.Init) + ")",
						Inline: true,
					})
			}
		}
		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:   "\u200B",
				Value:  "---------------------------------------------------------",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "[メ モ]",
				Value:  cd.Memo,
				Inline: false,
			})

		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Title:  cd.Personal.Name,
			URL:    opt[0],
			Color:  0x00ff00, // Green
			Fields: fields,
		}

	}
	return handlerResult
}

// CmdCharaNumCheck 能力値確認ハンドラ
func CmdCharaNumCheck(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var chara *CharacterOfCthulhu
	var exist bool
	var initNum string
	var startNum string
	var nowNum string
	var returnMes string

	if len(opt) == 0 {
		returnMes = "Invalid arguments."
	} else {
		if cs == nil {
			returnMes = "Character not registered."
		} else {
			if core.GetParentIDFromChildID(md.ChannelID) != "" {
				chara, exist = (*cs).Npc[md.AuthorID].(*CharacterOfCthulhu)
			} else {
				chara, exist = (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
			}
			if exist == false {
				returnMes = "Character not found."
			} else {
				initNum = GetSkillNum(chara, opt[0], "init")
				if initNum == "-1" {
					returnMes = "Skill not found."
				} else {
					startNum = GetSkillNum(chara, opt[0], "sum")
					nowNum = GetSkillNum(chara, opt[0], "now")
				}
			}
		}

	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = core.EnEmbed

	/* テキストメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Content = returnMes
	} else {
		handlerResult.Normal.Content = "[" + opt[0] + "] Init( " + initNum + " ), Start( " + startNum + "), Now( " + nowNum + " )"
	}

	/* Embedメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Description: returnMes,
			Color:       0xff0000, // Red
		}
	} else {
		var fields []*discordgo.MessageEmbedField
		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:   "[ Init ]",
				Value:  initNum,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "[ Start ]",
				Value:  startNum,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "[ Now ]",
				Value:  nowNum,
				Inline: true,
			},
		)
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Title:  "< " + opt[0] + " >",
			Color:  0x00ff00, // Green
			Fields: fields,
		}
	}

	return handlerResult
}

// CmdCharaNumControl 能力値操作ハンドラ
func CmdCharaNumControl(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var chara *CharacterOfCthulhu
	var exist bool
	var oldNum string
	var newNum string
	var diffCmd string
	var rollResultMessage string
	var returnMes string

	if len(opt) < 2 {
		returnMes = "Invalid arguments."
	} else {
		if cs == nil {
			returnMes = "Character not registered."
		} else {
			if core.GetParentIDFromChildID(md.ChannelID) != "" {
				chara, exist = (*cs).Npc[md.AuthorID].(*CharacterOfCthulhu)
			} else {
				chara, exist = (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
			}
			if exist == false {
				returnMes = "Character not found."
			} else {
				oldNum = GetSkillNum(chara, opt[0], "now")
				if oldNum == "-1" {
					returnMes = "Skill not found."
				} else {
					diffRegex := regexp.MustCompile("^[+-]?[0-9]+$")
					diffCmd = opt[1]
					if diffRegex.MatchString(diffCmd) == false {
						minusFlag := false
						if strings.Contains(diffCmd, "-") {
							diffCmd = strings.ReplaceAll(diffCmd, "-", "")
							minusFlag = true
						}
						rollResult, err := core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, diffCmd)
						rollResultMessage = rollResult.Result
						if err != nil {
							returnMes = "Invalid diff num."
							handlerResult.Error = err
						} else {
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
					}
					newNum = AddSkillNum(chara, opt[0], diffCmd)
				}
			}
		}
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = core.EnEmbed

	/* テキストメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Content = returnMes
	} else {
		handlerResult.Normal.Content = "[" + opt[0] + "] " + oldNum + " => " + newNum + " (Diff: " + diffCmd + ")"
	}

	/* Embedメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Description: returnMes,
			Color:       0xff0000, // Red
		}
	} else {
		var fields []*discordgo.MessageEmbedField
		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:   "[ Before ]",
				Value:  oldNum,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "[ After ]",
				Value:  newNum,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "[ Diff ]",
				Value:  diffCmd,
				Inline: true,
			},
		)
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Title:       "< " + opt[0] + " >",
			Description: rollResultMessage,
			Color:       0x00ff00, // Green
			Fields:      fields,
		}
	}

	return handlerResult
}

// CmdLinkRoll キャラシ連携ダイスロールハンドラ
func CmdLinkRoll(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var rollResult core.BCDiceRollResult
	var err error
	var returnMes string

	if len(opt) == 0 {
		returnMes = "Invalid arguments."
	} else {
		if cs == nil {
			returnMes = "PC not registered."
		} else {
			pc, exist := (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
			if exist == false {
				returnMes = "PC not found."
			} else {
				diceCmd := "CCB<=" + opt[0]
				exRegex := regexp.MustCompile("[^\\+\\-\\*\\/ 　]+")
				ignoreRegex := regexp.MustCompile("^[0-9]+$")
				for _, ex := range exRegex.FindAllString(opt[0], -1) {
					if ignoreRegex.MatchString(ex) == false {
						exNum := GetSkillNum(pc, ex, "now")
						if exNum == "-1" {
							returnMes = "Skill not found."
						} else {
							diceCmd = strings.Replace(diceCmd, ex, exNum, -1)
							rollResult, err = core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, diceCmd)
							if err != nil {
								handlerResult.Error = err
							} else {
								/* Non process */
							}
						}
					}
				}
			}
		}
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = core.EnEmbed

	/* テキストメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Content = returnMes
	} else {
		handlerResult.Normal.Content = rollResult.Result
	}

	/* Embedメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Description: returnMes,
			Color:       0xff0000, // Red
		}
	} else {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Title:       "< " + opt[0] + " >",
			Description: rollResult.Result,
			Color:       0x00ff00, // Green
		}
	}

	if (err == nil) && (len(opt) > 0) {
		//const format = "2006/01/02_15:04:05"
		//parsedTime, _ := mes.Timestamp.Parse()
		var cthulhuDiceResultLog DiceResultLogOfCthulhu

		cthulhuDiceResultLog.Player.ID = md.AuthorID
		cthulhuDiceResultLog.Player.Name = md.AuthorName
		//cthulhuDiceResultLog.Time = parsedTime.Format(format)
		cthulhuDiceResultLog.Command = opt[0]
		cthulhuDiceResultLog.Result = rollResult.Result
		DiceResultLogOfCthulhus = append(DiceResultLogOfCthulhus, cthulhuDiceResultLog)
	}

	return handlerResult
}

// CmdSecretLinkRoll キャラシ連携Secretダイスロールハンドラ
func CmdSecretLinkRoll(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var rollResult core.BCDiceRollResult
	var err error
	var returnMes string

	if len(opt) == 0 {
		returnMes = "Invalid arguments."
	} else {
		if cs == nil {
			returnMes = "NPC not registered."
		} else {
			pc, exist := (*cs).Npc[md.AuthorID].(*CharacterOfCthulhu)
			if exist == false {
				returnMes = "NPC not found."
			} else {
				diceCmd := "SCCB<=" + opt[0]
				exRegex := regexp.MustCompile("[^\\+\\-\\*\\/ 　]+")
				ignoreRegex := regexp.MustCompile("^[0-9]+$")
				for _, ex := range exRegex.FindAllString(opt[0], -1) {
					if ignoreRegex.MatchString(ex) == false {
						exNum := GetSkillNum(pc, ex, "now")
						if exNum == "-1" {
							returnMes = "Skill not found."
						} else {
							diceCmd = strings.Replace(diceCmd, ex, exNum, -1)
							rollResult, err = core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, diceCmd)
							if err != nil {
								handlerResult.Error = err
							} else {
								/* Non process */
							}
						}
					}
				}
			}
		}
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = core.EnEmbed

	/* テキストメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Content = returnMes
	} else {
		handlerResult.Normal.Content = rollResult.Result
	}

	/* Embedメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Description: returnMes,
			Color:       0xff0000, // Red
		}
	} else {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Title:       "< " + opt[0] + " >",
			Description: rollResult.Result,
			Color:       0x00ff00, // Green
		}
	}

	if rollResult.Secret == true {
		/* 有効にするメッセージタイプ */
		handlerResult.Secret.EnableType = core.EnEmbed

		/* テキストメッセージ */
		handlerResult.Secret.Content = "**SECRET DICE**"

		/* Embedメッセージ */
		handlerResult.Secret.Embed = &discordgo.MessageEmbed{
			Title: "SECRET DICE",
			Color: 0xffff00, // Yellow
		}
	}

	return handlerResult
}

// CmdSanCheckRoll SAN値チェック処理ハンドラ
func CmdSanCheckRoll(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var sanRollResult core.BCDiceRollResult
	var sucRollResult core.BCDiceRollResult
	var failRollResult core.BCDiceRollResult
	var err error
	var orgSanNum string
	var sanSub string
	var newNum string
	var returnMes string

	if len(opt) < 2 {
		returnMes = "Invalid arguments."
	} else {
		if cs == nil {
			returnMes = "PC not registered."
		} else {
			pc, exist := (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
			if exist == false {
				returnMes = "PC not found."
			} else {
				orgSanNum = GetSkillNum(pc, "san", "now")
				sanRollCmd := "SCCB<=" + orgSanNum
				sanRollResult, err = core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, sanRollCmd)
				if err != nil {
					returnMes = "Server error."
					handlerResult.Error = err
				} else {
					if strings.Contains(sanRollResult.Result, "成功") || strings.Contains(sanRollResult.Result, "スペシャル") {
						if strings.Contains(opt[0], "d") {
							sucRollResult, err = core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, opt[0])
							if err != nil {
								returnMes = "Server error."
								handlerResult.Error = err
							} else {
								sanSub = "-" + core.CalcDicesSum(sucRollResult.Dices)
							}

						} else {
							sanSub = "-" + opt[0]
						}
						newNum = AddSkillNum(pc, "san", sanSub)
					} else {
						if strings.Contains(opt[1], "d") {
							failRollResult, err = core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, opt[1])
							if err != nil {
								returnMes = "Server error."
								handlerResult.Error = err
							} else {
								sanSub = "-" + core.CalcDicesSum(failRollResult.Dices)
							}
						} else {
							sanSub = "-" + opt[1]
						}
						newNum = AddSkillNum(pc, "san", sanSub)
					}
				}
			}
		}
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = core.EnEmbed

	/* テキストメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Content = returnMes
	} else {
		handlerResult.Normal.Content = "sanc > [ " + sanRollResult.Result + " ] >> SAN: " + orgSanNum + " -> " + newNum + " ( " + sanSub + " )"
	}

	/* Embedメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Description: returnMes,
			Color:       0xff0000, // Red
		}
	} else {
		var fields []*discordgo.MessageEmbedField
		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:   "[ Before ]",
				Value:  orgSanNum,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "[ After ]",
				Value:  newNum,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "[ Diff ]",
				Value:  sanSub,
				Inline: true,
			},
		)
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Title:       "< SANc >",
			Description: sanRollResult.Result,
			Color:       0x00ff00, // Green
			Fields:      fields,
		}
	}
	return handlerResult
}

// CmdShowStatistics ダイスロール統計表示処理
func CmdShowStatistics(opt []string, cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var diceResultLogs = core.GetDiceResultLogs()
	var diceResultStatistics = map[string]DiceStatisticsOfCthulhu{}

	// 共通ダイスの集計
	for _, drl := range diceResultLogs {
		drs, isExist := diceResultStatistics[drl.Player.ID]

		if isExist == false {
			drs = DiceStatisticsOfCthulhu{}
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

		diceResultStatistics[drl.Player.ID] = drs
	}

	// クトゥルフダイスの集計
	for _, drl := range DiceResultLogOfCthulhus {
		drs, isExist := diceResultStatistics[drl.Player.ID]

		if isExist == false {
			drs = DiceStatisticsOfCthulhu{}
		}

		if diceResultStatistics[drl.Player.ID].Player.ID == "" {
			drs.Player.ID = drl.Player.ID
			drs.Player.Name = drl.Player.Name
		}
		if strings.Contains(drl.Result, "決定的成功") {
			drs.Critical = append(drs.Critical, drl.Command)
		} else if strings.Contains(drl.Result, "致命的失敗") {
			drs.Fumble = append(drs.Fumble, drl.Command)
		} else {

		}

		diceResultStatistics[drl.Player.ID] = drs
	}

	// 集計結果の構築

	if 0 < len(diceResultStatistics) {
		handlerResult.Normal.Content = "\r\n===================="
		for _, drs := range diceResultStatistics {
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
