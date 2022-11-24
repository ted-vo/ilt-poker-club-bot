package pkg

import "math/rand"

const (
	MARKDOWN string = "MarkdownV2"
	HTLM     string = "HTML"
)

func Rollem() int {
	return rand.Intn(12) + 1
}
