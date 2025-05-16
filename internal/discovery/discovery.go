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

	// resolve symlinks
	resolvedGobin, err := filepath.EvalSymlinks(gobin)
	if err != nil {
		return "", err
	}

	return resolvedGobin, nil
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

// returns local opt path, where webi installs go
// $HOME/.local/opt
func WebInstallPath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userHomeDir, "/.local/opt"), nil
}

func PathsMatch(gobinPath, webiPath string) (bool, error) {
	gobinPathInfo, err := os.Stat(gobinPath)
	if err != nil {
		return false, err
	}

	webiPathInfo, err := os.Stat(webiPath)
	if err != nil {
		return false, err
	}

	if gobinPath == webiPath && gobinPathInfo.Name() == webiPathInfo.Name() {
		return true, nil
	}

	// check two up from gobin path
	gobinPathParent := filepath.Dir(gobinPath)
	gobinPathGrandparent := filepath.Dir(gobinPathParent)
	gobinPathGrandparentInfo, err := os.Stat(gobinPathGrandparent)
	if err != nil {
		return false, err
	}

	if gobinPathGrandparent == webiPath && gobinPathGrandparentInfo.Name() == webiPathInfo.Name() {
		return true, nil
	}

	return false, nil
}
