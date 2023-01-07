package cthulhu

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// キャラクター能力名リスト
var cdAbilityNameList = []string{"STR", "CON", "POW", "DEX", "APP", "SIZ", "INT", "EDU", "HP", "MP", "SAN", "アイデア", "幸運", "知識"}

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/

// キャラクターシート情報取得処理
func GetCharSheetFromURL(url string) (*CharaSheet, error) {
	var cs CharaSheet

	resp, err := http.Get(url + ".json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&cs)
	if err != nil {
		return nil, err
	}

	return &cs, err
}

// キャラクターデータ生成処理
func GetCharDataFromCharSheet(cs *CharaSheet, plName string, plID string) *CharacterOfCthulhu {
	var cd CharacterOfCthulhu

	cd.Player.Name = plName
	cd.Player.ID = plID

	cd.Md = (*cs).Phrase
	cd.ID = strconv.Itoa((*cs).DataID)

	// パーソナルデータ
	cd.Personal.Name = (*cs).PcName
	cd.Personal.Job = (*cs).Shuzoku
	cd.Personal.Age, _ = strconv.Atoi((*cs).Age)
	cd.Personal.Sex = (*cs).Sex
	cd.Personal.Height = (*cs).PcHeight
	cd.Personal.Weight = (*cs).PcWeight
	cd.Personal.Country = (*cs).PcKigen
	cd.Personal.Haircl = (*cs).ColorHair
	cd.Personal.Eyecl = (*cs).ColorEye
	cd.Personal.Skincl = (*cs).ColorSkin
	cd.Memo = (*cs).PcMakingMemo

	// 能力値の登録
	cd.Ability = map[string]*Ability{}

	csAbilityPointerList := [][]*string{
		{&cs.StrInit, &cs.StrAdd, &cs.StrTemp, &cs.StrNow, &cs.StrNow},
		{&cs.ConInit, &cs.ConAdd, &cs.ConTemp, &cs.ConNow, &cs.ConNow},
		{&cs.PowInit, &cs.PowAdd, &cs.PowTemp, &cs.PowNow, &cs.PowNow},
		{&cs.DexInit, &cs.DexAdd, &cs.DexTemp, &cs.DexNow, &cs.DexNow},
		{&cs.AppInit, &cs.AppAdd, &cs.AppTemp, &cs.AppNow, &cs.AppNow},
		{&cs.SizInit, &cs.SizAdd, &cs.SizTemp, &cs.SizNow, &cs.SizNow},
		{&cs.IntInit, &cs.IntAdd, &cs.IntTemp, &cs.IntNow, &cs.IntNow},
		{&cs.EduInit, &cs.EduAdd, &cs.EduTemp, &cs.EduNow, &cs.EduNow},
		{&cs.HpInit, &cs.HpAdd, &cs.HpTemp, &cs.HpNow, &cs.HpNow},
		{&cs.MpInit, &cs.MpAdd, &cs.MpTemp, &cs.MpNow, &cs.MpNow},
		{&cs.SanInit, &cs.SanAdd, &cs.SanTemp, &cs.SanLeft, &cs.SanLeft},
		{&cs.IdeaInit, &cs.IdeaAdd, &cs.IdeaTemp, &cs.IdeaNow, &cs.IdeaNow},
		{&cs.LuckyInit, &cs.LuckyAdd, &cs.LuckyTemp, &cs.LuckyNow, &cs.LuckyNow},
		{&cs.KnowledgeInit, &cs.KnowledgeAdd, &cs.KnowledgeTemp, &cs.KnowledgeNow, &cs.KnowledgeNow},
	}

	for n, cdan := range cdAbilityNameList {
		var ta Ability
		ta.Name = cdan
		ta.Init, _ = strconv.Atoi(*csAbilityPointerList[n][0])
		ta.Add, _ = strconv.Atoi(*csAbilityPointerList[n][1])
		ta.Temp, _ = strconv.Atoi(*csAbilityPointerList[n][2])
		ta.Sum, _ = strconv.Atoi(*csAbilityPointerList[n][3])
		ta.Now, _ = strconv.Atoi(*csAbilityPointerList[n][4])
		cd.Ability[cdan] = &ta
	}

	// 技能の登録
	cd.Skill = map[string]*Skill{}
	atkSkillList := append([]string{"回避", "キック", "組み付き", "こぶし（パンチ）", "頭突き", "投擲", "マーシャルアーツ", "拳銃", "サブマシンガン", "ショットガン", "マシンガン", "ライフル"}, (*cs).SkillAtkExName...)
	searchSkillList := append([]string{"応急手当", "鍵開け", "隠す", "隠れる", "聞き耳", "忍び歩き", "写真術", "精神分析", "追跡", "登攀", "図書館", "目星"}, (*cs).SkillSearchExName...)
	actSkillList := append([]string{"運転", "機械修理", "重機械操作", "乗馬", "水泳", "製作", "操縦", "跳躍", "電気修理", "ナビゲート", "変装"}, (*cs).SkillActExName...)
	negSkillList := append([]string{"言いくるめ", "信用", "説得", "根切り", "母国語"}, (*cs).SkillNegExName...)
	knowSkillList := append([]string{"医学", "オカルト", "化学", "クトゥルフ神話", "芸術", "経理", "考古学", "コンピューター", "心理学", "人類学", "生物学", "地質学", "電子工学", "天文学", "博物学", "物理学", "法律", "薬学", "歴史"}, (*cs).SkillKnowExName...)
	allSkillList := []*[]string{&atkSkillList, &searchSkillList, &actSkillList, &negSkillList, &knowSkillList}

	csSkillPointerList := [][]*[]string{
		{&cs.SkillAtkInit, &cs.SkillAtkJob, &cs.SkillAtkHobby, &cs.SkillAtkGrowth,
			&cs.SkillAtkOther, &cs.SkillAtkSum, &cs.SkillAtkSum,
		},
		{&cs.SkillSearchSum, &cs.SkillSearchInit, &cs.SkillSearchJob, &cs.SkillSearchHobby,
			&cs.SkillSearchGrowth, &cs.SkillSearchOther, &cs.SkillSearchSum,
		},
		{&cs.SkillActSum, &cs.SkillActInit, &cs.SkillActJob, &cs.SkillActHobby,
			&cs.SkillActGrowth, &cs.SkillActOther, &cs.SkillActSum,
		},
		{&cs.SkillNegSum, &cs.SkillNegInit, &cs.SkillNegJob, &cs.SkillNegHobby,
			&cs.SkillNegGrowth, &cs.SkillNegOther, &cs.SkillNegSum,
		},
		{&cs.SkillKnowSum, &cs.SkillKnowInit, &cs.SkillKnowJob, &cs.SkillKnowHobby,
			&cs.SkillKnowGrowth, &cs.SkillKnowOther, &cs.SkillKnowSum,
		},
	}
	for asi, slist := range allSkillList {
		for uni, name := range *slist {
			var us Skill
			switch name {
			case "運転":
				us.Sub = (*cs).SkillActDriveType
				break
			case "製作":
				us.Sub = (*cs).SkillActMakeType
				break
			case "操縦":
				us.Sub = (*cs).SkillActOpeType
				break
			case "母国語":
				us.Sub = (*cs).SkillNegMylang
				break
			case "芸術":
				us.Sub = (*cs).SkillKnowArtType
				break
			default:
				us.Sub = ""
				break
			}
			us.Name = name
			us.Init, _ = strconv.Atoi((*csSkillPointerList[asi][0])[uni])
			us.Job, _ = strconv.Atoi((*csSkillPointerList[asi][1])[uni])
			us.Hobby, _ = strconv.Atoi((*csSkillPointerList[asi][2])[uni])
			us.Growth, _ = strconv.Atoi((*csSkillPointerList[asi][3])[uni])
			us.Other, _ = strconv.Atoi((*csSkillPointerList[asi][4])[uni])
			us.Sum, _ = strconv.Atoi((*csSkillPointerList[asi][5])[uni])
			us.Now, _ = strconv.Atoi((*csSkillPointerList[asi][6])[uni])
			us.Growflg = 0
			cd.Skill[name] = &us
		}
	}

	// 戦闘・武器・防具

	// 所持品・所持金

	return &cd
}

// キャラクターデータの能力リストを公開する
func GetCdAbilityNameList() []string {
	return cdAbilityNameList
}
