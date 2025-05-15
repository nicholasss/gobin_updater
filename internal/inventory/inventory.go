package inventory

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type GoVersion struct {
	Major int8
	Minor int8
	Patch int8
}

func (gov *GoVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", gov.Major, gov.Minor, gov.Patch)
}

func GetCurrentInstalledGoVersion() (*GoVersion, error) {
	versionStr := runtime.Version()
	versionStr = strings.TrimPrefix(versionStr, "go")

	versionArr := strings.Split(versionStr, ".")

	major, err := strconv.ParseInt(versionArr[0], 10, 8)
	if err != nil {
		return nil, err
	}

	minor, err := strconv.ParseInt(versionArr[1], 10, 8)
	if err != nil {
		return nil, err
	}

	patch, err := strconv.ParseInt(versionArr[2], 10, 8)
	if err != nil {
		return nil, err
	}

	foundGoVersion := GoVersion{
		Major: int8(major),
		Minor: int8(minor),
		Patch: int8(patch),
	}

	return &foundGoVersion, nil
}

// lists tools installed in gobin
func ListToolsInGoBin(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	toolList := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		toolList = append(toolList, entry.Name())
	}

	return toolList, nil
}

// specific to webinstall path
// $HOME/.local/opt/
func GetInstalledGoVersions() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	webinstallPath := filepath.Join(homeDir, "/.local/opt/")
	webinstallDir, err := os.ReadDir(webinstallPath)
	if err != nil {
		return nil, err
	}

	versionNames := make([]string, 0)
	for _, dir := range webinstallDir {
		if !strings.Contains(dir.Name(), "go-bin-") {
			versionNames = append(versionNames, dir.Name())
		}
	}

	return versionNames, nil
}
