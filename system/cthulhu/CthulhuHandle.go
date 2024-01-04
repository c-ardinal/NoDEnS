package cthulhu

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"Nodens/core"

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
/* 関数定義                                                                 */
/****************************************************************************/

// キャラクターシート連携ハンドラ
func CmdRegistryCharacter(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var cd *CharacterOfCthulhu
	var urlStr string
	var returnMes string

	if len(md.Options) == 0 {
		returnMes = "Invalid arguments."
		handlerResult.Error = errors.New(returnMes)
	} else {
		urlStr = md.Options[0].Value
		if core.CheckExistParentSession(md.ChannelID) == true {
			/* 親セッションでキャラクター登録コマンドが来た場合，PCとして登録する */
			if core.CheckExistCharacter(md.ChannelID, md.AuthorID) == true {
				returnMes = "Character already exists."
				handlerResult.Error = errors.New(returnMes)
			} else {
				cas, err := GetCharSheetFromURL(urlStr)
				if err != nil {
					returnMes = "Registry failed."
					handlerResult.Error = err
				} else {
					cd = GetCharDataFromCharSheet(cas, md.AuthorName, md.AuthorID)
					(*cd).URL = urlStr
					(*cs).Pc[md.AuthorID] = cd
				}
			}
		} else if core.GetParentIDFromChildID(md.ChannelID) != "" {
			/* 子セッションでキャラクター登録コマンドが来た場合，NPCとして登録する */
			if core.CheckExistNPCharacter(core.GetParentIDFromChildID(md.ChannelID), md.AuthorID) == true {
				returnMes = "Character already exists."
				handlerResult.Error = errors.New(returnMes)
			} else {
				cas, err := GetCharSheetFromURL(urlStr)
				if err != nil {
					returnMes = "Registry failed."
					handlerResult.Error = err
				} else {
					cd = GetCharDataFromCharSheet(cas, md.AuthorName, md.AuthorID)
					(*cd).URL = urlStr
					(*cs).Npc[md.AuthorID] = cd
				}
			}
		} else {
			returnMes = "Session not found."
			handlerResult.Error = errors.New(returnMes)
		}
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Normal.EnableType = core.EnEmbed

	/* テキストメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Content = returnMes
	} else {
		handlerResult.Normal.Content = "\n====================\n"
		handlerResult.Normal.Content += "**[名 前]** " + cd.Personal.Name + "\n"
		handlerResult.Normal.Content += "**[年 齢]** " + strconv.Itoa(cd.Personal.Age) + "歳\n"
		handlerResult.Normal.Content += "**[性 別]** " + cd.Personal.Sex + "\n"
		handlerResult.Normal.Content += "**[職 業]** " + cd.Personal.Job + "\n"
		for _, cdan := range GetCdAbilityNameList() {
			a := cd.Ability[cdan]
			if a.Now == a.Init {
				handlerResult.Normal.Content += "**[ " + a.Name + " ]** " + strconv.Itoa(a.Now) + "\n"
			} else {
				handlerResult.Normal.Content += "**[ " + a.Name + " ]** " + strconv.Itoa(a.Now) + " (Init: " + strconv.Itoa(a.Init) + ")\n"
			}
		}
		handlerResult.Normal.Content += "**[メ モ]** \n" + cd.Memo + "\n"
		handlerResult.Normal.Content += "====================\n"
	}

	/* Embedメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Description: returnMes,
			Color:       core.EnColorRed,
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

		for _, cdan := range GetCdAbilityNameList() {
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
			URL:    urlStr,
			Color:  core.EnColorGreen,
			Fields: fields,
		}

	}
	return handlerResult
}

// 能力値確認ハンドラ
func CmdCharaNumCheck(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var skillName string
	var initNum string
	var startNum string
	var nowNum string
	var returnMes string

	if len(md.Options[0].Value) == 0 {
		returnMes = "Invalid arguments."
		handlerResult.Error = errors.New(returnMes)
	} else {
		skillName = md.Options[0].Value
		if cs == nil {
			returnMes = "Character not registered."
			handlerResult.Error = errors.New(returnMes)
		} else {
			var chara *CharacterOfCthulhu
			var exist bool
			if core.GetParentIDFromChildID(md.ChannelID) != "" {
				chara, exist = (*cs).Npc[md.AuthorID].(*CharacterOfCthulhu)
			} else {
				chara, exist = (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
			}
			if exist == false {
				returnMes = "Character not found."
			} else {
				initNum = GetSkillNum(chara, skillName, "init")
				if initNum == "-1" {
					returnMes = "Skill not found."
				} else {
					startNum = GetSkillNum(chara, skillName, "sum")
					nowNum = GetSkillNum(chara, skillName, "now")
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
		handlerResult.Normal.Content = "[" + skillName + "] Init( " + initNum + " ), Start( " + startNum + "), Now( " + nowNum + " )"
	}

	/* Embedメッセージ */
	if returnMes != "" {
		handlerResult.Error = errors.New(returnMes)
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Description: returnMes,
			Color:       core.EnColorRed,
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
			Title:  "< " + skillName + " >",
			Color:  core.EnColorGreen,
			Fields: fields,
		}
	}

	return handlerResult
}

// 能力値操作ハンドラ
func CmdCharaNumControl(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var targetSkill string
	var oldNum string
	var newNum string
	var diffCmd string
	var rollResultMessage string
	var returnMes string

	if len(md.Options) < 2 {
		returnMes = "Invalid arguments."
	} else {
		if cs == nil {
			returnMes = "Character not registered."
		} else {
			var chara *CharacterOfCthulhu
			var exist bool
			var ctrlNum string

			for _, opt := range md.Options {
				if opt.Name == "0" || opt.Name == "target" {
					targetSkill = opt.Value
				} else {
					ctrlNum = opt.Value
				}
			}

			if core.GetParentIDFromChildID(md.ChannelID) != "" {
				chara, exist = (*cs).Npc[md.AuthorID].(*CharacterOfCthulhu)
			} else {
				chara, exist = (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
			}
			if exist == false {
				returnMes = "Character not found."
				handlerResult.Error = errors.New(returnMes)
			} else {
				oldNum = GetSkillNum(chara, targetSkill, "now")
				if oldNum == "-1" {
					returnMes = "Skill not found."
					handlerResult.Error = errors.New(returnMes)
				} else {
					diffRegex := regexp.MustCompile("^[+-]?[0-9]+$")
					diffCmd = ctrlNum
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
					newNum = AddSkillNum(chara, targetSkill, diffCmd)
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
		handlerResult.Normal.Content = "[" + targetSkill + "] " + oldNum + " => " + newNum + " (Diff: " + diffCmd + ")"
	}

	/* Embedメッセージ */
	if returnMes != "" {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Description: returnMes,
			Color:       core.EnColorRed,
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
			Title:       "< " + targetSkill + " >",
			Description: rollResultMessage,
			Color:       core.EnColorGreen,
			Fields:      fields,
		}
	}

	return handlerResult
}

// キャラクターシート連携ダイスロール共通処理
func jobLinkRoll(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var rollResult core.BCDiceRollResult
	var command string
	var err error
	var returnMes string

	if len(md.Options) == 0 {
		returnMes = "Invalid arguments."
		handlerResult.Error = errors.New(returnMes)
	} else {
		command = md.Options[0].Value
		if cs == nil {
			returnMes = "Character not registered."
			handlerResult.Error = errors.New(returnMes)
		} else {
			var characterData interface{}
			if cs.Chat.Parent.ID == md.ChannelID {
				characterData = (*cs).Pc[md.AuthorID]
			} else {
				characterData = (*cs).Npc[md.AuthorID]
			}
			if chara, exist := characterData.(*CharacterOfCthulhu); exist == false {
				returnMes = "Character not found."
				handlerResult.Error = errors.New(returnMes)
			} else {
				diceCmd := "CCB<=" + command
				exRegex := regexp.MustCompile("[^\\+\\-\\*\\/ 　]+")
				ignoreRegex := regexp.MustCompile("^[0-9]+$")
				for _, ex := range exRegex.FindAllString(command, -1) {
					if ignoreRegex.MatchString(ex) == false {
						exNum := GetSkillNum(chara, ex, "now")
						if exNum == "-1" {
							returnMes = "Skill not found."
							handlerResult.Error = errors.New(returnMes)
						} else {
							diceCmd = strings.Replace(diceCmd, ex, exNum, -1)
							rollResult, err = core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, diceCmd)
							if err != nil {
								returnMes = "Server internal error."
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
			Color:       core.EnColorRed,
		}
	} else {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Title:       "< " + command + " >",
			Description: rollResult.Result,
			Color:       core.EnColorGreen,
		}
	}
	return handlerResult
}

// キャラクターシート連携ダイスロールハンドラ
func CmdLinkRoll(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	handlerResult = jobLinkRoll(cs, md)
	return handlerResult
}

// キャラクターシート連携シークレットダイスロールハンドラ
func CmdSecretLinkRoll(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	handlerResult = jobLinkRoll(cs, md)
	/* 有効にするメッセージタイプ */
	handlerResult.Secret.EnableType = core.EnEmbed
	/* テキストメッセージ */
	handlerResult.Secret.Content = "**SECRET DICE**"
	/* Embedメッセージ */
	handlerResult.Secret.Embed = &discordgo.MessageEmbed{
		Title: "SECRET DICE",
		Color: core.EnColorYellow,
	}
	return handlerResult
}

// CmdSecretDiceRoll シークレットダイスロールハンドラ
func CmdSecretDiceRoll(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var returnMes string = ""
	var command string = md.Options[0].Value
	rollResult, err := core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, command)
	if err != nil {
		returnMes = "Server internal error."
		handlerResult.Error = err
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
			Color:       core.EnColorRed,
		}
	} else {
		handlerResult.Normal.Embed = &discordgo.MessageEmbed{
			Title:       "< " + command + " >",
			Description: rollResult.Result,
			Color:       core.EnColorGreen,
		}
	}

	/* 有効にするメッセージタイプ */
	handlerResult.Secret.EnableType = core.EnEmbed
	/* テキストメッセージ */
	handlerResult.Secret.Content = "**SECRET DICE**"
	/* Embedメッセージ */
	handlerResult.Secret.Embed = &discordgo.MessageEmbed{
		Title: "SECRET DICE",
		Color: core.EnColorYellow,
	}
	return handlerResult
}

// SAN値チェック処理ハンドラ
func CmdSanCheckRoll(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
	var successCommand string
	var failedCommand string
	var sanRollResult core.BCDiceRollResult
	var sucRollResult core.BCDiceRollResult
	var failRollResult core.BCDiceRollResult
	var err error
	var orgSanNum string
	var sanSub string
	var newNum string
	var returnMes string

	if len(md.Options) < 2 {
		returnMes = "Invalid arguments."
	} else {

		for _, opt := range md.Options {
			if opt.Name == "0" || opt.Name == "success" {
				successCommand = opt.Value
			} else {
				failedCommand = opt.Value
			}
		}

		if cs == nil {
			returnMes = "PC not registered."
			handlerResult.Error = errors.New(returnMes)
		} else {
			pc, exist := (*cs).Pc[md.AuthorID].(*CharacterOfCthulhu)
			if exist == false {
				returnMes = "PC not found."
				handlerResult.Error = errors.New(returnMes)
			} else {
				orgSanNum = GetSkillNum(pc, "san", "now")
				sanRollCmd := "SCCB<=" + orgSanNum
				sanRollResult, err = core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, sanRollCmd)
				if err != nil {
					returnMes = "Server error."
					handlerResult.Error = err
				} else {
					if strings.Contains(sanRollResult.Result, "成功") || strings.Contains(sanRollResult.Result, "スペシャル") {
						if strings.Contains(successCommand, "d") {
							sucRollResult, err = core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, successCommand)
							if err != nil {
								returnMes = "Server error."
								handlerResult.Error = err
							} else {
								sanSub = "-" + core.CalcDicesSum(sucRollResult.Dices)
							}

						} else {
							sanSub = "-" + successCommand
						}
						newNum = AddSkillNum(pc, "san", sanSub)
					} else {
						if strings.Contains(failedCommand, "d") {
							failRollResult, err = core.ExecuteDiceRollAndCalc(core.GetConfig().EndPoint, (*cs).Scenario.System, failedCommand)
							if err != nil {
								returnMes = "Server error."
								handlerResult.Error = err
							} else {
								sanSub = "-" + core.CalcDicesSum(failRollResult.Dices)
							}
						} else {
							sanSub = "-" + failedCommand
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
			Color:       core.EnColorRed,
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
			Color:       core.EnColorGreen,
			Fields:      fields,
		}
	}
	return handlerResult
}

// ダイスロール統計表示処理
//func CmdShowStatistics(cs *core.Session, md core.MessageData) (handlerResult core.HandlerResult) {
//
//}
