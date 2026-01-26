package util

import (
	"crypto/rand"
	"encoding/binary"
	mrand "math/rand"
	"strings"
	"time"
)

var syllables = []string{
	"ba", "be", "bi", "bo", "bu",
	"ca", "ce", "ci", "co", "cu",
	"da", "de", "di", "do", "du",
	"fa", "fe", "fi", "fo", "fu",
	"ga", "ge", "gi", "go", "gu",
	"ha", "he", "hi", "ho", "hu",
	"ka", "ke", "ki", "ko", "ku",
	"la", "le", "li", "lo", "lu",
	"ma", "me", "mi", "mo", "mu",
	"na", "ne", "ni", "no", "nu",
	"pa", "pe", "pi", "po", "pu",
	"ra", "re", "ri", "ro", "ru",
	"sa", "se", "si", "so", "su",
	"ta", "te", "ti", "to", "tu",
	"va", "ve", "vi", "vo", "vu",
	"za", "ze", "zi", "zo", "zu",
	"an", "en", "in", "on", "un",
	"ar", "er", "ir", "or", "ur",
	"al", "el", "il", "ol", "ul",
}

var lastParts = []string{
	"wood", "stone", "river", "brook", "field", "hill", "ford", "well",
	"hart", "wolf", "fox", "bear", "oak", "pine", "vale", "croft",
	"smith", "mason", "ward", "turner", "carter", "baker",
}

func cryptoSeed() int64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err == nil {
		return int64(binary.LittleEndian.Uint64(b[:]))
	}
	// Fallback if crypto/rand fails (rare): time-based.
	return time.Now().UnixNano()
}

func capFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func genFirstName(r *mrand.Rand, syllableCount int) string {
	if syllableCount < 2 {
		syllableCount = 2
	}
	if syllableCount > 5 {
		syllableCount = 5
	}

	var sb strings.Builder
	for i := 0; i < syllableCount; i++ {
		sb.WriteString(syllables[r.Intn(len(syllables))])
	}
	return capFirst(sb.String())
}

func GetFirstName() string {
	r := mrand.New(mrand.NewSource(cryptoSeed()))

	return genFirstName(r, 5)
}
