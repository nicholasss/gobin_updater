package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicholasss/gobin_updater/internal/discovery"
	"github.com/nicholasss/gobin_updater/internal/fetch"
	"github.com/nicholasss/gobin_updater/internal/inventory"
	_ "github.com/nicholasss/gobin_updater/internal/updater"
)

func main() {
	// ===
	// discovery stage
	// ===

	// check for currentGOBINPath
	currentGOBINPath, err := discovery.GetGoBinPath()
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
	webInstallPath, err := discovery.WebInstallPath()
	if err != nil {
		fmt.Printf("Error discovering Webinstall install path: %q", err)
		os.Exit(1)
	}

	// check for absolute paths
	ok, err = discovery.PathsMatch(currentGOBINPath, webInstallPath)
	if !ok {
		fmt.Println("Potential error found. Check that GOBIN and $HOME/.local/opt/ are the same.")
		fmt.Println("GOBIN path:", currentGOBINPath)
		fmt.Println("WebInstall path:", webInstallPath)
	}
	if err != nil {
		fmt.Printf("Error matching GOBIN and WebInstallPath: %q", err)
		os.Exit(1)
	}

	// ===
	// inventory stage
	// ===

	// inventory of current tools
	currentVersion, err := inventory.GetCurrentInstalledGoVersion()
	if err != nil {
		fmt.Printf("Error getting runtime Golang version: %q", err)
		os.Exit(1)
	}

	fmt.Println("=======================================")
	fmt.Println("Runtime Golang Version:", currentVersion.String())
	fmt.Println("")

	installedCurrentToolList, err := inventory.ListToolsInGoBin(currentGOBINPath)
	if err != nil {
		fmt.Printf("Error getting Go bin tools: %q", err)
		os.Exit(1)
	}

	fmt.Println("=======================================")
	fmt.Println("Current GOBIN Path:", currentGOBINPath)
	fmt.Println("Currently installed tools")
	for _, tool := range installedCurrentToolList {
		fmt.Printf(" - %s\n", tool)
	}
	fmt.Println("")

	fmt.Println("=======================================")
	fmt.Println("Additional Installed:")

	_, err = fetch.FetchVersions()
	if err != nil {
		fmt.Printf("Error fetching Golang version list: %q\n", err)
		os.Exit(1)
	}

	// inventory go versions installed in webinstall path
	GOBINPaths, err := inventory.GetInstalledGoVersionPaths()
	if err != nil {
		fmt.Printf("Error fetching GOBIN paths list: %q\n", err)
		os.Exit(1)
	}

	for version, path := range GOBINPaths {
		if version.IsEqualTo(currentVersion) {
			continue
		}

		fmt.Println(" + Version", version.String())

		binPath := filepath.Join(path, "bin")
		toolsInstalled, err := inventory.ListToolsInGoBin(binPath)
		if err != nil {
			fmt.Printf("Error taking inventory in path: %q, %q\n", path, err)
		}

		for _, tool := range toolsInstalled {
			fmt.Println("   -", tool)
		}

	}
}
