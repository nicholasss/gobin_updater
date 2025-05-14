package main

import (
	"fmt"
	"os"
	"os/exec"
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

func main() {
	gobin, err := getGoBinPath()
	if err != nil {
		fmt.Printf("Error getting go bin path: %q", err)
	}

	toolList, err := listToolsInGoBin(gobin)
	if err != nil {
		fmt.Printf("Error getting go bin tools: %q", err)
	}

	for _, tool := range toolList {
		fmt.Printf("%s\n", tool)
	}
}
