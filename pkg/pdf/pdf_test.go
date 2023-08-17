package pdf

import (
	"errors"
	"io"
	"testing"

	"github.com/one-eyed-willy/testdata"
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
			args:          args{files: testdata.LoadFilesWithInvalid(false, t)},
			resultIsEmpty: false,
			err:           nil,
		},
		{
			name:          "When some file is invalid",
			args:          args{files: testdata.LoadFilesWithInvalid(true, t)},
			resultIsEmpty: true,
			err:           errors.New("Could not merge pdfs. Some files can't be read"),
		},
		{
			name:          "When the files can't be merged",
			args:          args{files: [][]byte{testdata.UnreadableFile(), testdata.UnreadableFile()}},
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
			args:   args{files: testdata.LoadFilesWithInvalid(false, t)[0], password: "test"},
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
			args:   args{files: testdata.UnreadableFile(), password: "test"},
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
