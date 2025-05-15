package fetch

import (
	"encoding/json"
	"net/http"
)

const (
	GoVersionListURL = "https://go.dev/dl/?mode=json&include=all"
)

type GoVersionList struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

func FetchGoVersionList() (*[]GoVersionList, error) {
	resp, err := http.Get(GoVersionListURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var versionList []GoVersionList
	err = json.NewDecoder(resp.Body).Decode(&versionList)
	if err != nil {
		return nil, err
	}

	return &versionList, nil
}
