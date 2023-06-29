package pdf

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergePdfs(t *testing.T) {
	type args struct {
		files [][]byte
	}
	tests := []struct {
		name     string
		args     args
		fileSize int
		err      error
	}{
		{
			name:     "When all files could be read",
			args:     args{files: loadFiles(false, t)},
			fileSize: 319359,
			err:      nil,
		},
		{
			name:     "When some files are invalid",
			args:     args{files: loadFiles(true, t)},
			fileSize: 0,
			err:      errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MergePdfs(tt.args.files)
			assert.IsType(t, tt.err, err)
			assert.Equal(t, tt.fileSize, len(result))
		})
	}
}

func loadFiles(includeInvalid bool, t *testing.T) [][]byte {
	file1, err := os.ReadFile("../../testdata/file1.pdf")
	if err != nil {
		t.Fatal("testdata/file1.pdf not found")
	}

	file2, err := os.ReadFile("../../testdata/file2.pdf")
	if err != nil {
		t.Fatal("testdata/file2.pdf not found")
	}

	files := [][]byte{file1, file2}

	if includeInvalid == true {
		fileInvalid, err := os.ReadFile("../../testdata/file-invalid.pdf")
		if err != nil {
			t.Fatal("testdata/file-invalid.pdf not found")
		}
		files = append(files, fileInvalid)
	}

	return files
}
