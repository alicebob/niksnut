package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("hello.\n")
	configFile := "./config.nix"
	config, err := ReadConfig(configFile)
	if err != nil {
		fmt.Printf("error reading %s: %s\n", configFile, err)
		os.Exit(1)
	}

	fmt.Printf("found projects:\n")
	for _, p := range config.Projects {
		fmt.Printf("  - %s\n", p.ID)
	}

	p := config.Projects[0]
	out, err := Run(p)
	if err != nil {
		fmt.Printf("error running %s: %s\n", p.ID, err)
		os.Exit(1)
	}
	fmt.Printf("run of %s:\n", p.ID)
	fmt.Print(out.Stdout)
}
