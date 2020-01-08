package protector

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
)

func GetSessionKey(rSource rand.Source) string {
	r := rand.New(rSource)
	result := ""
	for i := 0; i < 10; i++ {
		result += strconv.Itoa(r.Intn(10))
	}
	return result
}

func GetHashStr(rSource rand.Source) string {
	r := rand.New(rSource)
	result := ""
	for i := 0; i < 5; i++ {
		result += strconv.Itoa(r.Intn(6))
	}
	return result
}

type SessionProtector struct {
	HashString string
}

func (protector *SessionProtector) NextSessionKey(sessionKey string) (*string, error) {
	if protector.HashString == "" {
		return nil, errors.New("Hash code is empty")
	}

	for _, ch := range protector.HashString {
		if !unicode.IsDigit(ch) {
			return nil, errors.New(fmt.Sprintf("Hash code contains non-digit letter \"%c\"", ch))
		}
	}
	num := 0
	for idx := 0; idx < len(protector.HashString); idx++ {
		i := protector.HashString[idx]
		hashNum, _ := strconv.Atoi(protector.calculateHash(sessionKey, int(i)))
		num += hashNum
	}
	str := strings.Repeat("0", 10) + strconv.Itoa(num)[0:10]
	result := str[len(str)-10:]
	return &result, nil
}

func (protector *SessionProtector) calculateHash(sessionKey string, value int) string {
	switch value {
	case 1:
		i, _ := strconv.Atoi(sessionKey[0:5])
		str := "00" + strconv.Itoa(i%97)
		return str[len(str)-2:]
	case 2:
		result := ""
		for i := 1; i < len(sessionKey); i++ {
			result += string(sessionKey[len(sessionKey)-i])
		}
		return result + string(sessionKey[0])
	case 3:
		return string(sessionKey[len(sessionKey)-5]) + sessionKey[0:5]
	case 4:
		num := 0
		for i := 1; i < 9; i++ {
			x, _ := strconv.Atoi(string(sessionKey[i]))
			num += x + 41
		}
		return strconv.Itoa(num)
	case 5:
		num := 0
		for i := 0; i < len(sessionKey); i++ {
			ch := rune(int(sessionKey[i]) + 41)
			if !unicode.IsDigit(ch) {
				ch = rune(strconv.Itoa(int(ch))[0])
			}
			num += int(ch)
		}
		return strconv.Itoa(num)
	default:
		i, _ := strconv.Atoi(sessionKey)
		return strconv.Itoa(i + value)
	}
}
