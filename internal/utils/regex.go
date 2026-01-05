package utils

import (
	"fmt"
	"os"
	"regexp"
)

var (
	SpaceRegex = func() *regexp.Regexp {
		expr, err := regexp.Compile(`\s+`)
		if err != nil {
			fmt.Println("Failed to compile regex:", err)
			os.Exit(1)
		}
		return expr
	}()
)
