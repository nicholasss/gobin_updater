package main

import (
	"fmt"
	"os"

	"github.com/nicholasss/gobin_updater/internal/discovery"
	"github.com/nicholasss/gobin_updater/internal/fetch"
	"github.com/nicholasss/gobin_updater/internal/inventory"
	_ "github.com/nicholasss/gobin_updater/internal/updater"
)

func main() {
	// ===
	// discovery stage
	// ===

	// check for GOBINPath
	GOBINPath, err := discovery.GetGoBinPath()
	if err != nil {
		fmt.Printf("Error discovering GOBIN: %q", err)
		os.Exit(1)
	}

	// check for webi
	ok, err := discovery.IsWebiUsed()
	if err != nil {
		fmt.Printf("Error discovering if Webinstall is used: %q", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Println("Webinstall is not detected. Exiting program.")
		os.Exit(1)
	}

	// get webinstall bin path
	webiInstallPath, err := discovery.WebiInstallPath()
	if err != nil {
		fmt.Printf("Error discovering Webinstall install path", err)
		os.Exit(1)
	}

	// check for absolute paths
	if discovery.PathsMatch(GOBINPath, webiInstallPath)

	// ===
	// inventory stage
	// ===

	toolList, err := inventory.ListToolsInGoBin(GOGOBINPath)
	if err != nil {
		fmt.Printf("Error getting go bin tools: %q", err)
		os.Exit(1)
	}

	for _, tool := range toolList {
		fmt.Printf("%s\n", tool)
	}

	gov, err := inventory.GetCurrentInstalledGoVersion()
	if err != nil {
		fmt.Printf("Error getting current go version: %q", err)
		os.Exit(1)
	}

	fmt.Println("go version:", gov.String())

	versionList, err := fetch.FetchGoVersionList()
	if err != nil {
		fmt.Printf("Error fetching go versions: %q", err)
		os.Exit(1)
	}

	fmt.Println("Go Versions")
	for _, version := range *versionList {
		if !version.Stable {
			continue
		}

		fmt.Printf("Version: %s\n", version.Version)
	}
}
