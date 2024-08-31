package env

import (
	"os"
	"sort"
	"strconv"
	"time"

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

func Get(name string) string {

	return os.Getenv(name)
}

func GetStr(name string, def string) string {
	ret := Get(name)

	if ret == "" {
		return ret
	}

	return def
}

func GetUint32(name string, def uint) uint {
	v := Get(name)

	if v != "" {
		c, err := strconv.ParseUint(v, 10, 32)

		if err == nil {
			return uint(c)
		}
	}

	return def
}

// Interpret an env variable as a duration or return
// a default if the variable is not found
func GetMin(name string, def time.Duration) time.Duration {
	v := Get(name)

	if v != "" {
		c, err := strconv.ParseUint(v, 10, 32)

		if err == nil {
			return time.Duration(c) * time.Minute
		}
	}

	return def
}
