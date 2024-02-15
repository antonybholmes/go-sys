package env

import (
	"os"
	"strconv"
)

func GetStr(name string, def string) string {
	ret := os.Getenv(name)

	if ret == "" {
		ret = def
	}

	return ret
}

func GetUint32(name string, def uint) uint {
	v := os.Getenv(name)

	var ret uint

	if v != "" {
		c, err := strconv.ParseUint(v, 10, 32)

		if err == nil {
			ret = uint(c)
		} else {
			ret = def
		}
	}

	return ret
}
