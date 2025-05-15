package discovery

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ===
// */discovery package is utilized to discover how go was installed and reading env vars
// 		as well as discovering tool directories
// ===

// gets gobin from os env vars
func GetGoBinPath() (string, error) {
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

// discovers if current user has webi installed
func IsWebiUsed() (bool, error) {
	var goBinExists bool
	var webiExecExists bool

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	// looking at $HOME/.local/opt
	localOptPath := filepath.Join(userHomeDir, "/.local/opt")
	optContents, err := os.ReadDir(localOptPath)
	if err != nil {
		return false, err
	}

	for _, optContent := range optContents {
		if strings.Contains(optContent.Name(), "go") {
			goBinExists = true
			break
		}
	}

	// looking at $HOME/.local/bin
	localBinPath := filepath.Join(userHomeDir, "/.local/bin")
	webiPath := filepath.Join(localBinPath, "/webi")

	webiBinInfo, err := os.Stat(webiPath)
	if err != nil {
		return false, err
	}

	if webiBinInfo.Size() != 0 {
		webiExecExists = true
	}

	return goBinExists && webiExecExists, nil
}
