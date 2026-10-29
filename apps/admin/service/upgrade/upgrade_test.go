package upgrade

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"v1.0.1", "v1.0.0", 1},
		{"v1.0.0", "v1.0.1", -1},
		{"v1.0.0", "v1.0.0", 0},
		{"1.0.10", "1.0.2", 1},
		{"v1.1.0", "v1.0.9", 1},
		{"1.0", "1.0.0", 0},
		{"v2.0.0", "v1.9.9", 1},
	}

	for _, tt := range tests {
		res := CompareVersions(tt.v1, tt.v2)
		if res != tt.expected {
			t.Errorf("CompareVersions(%s, %s) = %d, expected %d", tt.v1, tt.v2, res, tt.expected)
		}
	}
}

func TestMatchTargetAsset(t *testing.T) {
	assets := []ReleaseAsset{
		{Name: "gbootcms-linux.tar.gz", BrowserDownloadURL: "https://example.com/linux.tar.gz"},
		{Name: "gbootcms-windows.zip", BrowserDownloadURL: "https://example.com/windows.zip"},
	}

	linuxAsset := MatchTargetAsset(assets, "linux", "amd64")
	if linuxAsset == nil || linuxAsset.Name != "gbootcms-linux.tar.gz" {
		t.Errorf("Expected linux asset, got %v", linuxAsset)
	}

	winAsset := MatchTargetAsset(assets, "windows", "amd64")
	if winAsset == nil || winAsset.Name != "gbootcms-windows.zip" {
		t.Errorf("Expected windows asset, got %v", winAsset)
	}

	unknownAsset := MatchTargetAsset(assets, "darwin", "arm64")
	if unknownAsset != nil {
		t.Errorf("Expected nil for darwin, got %v", unknownAsset)
	}
}

func TestExtractBinaryFromTarGz(t *testing.T) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	content := []byte("#!/bin/sh\necho ok")
	hdr := &tar.Header{
		Name: "bin/gbootcms",
		Mode: 0755,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gw.Close()

	extracted, err := ExtractBinaryFromTarGz(buf.Bytes())
	if err != nil {
		t.Fatalf("ExtractBinaryFromTarGz failed: %v", err)
	}
	if !bytes.Equal(extracted, content) {
		t.Errorf("Extracted bytes mismatch: got %s, expected %s", string(extracted), string(content))
	}
}

func TestExtractBinaryFromZip(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	content := []byte("MZdummyexe")
	f, err := zw.Create("bin/gbootcms.exe")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(content); err != nil {
		t.Fatal(err)
	}
	zw.Close()

	extracted, err := ExtractBinaryFromZip(buf.Bytes())
	if err != nil {
		t.Fatalf("ExtractBinaryFromZip failed: %v", err)
	}
	if !bytes.Equal(extracted, content) {
		t.Errorf("Extracted bytes mismatch: got %s, expected %s", string(extracted), string(content))
	}
}
