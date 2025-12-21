package sys

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	Sqlite3DB  = "sqlite3"
	PostgresDB = "postgres"
	MySQLDB    = "mysql"
	BlankUUID  = "00000000-0000-0000-0000-000000000000"
)

// panics on error or returns value, similar to google's
// must in uuid but this is generic
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// panics on error for functions that return only error
func VoidMust(err error) {
	if err != nil {
		panic(err)
	}

}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Uuidv7() (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func B64Uuid() string {
	uuid := uuid.New()
	b := (*[]byte)(unsafe.Pointer(&uuid))
	return base64.RawURLEncoding.EncodeToString(*b)
}

func DecodeB64Uuid(id string) (*uuid.UUID, error) {
	dec, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return nil, err
	}
	decID, err := uuid.FromBytes(dec)
	if err != nil {
		return nil, err
	}
	return &decID, nil
}

// Generates a 12 char random id, which can used as guid for
// most purposes. It's good enough for Planetscale
// https://planetscale.com/blog/why-we-chose-nanoids-for-planetscales-api
func NanoId() string {
	id, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 12)

	if err != nil {
		id = ""
	}

	return id
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// Replace runs of spaces, tabs, newlines with a single space
func NormalizeSpaces(s string) string {
	var b strings.Builder
	inSpace := false

	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
		} else {
			b.WriteRune(r)
			inSpace = false
		}
	}

	return b.String()
}

func MinDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

// Atoi converts a string to int, ignoring commas
func Atoi(s string) (int, error) {
	return strconv.Atoi(strings.ReplaceAll(s, ",", ""))
}

func IsAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func IsUpperAlpha(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
