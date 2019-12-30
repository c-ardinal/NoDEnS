package nodens

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
