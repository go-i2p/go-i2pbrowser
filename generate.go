//go:build generate
// +build generate

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
)

var (
	owner = "eyedeekay"
	repo  = "i2p.plugins.firefox"
)

func validVersion(name string) bool {
	vers := strings.Split(name, ".")
	if len(vers) == 3 {
		return true
	}
	return false
}

func profileVersion() string {
	client := github.NewClient(nil)
	tags, _, err := client.Repositories.ListTags(context.Background(), owner, repo, nil)
	if err != nil {
		fmt.Println(err)
		return "2.8.3"
	}
	if len(tags) > 0 {
		for _, tag := range tags {
			if validVersion(*tag.Name) {
				latestTag := tag
				fmt.Printf("Latest tag '%s', (SHA-1: %s)\n", *latestTag.Name, *latestTag.Commit.SHA)
				return *latestTag.Name
			}
		}
	} else {
		fmt.Printf("No tags yet\n")
	}
	return "2.8.3"
}

func DownloadURL(version, mode string) string {
	return "https://github.com/" + owner + "/" + repo + "/releases/download/" + version + "/i2p.firefox." + mode + ".profile.zip"
}

// Downloads a file from a URL to path output. Outputs absolute path to download or an error
func DownloadFile(URL, output string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(wd, output)
	log.Println("Downloading:", URL, "to", filePath)
	resp, err := http.Get(URL)
	if err != nil {
		return filePath, err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) // drain before close for connection reuse
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return filePath, err
	}
	err = os.WriteFile(filePath, bodyBytes, 0o644)
	if err != nil {
		return filePath, err
	}
	return filePath, nil
}

func main() {
	log.Println("Downloading profiles from upstream bundles")
	version := profileVersion()
	log.Println("Version is:", version)
	basePath, err := DownloadFile(DownloadURL(version, "base"), "i2p.firefox.base.profile.zip")
	if err != nil {
		panic(err)
	}
	log.Println("Downloaded:", basePath)
	usabilityPath, err := DownloadFile(DownloadURL(version, "usability"), "i2p.firefox.usability.profile.zip")
	if err != nil {
		panic(err)
	}
	log.Println("Downloaded:", usabilityPath)
}
