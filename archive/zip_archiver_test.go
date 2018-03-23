package archive

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"testing"
)

func TestZipArchiver_Content(t *testing.T) {
	zipfilepath := "archive-content.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveContent([]byte("This is some content"), "content.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"content.txt": []byte("This is some content"),
	})

	ensureChecksum(t, zipfilepath, "P7VckxoEiUO411WN3nwuS/yOBL4zsbVWkQU9E1I5H6c=")
}

func TestZipArchiver_File(t *testing.T) {
	zipfilepath := "archive-file.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveFile("./test-fixtures/test-file.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"test-file.txt": []byte("This is test content"),
	})

	ensureChecksum(t, zipfilepath, "7Ozdhchkz12Ae7ZMtXQ4jqKlV5NWUjY2qgAxRflv0UA=")
}

func TestZipArchiver_Dir(t *testing.T) {
	zipfilepath := "archive-dir.zip"
	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	})

	ensureChecksum(t, zipfilepath, "9tByVm9Ik6q8Zodh5FoMssxapEdlWptHuxxcUMk+j4w=")
}

func TestZipArchiver_Multiple(t *testing.T) {
	zipfilepath := "archive-content.zip"
	content := map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	}

	archiver := NewZipArchiver(zipfilepath)
	if err := archiver.ArchiveMultiple(content); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipfilepath, content)

	ensureChecksum(t, zipfilepath, "LrddibLkkT50VJied6+dNmh8hLzADtihWtxSNL/UdYY=")
}

func ensureContents(t *testing.T, zipfilepath string, wants map[string][]byte) {
	r, err := zip.OpenReader(zipfilepath)
	if err != nil {
		t.Fatalf("could not open zip file: %s", err)
	}
	defer r.Close()

	if len(r.File) != len(wants) {
		t.Errorf("mismatched file count, got %d, want %d", len(r.File), len(wants))
	}
	for _, cf := range r.File {
		ensureContent(t, wants, cf)
	}
}

func ensureContent(t *testing.T, wants map[string][]byte, got *zip.File) {
	want, ok := wants[got.Name]
	if !ok {
		t.Errorf("additional file in zip: %s", got.Name)
		return
	}

	r, err := got.Open()
	if err != nil {
		t.Errorf("could not open file: %s", err)
	}
	defer r.Close()
	gotContentBytes, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("could not read file: %s", err)
	}

	wantContent := string(want)
	gotContent := string(gotContentBytes)
	if gotContent != wantContent {
		t.Errorf("mismatched content\ngot\n%s\nwant\n%s", gotContent, wantContent)
	}
}

func ensureChecksum(t *testing.T, zipfilepath string, wantChecksum string) {
	data, err := ioutil.ReadFile(zipfilepath)
	if err != nil {
		t.Errorf("could not compute file '%s' checksum: %s", zipfilepath, err)
	}
	h256 := sha256.New()
	h256.Write([]byte(data))
	shaSum := h256.Sum(nil)
	gotChecksum := base64.StdEncoding.EncodeToString(shaSum[:])
	if gotChecksum != wantChecksum {
		t.Errorf("mismatched checksum\ngot\n%s\nwant\n%s", gotChecksum, wantChecksum)
	}
}
