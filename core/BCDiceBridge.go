package core

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// BCDiceVersionResult BCDiceバージョン情報格納構造体
type BCDiceVersionResult struct {
	API    string `json:"api"`
	BCDice string `json:"bcdice"`
}

// ExecuteVersionCheck BCDiceバージョン情報取得処理
func ExecuteVersionCheck(endpoint string) (vr BCDiceVersionResult, err error) {
	urlStr := endpoint + "/version"
	resp, err := http.Get(urlStr)
	log.Printf("\"[URL]: %s\"", urlStr)
	if err != nil {
		log.Println(err)
		return vr, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
		log.Println(err)
		return vr, err
	}
	return vr, nil
}

// BCDiceSystemsResult BCDiceシステム一覧取得結果格納構造体
type BCDiceSystemsResult struct {
	Systems []BCDiceSystem `json:"game_system"`
}

// BCDiceSystem BCDiceシステム情報格納構造体
type BCDiceSystem struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	SortKey string `json:"sort_key"`
}

// ExecuteGetSystems BCDiceシステム一覧取得
func ExecuteGetSystems(endpoint string) (sr BCDiceSystemsResult, err error) {
	urlStr := endpoint + "/game_system"
	resp, err := http.Get(urlStr)
	log.Printf("\"[URL]: %s\"", urlStr)
	if err != nil {
		log.Println(err)
		return sr, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		log.Println(err)
		return sr, err
	}
	return sr, nil
}

// CheckContainsSystem BCDiceシステム対応有無チェック
func CheckContainsSystem(endpoint string, system string) (result bool) {
	systemsList, err := ExecuteGetSystems(endpoint)
	if err == nil {
		for _, sys := range systemsList.Systems {
			if strings.ToLower(sys.Id) == strings.ToLower(system) {
				return true
			}
		}
	}

	return false
}

// BCDiceRollResult ダイスロール実行結果格納構造体
type BCDiceRollResult struct {
	Ok       bool   `json:"ok"`
	Result   string `json:"text"`
	Secret   bool   `json:"secret"`
	Success  bool   `json:"success"`
	Failure  bool   `json:"failure"`
	Critical bool   `json:"critical"`
	Fumble   bool   `json:"fumble"`
	Dices    []Dice `json:"rands"`
}

// Dice ダイス情報格納構造体
type Dice struct {
	Kind  string `json:"kind"`
	Faces int    `json:"sides"`
	Value int    `json:"value"`
}

// ExecuteDiceRollAndCalc ダイスロール+演算実行
func ExecuteDiceRollAndCalc(endpoint string, system string, dice string) (rr BCDiceRollResult, err error) {
	/* 不等号を堺に文字列分割 */
	rep := regexp.MustCompile("(.+)(<=|>=|=|>|<)(.+)")
	splitTargetStr := rep.ReplaceAllString(dice, "$1$2{SPLIT}$3")
	splitedStrArray := strings.Split(splitTargetStr, "{SPLIT}")

	var diceCalcStr string
	if len(splitedStrArray) > 1 {
		diceCalcStr = splitedStrArray[1]
	} else {
		diceCalcStr = splitedStrArray[0]
	}

	/* 計算式が含まれているか確認 */
	calcCheckRegp := regexp.MustCompile("[\\( ]*[\\+\\- ]*[a-zA-Z0-9 ]+[\\) ]*[\\+\\-\\/\\* ]{1}[\\( ]*[\\+\\- ]*[a-zA-Z0-9 ]+[\\) ]*")
	isCalcMutch := calcCheckRegp.MatchString(diceCalcStr)
	if isCalcMutch {
		/* 計算式が含まれていた場合 */
		var calAnswer string
		var workingFormula string
		calAnswer, workingFormula, err = CalcStr2Ans(diceCalcStr, system)
		if err != nil {
			return rr, err
		}
		if len(splitedStrArray) > 1 {
			var rrtmp BCDiceRollResult
			strIntegDiceCmd := splitedStrArray[0] + calAnswer
			rrtmp, err = ExecuteDiceRoll(endpoint, system, strIntegDiceCmd)
			rr.Ok = rrtmp.Ok
			if "" != workingFormula {
				rr.Result = "calc(" + dice + ") \n＞ calc(" + workingFormula + ") \n＞ " + strings.Replace(rrtmp.Result, "＞", "\n＞", -1)
			} else {
				rr.Result = "calc(" + dice + ") \n＞ " + strings.Replace(rrtmp.Result, "＞", "\n＞", -1)
			}
			rr.Secret = rrtmp.Secret
			rr.Success = rrtmp.Success
			rr.Failure = rrtmp.Failure
			rr.Critical = rrtmp.Critical
			rr.Fumble = rrtmp.Fumble
			rr.Dices = rrtmp.Dices
		} else {
			rr.Ok = true
			if "" != workingFormula {
				rr.Result = "calc(" + dice + ") \n＞ calc(" + workingFormula + ") \n＞ " + strings.Replace(calAnswer, "＞", "\n＞", -1)
			} else {
				rr.Result = "calc(" + dice + ") \n＞ " + strings.Replace(calAnswer, "＞", "\n＞", -1)
			}
			rr.Secret = false
			rr.Success = false
			rr.Failure = false
			rr.Critical = false
			rr.Fumble = false
			rr.Dices = make([]Dice, 1)
			rr.Dices[0].Value, _ = strconv.Atoi(calAnswer)
		}

	} else {
		/* 計算式が含まれていなかった場合 */
		rr, _ = ExecuteDiceRoll(endpoint, system, dice)
		rr.Result = strings.Replace(rr.Result, "＞", "\n＞", -1)
	}

	return rr, err
}

// ExecuteDiceRoll ダイスロール実行
func ExecuteDiceRoll(endpoint string, system string, dice string) (rr BCDiceRollResult, err error) {
	urlStr := endpoint + "/game_system/" + system + "/roll?command=" + url.QueryEscape(dice)
	resp, err := http.Get(urlStr)
	log.Printf("\"[URL]: %s\"", urlStr)
	if err != nil {
		log.Println(err)
		return rr, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		log.Println(err)
		return rr, err
	}
	log.Println(rr)
	return rr, nil
}

// CalcDicesSum ダイス合計値算出
func CalcDicesSum(dices []Dice) string {
	var diceSum int
	for _, d := range dices {
		diceSum += d.Value
	}
	return strconv.Itoa(diceSum)
}
