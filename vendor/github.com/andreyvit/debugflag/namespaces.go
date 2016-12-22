package debugflag

import (
	"strings"
)

func parentNamespace(flag string) string {
	pos := strings.LastIndex(flag, ":")
	if pos < 0 {
		return ""
	} else {
		return flag[0:pos]
	}
}
