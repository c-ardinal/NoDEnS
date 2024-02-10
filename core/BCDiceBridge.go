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

// BCDiceバージョン情報取得処理
func ExecuteVersionCheck(endpoint string) (vr BCDiceVersionResult, err error) {
	urlStr := endpoint + "/version"
	resp, err := http.Get(urlStr)
	log.Printf("[Event]: BCDice-API call > %s", urlStr)
	if err != nil {
		log.Printf("[Warning]: %v", err)
		return vr, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
		log.Printf("[Warning]: %v", err)
		return vr, err
	}
	return vr, nil
}

// BCDiceシステム一覧取得
func ExecuteGetSystems(endpoint string) (sr BCDiceSystemsResult, err error) {
	urlStr := endpoint + "/game_system"
	resp, err := http.Get(urlStr)
	log.Printf("[Event]: BCDice-API call > %s", urlStr)
	if err != nil {
		log.Printf("[Warning]: %v", err)
		return sr, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		log.Printf("[Warning]: %v", err)
		return sr, err
	}
	return sr, nil
}

// BCDiceシステム対応有無チェック
func CheckContainsSystem(endpoint string, system string) (result bool) {
	systemsList, err := ExecuteGetSystems(endpoint)
	if err == nil {
		for _, sys := range systemsList.Systems {
			if strings.EqualFold(strings.ToLower(sys.Id), strings.ToLower(system)) {
				return true
			}
		}
	}

	return false
}

// ダイスロール+演算実行
func ExecuteDiceRollAndCalc(endpoint string, system string, dice string) (rr BCDiceRollResult, err error) {
	/* 不等号を堺に文字列分割 */
	rep := regexp.MustCompile("(.+)(<=|>=|=|>|<)(.+)")
	splitTargetStr := rep.ReplaceAllString(dice, "$1$2{SPLIT}$3")
	splitStrArray := strings.Split(splitTargetStr, "{SPLIT}")

	var diceCalcStr string
	if len(splitStrArray) > 1 {
		diceCalcStr = splitStrArray[1]
	} else {
		diceCalcStr = splitStrArray[0]
	}

	/* 計算式が含まれているか確認 */
	calcCheckRegp := regexp.MustCompile(`[\( ]*[\+\- ]*[a-zA-Z0-9 ]+[\) ]*[\+\-\/\* ]{1}[\( ]*[\+\- ]*[a-zA-Z0-9 ]+[\) ]*`)
	isCalcMatch := calcCheckRegp.MatchString(diceCalcStr)
	if isCalcMatch {
		/* 計算式が含まれていた場合 */
		var calAnswer string
		var workingFormula string
		calAnswer, workingFormula, err = CalcStr2Ans(diceCalcStr, system)
		if err != nil {
			return rr, err
		}
		if len(splitStrArray) > 1 {
			var rrtmp BCDiceRollResult
			strIntegDiceCmd := splitStrArray[0] + calAnswer
			rrtmp, err = ExecuteDiceRoll(endpoint, system, strIntegDiceCmd)
			rr.Ok = rrtmp.Ok
			if workingFormula != "" {
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
			if workingFormula != "" {
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

// ダイスロール実行
func ExecuteDiceRoll(endpoint string, system string, dice string) (rr BCDiceRollResult, err error) {
	urlStr := endpoint + "/game_system/" + system + "/roll?command=" + url.QueryEscape(dice)
	resp, err := http.Get(urlStr)
	log.Printf("[Event]: BCDice-API call > %s", urlStr)
	if err != nil {
		log.Printf("[Warning]: %v", err)
		return rr, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		log.Printf("[Warning]: %v", err)
		return rr, err
	}
	log.Printf("[Event]: Dice roll result > '%v'", rr)
	return rr, nil
}

// ダイス合計値算出
func CalcDicesSum(dices []Dice) string {
	var diceSum int
	for _, d := range dices {
		diceSum += d.Value
	}
	return strconv.Itoa(diceSum)
}
