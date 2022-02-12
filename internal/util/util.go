package util

import (
	"strings"
)

func GetURLFromZoneName(s string) string {
	return strings.TrimSuffix(s, ".")
}
