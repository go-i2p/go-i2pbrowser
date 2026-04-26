package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	goi2pbrowser "github.com/go-i2p/go-i2pbrowser"
)

var (
	u = flag.Bool("usability", false, "Launch in usability mode")
	a = flag.Bool("application", false, "Launch in app mode")
	d = flag.String("directory", "", "Directory to store profiles in")
)

func defaultDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	ret := filepath.Join(wd, "i2p-profiles")
	userHome, err := os.UserHomeDir()
	if err != nil {
		return ret
	}
	if wd == userHome {
		ret = filepath.Join(userHome, ".i2p/plugins/i2pbrowser")
	}
	return ret
}

func browse(profileDir, url string) error {
	switch {
	case *a:
		return goi2pbrowser.BrowseApp(profileDir, url)
	case *u:
		return goi2pbrowser.BrowseUsability(profileDir, url)
	default:
		return goi2pbrowser.BrowseStrict(profileDir, url)
	}
}

func main() {
	flag.Parse()
	profileDir := *d
	if profileDir == "" {
		profileDir = defaultDir()
	}
	if profileDir == "" {
		log.Fatal("could not determine profile directory: set -directory explicitly")
	}
	url := "http://127.0.0.1:7657"
	if len(flag.Args()) > 0 {
		url = flag.Arg(0)
	}
	if err := browse(profileDir, url); err != nil {
		log.Fatal(err)
	}
}
