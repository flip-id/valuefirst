package valuefirst

import (
	"fmt"
	"github.com/fairyhunter13/pool"
	"html"
)

// These constants represents the lower and upper limit of the characters that don't have to be encoded.
const (
	LowerLimitEncoding = 32
	UpperLimitEncoding = 128
)

func isEncoded(chr rune) bool {
	return chr > UpperLimitEncoding ||
		chr < LowerLimitEncoding ||
		chr == '*' ||
		chr == '#' ||
		chr == '%' ||
		chr == '<' ||
		chr == '>' ||
		chr == '+'
}

// Encode encodes the message using the step 2 encoding of ValueFirst documentation.
func Encode(msg string) (res string) {
	builder := pool.GetStrBuilder()
	defer pool.Put(builder)

	for _, chr := range msg {
		writeStr := string(chr)
		if isEncoded(chr) {
			writeStr = fmt.Sprintf("%%%X", chr)
		}

		_, _ = builder.WriteString(writeStr)
	}

	res = builder.String()
	return
}

var (
	htmlSpecialChars = map[rune]string{
		rune(39): "&apos",
		rune(32): "&#032",
		rune(34): "&quot",
		'>':      "&gt",
		'<':      "&lt",
		rune(13): "&#013",
		rune(10): "&#010",
		rune(9):  "&#009",
	}
)

// EncodeHTML encodes the message using the step 1 encoding of ValueFirst documentation.
// This encoding is used for HTML special characters.
func EncodeHTML(msg string) (res string) {
	builder := pool.GetStrBuilder()
	defer pool.Put(builder)

	msg = html.UnescapeString(msg)
	var msgStr string
	for _, chr := range msg {
		msgStr = string(chr)
		if str, ok := htmlSpecialChars[chr]; ok {
			msgStr = str
		}

		_, _ = builder.WriteString(msgStr)
	}

	res = builder.String()
	return
}
