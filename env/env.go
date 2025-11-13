package env

import (
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/antonybholmes/go-sys/log"
	"github.com/joho/godotenv"
)

func init() {
	log.Debug().Msgf("loading envs...")

	// force reloading
	err := godotenv.Overload()

	if err != nil {
		log.Debug().Msgf("error loading .env file(s)")
	}

}

func Load(file string) {
	godotenv.Load(file)
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

	if ret != "" {
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
	return GetTime(name, time.Minute, def)
}

func GetHour(name string, def time.Duration) time.Duration {
	return GetTime(name, time.Hour, def)
}

func GetTime(name string, unit time.Duration, def time.Duration) time.Duration {
	v := Get(name)

	if v != "" {
		c, err := strconv.ParseInt(v, 10, 32)

		if err == nil {
			log.Debug().Msgf("found %v with value %d", name, c)
			return time.Duration(c) * unit
		}
	}

	return def
}
