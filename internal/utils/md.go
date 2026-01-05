package utils

import (
	"fmt"
	"strings"
)

func EscapeMd2(s string) string {
	rep := []string{
		`_`, `*`, `[`, `]`, `(`, `)`, `~`, "`", `>`, `#`, `+`, `-`, `=`, `|`, `{`, `}`, `.`, `!`,
	}
	for _, r := range rep {
		s = strings.ReplaceAll(s, r, `\`+r)
	}
	fmt.Println(s)

	//s = strings.ReplaceAll(s, ".", "\\.")
	//s = strings.ReplaceAll(s, "-", "\\-")

	return s
}
