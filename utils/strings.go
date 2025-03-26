package utils

import (
	"fmt"
	"strings"
)

func PrintWithPrefix(prefix string, message string) {
	for _, line := range strings.Split(strings.Trim(message, "\n"), "\n") {
		fmt.Println(prefix, line)
	}
}
