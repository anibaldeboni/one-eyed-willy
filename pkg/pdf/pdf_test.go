package pdf

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	type args struct {
		files [][]byte
	}
	tests := []struct {
		name          string
		args          args
		resultIsEmpty bool
		err           error
	}{
		{
			name:          "When all files could be read",
			args:          args{files: loadFiles(false, t)},
			resultIsEmpty: false,
			err:           nil,
		},
		{
			name:          "When some file is invalid",
			args:          args{files: loadFiles(true, t)},
			resultIsEmpty: true,
			err:           errors.New("Could not merge pdfs. Some files can't be read"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Merge(tt.args.files)
			assert.IsType(t, tt.err, err)
			assert.Equal(t, tt.resultIsEmpty, result == nil)
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
