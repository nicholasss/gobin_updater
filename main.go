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
	gobin, err := discovery.GetGoBinPath()
	if err != nil {
		fmt.Printf("Error getting go bin path: %q", err)
		os.Exit(1)
	}

	toolList, err := inventory.ListToolsInGoBin(gobin)
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
