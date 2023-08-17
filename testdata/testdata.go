package testdata

import (
	"os"
	"testing"
)

func UnreadableFile() []byte {
	return []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x34, 0x25, 0x25, 0x45, 0x4F, 0x46, 0x0A}
}

func LoadFile(t *testing.T) []byte {
	file, err := os.ReadFile("../../testdata/file1.pdf")
	if err != nil {
		t.Fatal("testdata/file1.pdf not found")
	}

	return file
}

func LoadFilesWithInvalid(includeInvalid bool, t *testing.T) [][]byte {
	files := [][]byte{LoadFile(t), LoadFile(t)}

	if includeInvalid {
		files = append(files, UnreadableFile())
	}

	return files
}
