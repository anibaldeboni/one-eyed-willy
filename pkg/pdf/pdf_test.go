package pdf

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFromHTML(t *testing.T) {
	pdfRender := NewRender()
	type args struct {
		html string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "When html is valid",
			args: args{html: "<h1>Hello World</h1>"},
		},
		{
			name: "When html is empty",
			args: args{html: ""},
		},
		{
			name: "When html is invalid",
			args: args{html: "invalid-html"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateFromHTML(pdfRender.Context, tt.args.html)
			fBytes, _ := io.ReadAll(result)
			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.NoError(t, ValidatePdf(fBytes))
		})
	}
}

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
		{
			name:          "When the files can't be merged",
			args:          args{files: [][]byte{unreadableFile(), unreadableFile()}},
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

func TestEncrypt(t *testing.T) {
	type args struct {
		files    []byte
		password string
	}
	tests := []struct {
		name   string
		args   args
		err    error
		result bool
	}{
		{
			name:   "When the file is valid",
			args:   args{files: loadFiles(false, t)[0], password: "test"},
			err:    nil,
			result: true,
		},
		{
			name:   "When the file is invalid",
			args:   args{files: []byte(`not-a-pdf`), password: "test"},
			err:    errors.New("This is not a pdf file"),
			result: false,
		},
		{
			name:   "When the file is not encryptable",
			args:   args{files: unreadableFile(), password: "test"},
			err:    errors.New("This file can't be encrypted"),
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.files, tt.args.password)
			assert.IsType(t, tt.err, err)
			assert.Equal(t, tt.result, got != nil)
		})
	}
}

func unreadableFile() []byte {
	return []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x34, 0x25, 0x25, 0x45, 0x4F, 0x46, 0x0A}
}

func loadFile(t *testing.T) []byte {
	file, err := os.ReadFile("../../testdata/file1.pdf")
	if err != nil {
		t.Fatal("testdata/file1.pdf not found")
	}

	return file
}

func loadFiles(includeInvalid bool, t *testing.T) [][]byte {
	files := [][]byte{loadFile(t), loadFile(t)}

	if includeInvalid == true {
		// fileInvalid := []byte(`invalid-file`)
		files = append(files, unreadableFile())
	}

	return files
}
