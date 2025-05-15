package fetch

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	GoVersionListURL     = "https://go.dev/dl/?mode=json&include=all"
	CacheDirName         = "goblin"
	CacheFileName        = "goblin_versions.json"
	CacheExpiresDuration = time.Hour * 120 // five days
)

type GoVersionList struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

// gets go versions
func FetchVersions() (*[]GoVersionList, error) {
	// first checks for cache
	exists, err := cacheFileExists()
	if err != nil {
		return nil, err
	}

	// check age
	if exists {
		cacheExpiryTime := time.Now().Add(-CacheExpiresDuration)
		cacheAge, err := cacheLastModified()
		if err != nil {
			return nil, err
		}

		if !cacheAge.Before(cacheExpiryTime) {
			// cache is not expired, load and return
			versionList, err := loadFromCache()
			if err != nil {
				return nil, err
			}

			return versionList, nil
		}

	}

	// cache either does not exist or is expired, fetch new
	versionList, err := fetchGoVersionList()
	if err != nil {
		return nil, err
	}

	// save to disk
	err = saveToCache(versionList)
	if err != nil {
		return nil, err
	}

	return versionList, nil
}

// utility function to get the filepath of the cache
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

func cacheFileExists() (bool, error) {
	toolCachePath, err := getCacheFilePath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(toolCachePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// checks the cache for when it was last modified
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

// loads the file from cache
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

// saves the go versions to cache
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

// fetches version list from online
func fetchGoVersionList() (*[]GoVersionList, error) {
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
