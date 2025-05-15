package discovery

import (
	"os"
	"os/exec"
	"strings"
)

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

// validates that this is where gobin is located
func ValidateGoBinPath() {}
