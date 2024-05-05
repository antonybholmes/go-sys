package env

import (
	"os"
	"sort"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func Load() {
	// force reloading
	err := godotenv.Overload()

	if err != nil {
		log.Fatal().Msgf("Error loading .env file(s)")
	}

}

func Ls() {

	envs := os.Environ()

	sort.Strings(envs)

	for _, e := range envs {
		log.Debug().Msgf("%s", e)
	}
}

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
