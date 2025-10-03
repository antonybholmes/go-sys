package sys

import (
	"encoding/json"
	"os"
)

type VersionInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Copyright string `json:"copyright"`
	License   string `json:"license"`
	Updated   string `json:"updated"`
}

func LoadVersionInfo(path string) (VersionInfo, error) {
	var info VersionInfo

	file, err := os.Open(path)

	if err != nil {
		return info, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&info)

	return info, err
}
