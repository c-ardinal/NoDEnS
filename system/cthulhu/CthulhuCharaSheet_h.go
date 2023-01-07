package cthulhu

/****************************************************************************/
/* 外部公開型定義                                                           */
/****************************************************************************/

// キャラクターシート情報格納構造体
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

/****************************************************************************/
/* 外部公開定数定義                                                         */
/****************************************************************************/
