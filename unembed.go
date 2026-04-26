package goi2pbrowser

import (
	"log"
	"os"
	"path/filepath"

	"github.com/artdarek/go-unzip/pkg/unzip"
)

func existsAlready(profileDir string) bool {
	if _, err := os.Stat(filepath.Join(profileDir, "user.js")); err == nil {
		return true
	}
	return false
}

func baseProfilePath(profileDir string) string {
	return filepath.Join(profileDir, "i2p.firefox.base.profile")
}

func usabilityProfilePath(profileDir string) string {
	return filepath.Join(profileDir, "i2p.firefox.usability.profile")
}

// UnpackBase unpacks a "Strict" mode profile into the "profileDir" and returns the
// path to the profile and possibly, an error if something goes wrong. If everything
// works, the error will be nil.
//
// Note: a ZIP archive is written to the parent directory of profileDir as a side
// effect and is not removed after extraction.
func UnpackBase(profileDir string) (string, error) {
	profilePath := baseProfilePath(profileDir)
	if existsAlready(profilePath) {
		log.Println(profilePath, "exists already")
		return profilePath, nil
	}
	if err := os.MkdirAll(filepath.Dir(profileDir), 0o755); err != nil {
		return profilePath, err
	}
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		return profilePath, err
	}
	zipFile := filepath.Join(filepath.Dir(profileDir), "i2p.firefox.base.profile.zip")
	err := os.WriteFile(zipFile, BaseProfile, 0o644)
	if err != nil {
		return profilePath, err
	}
	uz := unzip.New()
	_, err = uz.Extract(zipFile, profileDir)
	if err != nil {
		return profilePath, err
	}
	return profilePath, nil
}

// UnpackUsability unpacks a "Usability" mode profile into the "profileDir" and returns the
// path to the profile and possibly, an error if something goes wrong. If everything
// works, the error will be nil.
//
// Note: a ZIP archive is written to the parent directory of profileDir as a side
// effect and is not removed after extraction.
func UnpackUsability(profileDir string) (string, error) {
	profilePath := usabilityProfilePath(profileDir)
	if existsAlready(profilePath) {
		log.Println(profilePath, "exists already")
		return profilePath, nil
	}
	if err := os.MkdirAll(filepath.Dir(profileDir), 0o755); err != nil {
		return profilePath, err
	}
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		return profilePath, err
	}
	zipFile := filepath.Join(filepath.Dir(profileDir), "i2p.firefox.usability.profile.zip")
	err := os.WriteFile(zipFile, UsabilityProfile, 0o644)
	if err != nil {
		return profilePath, err
	}
	uz := unzip.New()
	_, err = uz.Extract(zipFile, profileDir)
	if err != nil {
		return profilePath, err
	}
	return profilePath, nil
}
