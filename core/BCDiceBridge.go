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
	resp, err := http.Get(endpoint + "/version")
	log.Printf("\"%s\"", endpoint)
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
	Systems []string `json:"systems"`
}

// ExecuteGetSystems BCDiceシステム一覧取得
func ExecuteGetSystems(endpoint string) (sr BCDiceSystemsResult, err error) {
	resp, err := http.Get(endpoint + "/systems")
	log.Printf("\"%s\"", endpoint)
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
			if strings.ToLower(sys) == strings.ToLower(system) {
				return true
			}
		}
	}

	return false
}

// BCDiceRollResult ダイスロール実行結果格納構造体
type BCDiceRollResult struct {
	Ok     bool   `json:"ok"`
	Result string `json:"result"`
	Secret bool   `json:"secret"`
	Dices  []Dice `json:"dices"`
}

// Dice ダイス情報格納構造体
type Dice struct {
	Faces int `json:"faces"`
	Value int `json:"value"`
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
	calcCheckRegp := regexp.MustCompile("[\\+-/\\*\\(\\)]")
	isCalcMutch := calcCheckRegp.MatchString(diceCalcStr)
	if isCalcMutch {
		/* 計算式が含まれていた場合 */
		var calAnswer string
		calAnswer, err = CalcStr2Ans(diceCalcStr, system)
		if err != nil {
			return rr, err
		}
		if len(splitedStrArray) > 1 {
			var rrtmp BCDiceRollResult
			strIntegDiceCmd := splitedStrArray[0] + calAnswer
			rrtmp, err = ExecuteDiceRoll(endpoint, system, strIntegDiceCmd)
			rr.Ok = rrtmp.Ok
			rr.Result = "calc(" + dice + ") ＞ " + rrtmp.Result
			rr.Secret = rrtmp.Secret
			rr.Dices = rrtmp.Dices
		} else {
			rr.Ok = true
			rr.Result = "calc(" + dice + ") ＞ " + calAnswer
			rr.Secret = false
			rr.Dices = make([]Dice, 1)
			rr.Dices[0].Value, _ = strconv.Atoi(calAnswer)
		}

	} else {
		/* 計算式が含まれていなかった場合 */
		rr, _ = ExecuteDiceRoll(endpoint, system, dice)
	}

	return rr, err
}

// ExecuteDiceRoll ダイスロール実行
func ExecuteDiceRoll(endpoint string, system string, dice string) (rr BCDiceRollResult, err error) {
	resp, err := http.Get(endpoint + "/diceroll?system=" + system + "&command=" + url.QueryEscape(dice))
	log.Printf("\"%s\"", endpoint)
	if err != nil {
		log.Println(err)
		return rr, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		log.Println(err)
		return rr, err
	}
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
