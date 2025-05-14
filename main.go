package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func getGoBinPath() (string, error) {
	gobin := os.Getenv("GOBIN")

	if gobin == "" {
		// If GOBIN isn't set, fallback to GOPATH/bin
		gopath := os.Getenv("GOPATH")

		if gopath == "" {
			// Use `go env GOPATH` if GOPATH isn't explicitly set
			out, err := exec.Command("go", "env", "GOPATH").Output()
			if err != nil {
				return "", err
			}

			gopath = strings.TrimSpace(string(out))
		}

		gobin = gopath + "/bin"
	}

	return gobin, nil
}

func listToolsInGoBin(dir string) ([]string, error) {
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

type GoVersion struct {
	Major int8
	Minor int8
	Patch int8
}

func (gov *GoVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", gov.Major, gov.Minor, gov.Patch)
}

func GetCurrentGoVersion() (*GoVersion, error) {
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

func main() {
	gobin, err := getGoBinPath()
	if err != nil {
		fmt.Printf("Error getting go bin path: %q", err)
		os.Exit(1)
	}

	toolList, err := listToolsInGoBin(gobin)
	if err != nil {
		fmt.Printf("Error getting go bin tools: %q", err)
		os.Exit(1)
	}

	for _, tool := range toolList {
		fmt.Printf("%s\n", tool)
	}

	gov, err := GetCurrentGoVersion()
	if err != nil {
		fmt.Printf("Error getting current go version: %q", err)
		os.Exit(1)
	}

	fmt.Println("go version:", gov.String())
}
