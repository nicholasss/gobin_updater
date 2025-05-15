package fetch

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

const (
	GoVersionListURL = "https://go.dev/dl/?mode=json&include=all"
	CacheDirName     = "goblin"
	CacheFileName    = "goblin_versions.json"
)

type GoVersionList struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

func checkForCache()

func loadFromCache()

func saveToCache(versionList *[]GoVersionList) error {
	// get cache path
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	toolCacheDir := filepath.Join(cacheDir, CacheDirName)
	err = os.MkdirAll(toolCacheDir, 0755)
	if err != nil {
		return err
	}

	// create dir if needed
	toolCachePath := filepath.Join(toolCacheDir, CacheFileName)
	toolCacheFile, err := os.Create(toolCachePath)
	if err != nil {
		return err
	}
	defer toolCacheFile.Close()

	// save to json
	err = json.NewEncoder(toolCacheFile).Encode(versionList)
	if err != nil {
		return err
	}

	return nil
}

// fetches version list
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
