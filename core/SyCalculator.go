package core

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/stack"
)

/****************************************************************************/
/* 内部型定義                                                               */
/****************************************************************************/

// トークン型
type tokenTypeT int8

// トークン情報構造体
type tokenT struct {
	token     string
	tokentype tokenTypeT
	pin       int8
	pst       int8
}

/****************************************************************************/
/* 内部定数定義                                                             */
/****************************************************************************/

// トークン型ダイナミクス定義
const (
	ESCAPE tokenTypeT = iota
	NUMBER
	COMMAND
	OPERATOR
	LEFTPAREN
	RIGHTPAREN
)

/****************************************************************************/
/* 内部変数定義                                                             */
/****************************************************************************/

// トークン情報定義
var tokensDict = map[string]tokenT{
	"+":   {"+", OPERATOR, 4, 5},
	"-":   {"-", OPERATOR, 4, 5},
	"/":   {"/", OPERATOR, 6, 7},
	"*":   {"*", OPERATOR, 6, 7},
	"(":   {"(", LEFTPAREN, 8, 1},
	")":   {")", RIGHTPAREN, 1, 0},
	"NUM": {"NUM", NUMBER, 9, 9},
	"CMD": {"CMD", COMMAND, 9, 9},
	"$$$": {"$$$", ESCAPE, -1, -1},
}

/****************************************************************************/
/* 関数定義                                                                 */
/****************************************************************************/

// 文字列として受け取った計算式を計算し，計算結果を文字列として返す
func CalcStr2Ans(s string, system string) (result string, numOnlyFormula string, err error) {
	var numOnlyTokens []tokenT
	err = nil
	tokens := convStr2Tokens(s)
	isError, errorCol, errorMes, isContCmd := evalTokens(tokens)
	if isError {
		err = errors.New("Syntax error [ " + strconv.Itoa(errorCol) + ", " + errorMes + " ]")
	} else {
		if isContCmd {
			numOnlyTokens, err = convDiceTokens2NumTokens(tokens, system)
			if err != nil {
				return "0", "0", err
			}
			for _, s := range numOnlyTokens {
				numOnlyFormula += s.token
			}
		} else {
			numOnlyTokens = tokens
		}
		syConvedTokens := convTokens2ShuntingYardTokens(numOnlyTokens)
		result, err = calFromTokens(syConvedTokens)
	}
	return result, numOnlyFormula, err
}

// 文字列をトークン列へ変換する
func convStr2Tokens(str string) (result []tokenT) {
	/* 文字列をトークン列に変換するために整理 */
	strTrimSpaces := strings.Trim(str, " ")
	strTrimSpaces = strings.Trim(strTrimSpaces, "　")

	var strInserSpaces string = strTrimSpaces
	for _, dict := range tokensDict {
		strInserSpaces = strings.Replace(strInserSpaces, dict.token, " "+dict.token+" ", -1)
	}

	strBase := strings.TrimSpace(strInserSpaces)

	splitRegp := regexp.MustCompile(" +")
	strArray := splitRegp.Split(strBase, -1)

	/* 文字配列をトークン化 */
	for _, t := range strArray {
		tmpToken, isExist := tokensDict[t]
		if !isExist {
			numOrCmdRegp := regexp.MustCompile("^[0-9]+$")
			isNum := numOrCmdRegp.MatchString(t)
			if isNum {
				tmpToken = tokensDict["NUM"]
			} else {
				tmpToken = tokensDict["CMD"]
			}
			tmpToken.token = t
		}
		result = append(result, tmpToken)
	}
	return result
}

// トークン列を評価し、構文誤りが無いかチェックする
func evalTokens(tknArray []tokenT) (result bool, errorCol int, errorMes string, isContCmd bool) {
	var parenPairNum int = 0

	errorCol = -1
	isContCmd = false
	errorMes = ""

	for i, t := range tknArray {
		switch t.tokentype {
		case COMMAND:
			isContCmd = true
			fallthrough
		case NUMBER:
			if i == 0 {
				//.先頭が数字かコマンドの場合は正常
				result = false
			} else if tknArray[i-1].tokentype == OPERATOR || tknArray[i-1].tokentype == LEFTPAREN {
				// 記号か（の後でも正常
				result = false
			} else {
				result = true
			}
		case OPERATOR:
			if i == 0 {
				//.先頭が数字かコマンドの場合は異常
				result = true
			} else if tknArray[i-1].tokentype == RIGHTPAREN || tknArray[i-1].tokentype == NUMBER || tknArray[i-1].tokentype == COMMAND {
				// ）か数字かコマンドの後に有れば正常
				result = false
			} else {
				result = true
			}
		case LEFTPAREN:
			parenPairNum++
			if i == 0 {
				//.先頭が（の場合は正常
				result = false
			} else if tknArray[i-1].tokentype == OPERATOR || tknArray[i-1].tokentype == LEFTPAREN {
				// 記号か（の後でも正常
				result = false
			} else {
				result = true
			}
		case RIGHTPAREN:
			parenPairNum--
			if i == 0 {
				//.先頭が数字かコマンドの場合は異常
				result = true
			} else if tknArray[i-1].tokentype == RIGHTPAREN || tknArray[i-1].tokentype == NUMBER || tknArray[i-1].tokentype == COMMAND {
				// ）か数字かコマンドの後に有れば正常
				result = false
			} else {
				result = true
			}
		default:
			result = true
		}

		if result {
			errorCol = int(i)
			errorMes = t.token
			break
		}
	}

	if parenPairNum > 0 {
		result = true
		errorCol = int(len(tknArray))
		errorMes = "Missing ')'"
	} else if parenPairNum < 0 {
		result = true
		errorCol = int(len(tknArray))
		errorMes = "Missing '('"
	}

	return result, errorCol, errorMes, isContCmd
}

// トークン列をシャンティングヤード法に従って並び替える
func convTokens2ShuntingYardTokens(tknArray []tokenT) (result []tokenT) {
	var convedTokens []tokenT
	var stk = stack.New()

	escapeToken := tokensDict["$$$"]
	stk.Push(escapeToken)

	tknArray = append(tknArray, escapeToken)

	for _, t := range tknArray {
		convTokens(t, stk, &convedTokens)
	}

	for _, r := range convedTokens {
		if r.tokentype != LEFTPAREN && r.tokentype != RIGHTPAREN {
			result = append(result, r)
		}
	}

	return result
}

// トークンスタックの操作を行う
func convTokens(t tokenT, stk *stack.Stack, result *[]tokenT) {
	if t.pin > stk.Peek().(tokenT).pst {
		stk.Push(t)
	} else if t.pin < stk.Peek().(tokenT).pst {
		*result = append(*result, stk.Pop().(tokenT))
		convTokens(t, stk, result)
	} else {
		if t.tokentype != ESCAPE {
			*result = append(*result, stk.Pop().(tokenT))
			*result = append(*result, t)
		}
	}
}

// ダイスコマンドを含んだトークン列を、数字のみのトークン列へ変換する
func convDiceTokens2NumTokens(tknArray []tokenT, system string) (result []tokenT, err error) {
	var rollResult BCDiceRollResult
	result = tknArray
	for i, t := range tknArray {
		if t.tokentype == COMMAND {
			rollResult, err = ExecuteDiceRoll(GetConfig().EndPoint, system, t.token)
			var sum int = 0
			for _, d := range rollResult.Dices {
				sum += int(d.Value)
			}
			result[i].token = strconv.Itoa(sum)
			result[i].tokentype = NUMBER
		}
	}

	return result, err
}

// トークン列を元に、演算を実行する
func calFromTokens(tknArray []tokenT) (result string, err error) {
	var stk = stack.New()

	err = nil

	for _, t := range tknArray {
		i, _ := strconv.Atoi(t.token)
		if t.tokentype == NUMBER {
			stk.Push(i)
		} else if t.tokentype == OPERATOR {
			first := stk.Pop().(int)
			end := stk.Pop().(int)
			switch t.token {
			case "+":
				stk.Push(end + first)
			case "-":
				stk.Push(end - first)
			case "*":
				stk.Push(end * first)
			case "/":
				if first == 0 {
					err = errors.New("divide zero")
					break
				} else {
					stk.Push(end / first)
				}
			default:
			}
		}
	}

	if err == nil {
		ans := stk.Pop().(int)
		if ans < 0 {
			ans = 0
		}
		result = strconv.Itoa(ans)
	} else {
		result = "0"
	}

	return result, err
}
