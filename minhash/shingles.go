package minhash

import (
	"bytes"
	"github.com/kljensen/snowball"
	"strings"
	"unicode"
)

const length = 3
const lang  = "russian"


func SplitShingles(line string) (shingles []string) {
	line, err := snowball.Stem(line, string(lang), true)
	if err != nil {
		return
	}

	var filtered bytes.Buffer

	for _, r := range []rune(line) {
		if unicode.IsLetter(r) || r == ' ' {
			filtered.WriteRune(r)
		}
	}

	var words = strings.Split(filtered.String(), " ")

	max := len(words) - length
	for i := 0; i < max; i++ {
		shingles = append(shingles, strings.Join(words[i:i + length], ""))
	}

	return
}
