package main

import (
	"fmt"
	"os"

	"github.com/nicholasss/gobin_updater/internal/discovery"
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
	webInstallPath, err := discovery.WebInstallPath()
	if err != nil {
		fmt.Printf("Error discovering Webinstall install path: %q", err)
		os.Exit(1)
	}

	// check for absolute paths
	ok, err = discovery.PathsMatch(GOBINPath, webInstallPath)
	if !ok {
		fmt.Println("Potential error found. Check that GOBIN and $HOME/.local/opt/ are the same.")
		fmt.Println("GOBIN path:", GOBINPath)
		fmt.Println("WebInstall path:", webInstallPath)
	}
	if err != nil {
		fmt.Printf("Error matching GOBIN and WebInstallPath: %q", err)
		os.Exit(1)
	}

	// ===
	// inventory stage
	// ===

	currentGoVersion, err := inventory.GetCurrentInstalledGoVersion()
	if err != nil {
		fmt.Printf("Error getting runtime Golang version: %q", err)
		os.Exit(1)
	}

	fmt.Println("=======================================")
	fmt.Println("Runtime Golang Version:", currentGoVersion.String())
	fmt.Println("")

	installedGoToolList, err := inventory.ListToolsInGoBin(GOBINPath)
	if err != nil {
		fmt.Printf("Error getting Go bin tools: %q", err)
		os.Exit(1)
	}

	fmt.Println("=======================================")
	fmt.Println("Current GOBIN Path:", GOBINPath)
	fmt.Println("Currently installed tools")
	for _, tool := range installedGoToolList {
		fmt.Printf(" - %s\n", tool)
	}
	fmt.Println("")

	fmt.Println("=======================================")
	fmt.Println("Additional Installed:")
	// TODO:
	//  : version 1.24.1
	//    - sqlc
	//    - goose
	//  : version 1.23.7
	//    - sqlc
	//    - goose
	// etc.

	// versionList, err := fetch.FetchGoVersionList()
	// if err != nil {
	// 	fmt.Printf("Error fetching go versions: %q", err)
	// 	os.Exit(1)
	// }

	// fmt.Println("Go Versions")
	// for _, version := range *versionList {
	// 	if !version.Stable {
	// 		continue
	// 	}

	// 	fmt.Printf("Version: %s\n", version.Version)
	// }
}
