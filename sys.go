package sys

import (
	"encoding/base64"
	"unsafe"

	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// panics on error or returns value, similar to google's
// must in uuid but this is generic
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Uuid() string {
	return uuid.New().String() // strings.ReplaceAll(u1.String(), "-", ""), nil
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

func NanoId() string {
	// good enough for Planetscale https://planetscale.com/blog/why-we-chose-nanoids-for-planetscales-api
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
