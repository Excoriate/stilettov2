package utils

import (
	"math/rand"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytesWithSpecialChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()"

type RandomStringOptions struct {
	Length            int
	AllowSpecialChars bool
	UpperLowerRandom  bool
}

func GenerateRandomString(opt RandomStringOptions) string {
	var letterBytesRunes []rune
	if opt.AllowSpecialChars {
		letterBytesRunes = []rune(letterBytesWithSpecialChars)
	} else {
		letterBytesRunes = []rune(letterBytes)
	}

	b := make([]rune, opt.Length)
	for i := range b {
		b[i] = letterBytesRunes[rand.Intn(len(letterBytesRunes))]
	}

	if opt.UpperLowerRandom {
		for i, v := range b {
			if rand.Intn(2) == 0 {
				b[i] = []rune(strings.ToUpper(string(v)))[0]
			} else {
				b[i] = []rune(strings.ToLower(string(v)))[0]
			}
		}
	}

	return string(b)
}
