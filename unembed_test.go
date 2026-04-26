package goi2pbrowser

import (
	"os"
	"testing"
)

func TestUnEmbed(t *testing.T) {
	t.Log("testing base")
	basePath, err := UnpackBase("testing/Base")
	if err != nil {
		t.Fatalf("UnpackBase failed: %v", err)
	}
	if _, err := os.Stat(basePath); err != nil {
		t.Fatalf("base profile path %q does not exist after unpack: %v", basePath, err)
	}

	t.Log("testing usability")
	usabilityPath, err := UnpackUsability("testing/Usability")
	if err != nil {
		t.Fatalf("UnpackUsability failed: %v", err)
	}
	if _, err := os.Stat(usabilityPath); err != nil {
		t.Fatalf("usability profile path %q does not exist after unpack: %v", usabilityPath, err)
	}
}
