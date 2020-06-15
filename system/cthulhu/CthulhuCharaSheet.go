package cthulhu

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// CharaSheet キャラシート情報格納構造体
type CharaSheet struct {
	StrInit       string `json:"NA1"`
	ConInit       string `json:"NA2"`
	PowInit       string `json:"NA3"`
	DexInit       string `json:"NA4"`
	AppInit       string `json:"NA5"`
	SizInit       string `json:"NA6"`
	IntInit       string `json:"NA7"`
	EduInit       string `json:"NA8"`
	HpInit        string `json:"NA9"`
	MpInit        string `json:"NA10"`
	SanInit       string `json:"NA11"`
	IdeaInit      string `json:"NA12"`
	LuckyInit     string `json:"NA13"`
	KnowledgeInit string `json:"NA14"`

	StrAdd       string `json:"NS1"`
	ConAdd       string `json:"NS2"`
	PowAdd       string `json:"NS3"`
	DexAdd       string `json:"NS4"`
	AppAdd       string `json:"NS5"`
	SizAdd       string `json:"NS6"`
	IntAdd       string `json:"NS7"`
	EduAdd       string `json:"NS8"`
	HpAdd        string `json:"NS9"`
	MpAdd        string `json:"NS10"`
	SanAdd       string `json:"NS11"`
	IdeaAdd      string `json:"NS12"`
	LuckyAdd     string `json:"NS13"`
	KnowledgeAdd string `json:"NS14"`

	StrTemp       string `json:"NM1"`
	ConTemp       string `json:"NM2"`
	PowTemp       string `json:"NM3"`
	DexTemp       string `json:"NM4"`
	AppTemp       string `json:"NM5"`
	SizTemp       string `json:"NM6"`
	IntTemp       string `json:"NM7"`
	EduTemp       string `json:"NM8"`
	HpTemp        string `json:"NM9"`
	MpTemp        string `json:"NM10"`
	SanTemp       string `json:"NM11"`
	IdeaTemp      string `json:"NM12"`
	LuckyTemp     string `json:"NM13"`
	KnowledgeTemp string `json:"NM14"`

	StrNow       string `json:"NP1"`
	ConNow       string `json:"NP2"`
	PowNow       string `json:"NP3"`
	DexNow       string `json:"NP4"`
	AppNow       string `json:"NP5"`
	SizNow       string `json:"NP6"`
	IntNow       string `json:"NP7"`
	EduNow       string `json:"NP8"`
	HpNow        string `json:"NP9"`
	MpNow        string `json:"NP10"`
	SanNow       string `json:"NP11"`
	IdeaNow      string `json:"NP12"`
	LuckyNow     string `json:"NP13"`
	KnowledgeNow string `json:"NP14"`

	SanLeft   string `json:"SAN_Left"`
	SanMax    string `json:"SAN_Max"`
	SanDanger string `json:"SAN_Danger"`

	TSTotal   string `json:"TS_Total"`
	TSMaximum string `json:"TS_Maximum"`
	TSAdd     string `json:"TS_Add"`
	TKTotal   string `json:"TK_Total"`
	TKMaximum string `json:"TK_Maximum"`
	TKAdd     string `json:"TK_Add"`

	SkillAtkGrowflg []string `json:"TBAU"`
	SkillAtkInit    []string `json:"TBAD"`
	SkillAtkJob     []string `json:"TBAS"`
	SkillAtkHobby   []string `json:"TBAK"`
	SkillAtkGrowth  []string `json:"TBAA"`
	SkillAtkOther   []string `json:"TBAO"`
	SkillAtkSum     []string `json:"TBAP"`
	SkillAtkExName  []string `json:"TBAName"`

	SkillSearchGrowflg []string `json:"TFAU"`
	SkillSearchInit    []string `json:"TFAD"`
	SkillSearchJob     []string `json:"TFAS"`
	SkillSearchHobby   []string `json:"TFAK"`
	SkillSearchGrowth  []string `json:"TFAA"`
	SkillSearchOther   []string `json:"TFAO"`
	SkillSearchSum     []string `json:"TFAP"`
	SkillSearchExName  []string `json:"TFAName"`

	SkillActGrowflg   []string `json:"TAAU"`
	SkillActDriveType string   `json:"unten_bunya"`
	SkillActInit      []string `json:"TAAD"`
	SkillActJob       []string `json:"TAAS"`
	SkillActHobby     []string `json:"TAAK"`
	SkillActGrowth    []string `json:"TAAA"`
	SkillActOther     []string `json:"TAAO"`
	SkillActSum       []string `json:"TAAP"`
	SkillActMakeType  string   `json:"seisaku_bunya"`
	SkillActOpeType   string   `json:"main_souju_norimono"`
	SkillActExName    []string `json:"TAAName"`

	SkillNegGrowflg []string `json:"TCAU"`
	SkillNegInit    []string `json:"TCAD"`
	SkillNegJob     []string `json:"TCAS"`
	SkillNegHobby   []string `json:"TCAK"`
	SkillNegGrowth  []string `json:"TCAA"`
	SkillNegOther   []string `json:"TCAO"`
	SkillNegSum     []string `json:"TCAP"`
	SkillNegMylang  string   `json:"mylang_name"`
	SkillNegExName  []string `json:"TCAName"`

	SkillKnowGrowflg []string `json:"TKAU"`
	SkillKnowInit    []string `json:"TKAD"`
	SkillKnowJob     []string `json:"TKAS"`
	SkillKnowHobby   []string `json:"TKAK"`
	SkillKnowGrowth  []string `json:"TKAA"`
	SkillKnowOther   []string `json:"TKAO"`
	SkillKnowSum     []string `json:"TKAP"`
	SkillKnowArtType string   `json:"geijutu_bunya"`
	SkillKnowExName  []string `json:"TKAName"`

	DmgBonus        string   `json:"dmg_bonus"`
	ArmsName        []string `json:"arms_name"`
	ArmsHit         []string `json:"arms_hit"`
	ArmsDamage      []string `json:"arms_damage"`
	ArmsRange       []string `json:"arms_range"`
	ArmsAttackCount []string `json:"arms_attack_count"`
	ArmsLastShot    []string `json:"arms_last_shot"`
	ArmsVitality    []string `json:"arms_vitality"`
	ArmsSonota      []string `json:"arms_sonota"`

	ItemName     []string `json:"item_name"`
	ItemTanka    []string `json:"item_tanka"`
	ItemNum      []string `json:"item_num"`
	ItemPrice    []string `json:"item_price"`
	ItemMemo     []string `json:"item_memo"`
	PriceItemSum string   `json:"price_item_sum"`
	Money        string   `json:"money"`
	Debt         string   `json:"debt"`

	PcName       string `json:"pc_name"`
	PcTags       string `json:"pc_tags"`
	Shuzoku      string `json:"shuzoku"`
	Age          string `json:"age"`
	Sex          string `json:"sex"`
	PcHeight     string `json:"pc_height"`
	PcWeight     string `json:"pc_weight"`
	PcKigen      string `json:"pc_kigen"`
	ColorHair    string `json:"color_hair"`
	ColorEye     string `json:"color_eye"`
	ColorSkin    string `json:"color_skin"`
	PcMakingMemo string `json:"pc_making_memo"`
	Message      string `json:"message"`
	Game         string `json:"game"`
	DataID       int    `json:"data_id"`
	Phrase       string `json:"phrase"`
}

// CdAbilityNameList 能力名リスト
var CdAbilityNameList = []string{"STR", "CON", "POW", "DEX", "APP", "SIZ", "INT", "EDU", "HP", "MP", "SAN", "アイデア", "幸運", "知識"}

// GetCharSheetFromURL キャラシート情報取得処理
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

// GetCharDataFromCharSheet キャラデータ生成処理
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

	for n, cdan := range CdAbilityNameList {
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
	actSkillList := append([]string{"運転", "機械修理", "重機械操作", "乗馬", "水泳", "製作", "操縦", "跳躍", "電気修理", "ナビゲート", "返送"}, (*cs).SkillActExName...)
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
