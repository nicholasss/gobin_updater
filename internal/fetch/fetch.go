package fetch

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

func getCacheFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	toolCacheDir := filepath.Join(cacheDir, CacheDirName)
	err = os.MkdirAll(toolCacheDir, 0755)
	if err != nil {
		return "", err
	}

	toolCachePath := filepath.Join(toolCacheDir, CacheFileName)
	return toolCachePath, nil
}

func cacheLastModified() (time.Time, error) {
	toolCachePath, err := getCacheFilePath()
	if err != nil {
		return time.Time{}, err
	}

	cacheInfo, err := os.Stat(toolCachePath)
	if err != nil {
		return time.Time{}, err
	}

	return cacheInfo.ModTime(), nil
}

func loadFromCache() (*[]GoVersionList, error) {
	toolCachePath, err := getCacheFilePath()
	if err != nil {
		return nil, err
	}

	cacheFile, err := os.Open(toolCachePath)
	if err != nil {
		return nil, err
	}

	var goVersionList *[]GoVersionList
	err = json.NewDecoder(cacheFile).Decode(goVersionList)
	if err != nil {
		return nil, err
	}

	return goVersionList, nil
}

func saveToCache(versionList *[]GoVersionList) error {
	toolCachePath, err := getCacheFilePath()
	if err != nil {
		return err
	}

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
