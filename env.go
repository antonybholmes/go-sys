package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func Load() error {
	err := godotenv.Load()

	if err != nil {
		log.Error().Msgf("Error loading .env file(s)")
	}

	return err
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
