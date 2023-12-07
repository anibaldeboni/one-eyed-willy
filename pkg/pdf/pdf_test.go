package pdf

import (
	"testing"

	"github.com/one-eyed-willy/testdata"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var logger = zap.NewNop()

func TestGenerateFromHTML(t *testing.T) {
	pdfRender := NewRender(logger)
	mockChromeApi := new(MockChromeApi)

	type args struct {
		html string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		pdfRender *PdfRender
	}{
		{
			name:      "When html is valid",
			args:      args{html: "<h1>Hello World</h1>"},
			wantErr:   false,
			pdfRender: pdfRender,
		},
		{
			name:      "When html is empty",
			args:      args{html: ""},
			wantErr:   false,
			pdfRender: pdfRender,
		},
		{
			name:      "When chromedp returns an error",
			args:      args{html: "<h1>Hello World</h1>"},
			wantErr:   true,
			pdfRender: &PdfRender{chromeApi: mockChromeApi},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.pdfRender.GenerateFromHTML(tt.args.html)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantErr, result == nil)
		})
	}
}

func BenchmarkGenerateFromHTML(b *testing.B) {
	pdfRender := NewRender(logger)
	for i := 0; i < b.N; i++ {
		_, _ = pdfRender.GenerateFromHTML("<h1>Hello World</h1>")
	}
}

func TestMerge(t *testing.T) {
	type args struct {
		files [][]byte
	}
	tests := []struct {
		name    string
		args    args
		wantGot bool
		wantErr bool
	}{
		{
			name:    "When all files could be read",
			args:    args{files: testdata.LoadFilesWithInvalid(false, t)},
			wantGot: true,
			wantErr: false,
		},
		{
			name:    "When some file is invalid",
			args:    args{files: testdata.LoadFilesWithInvalid(true, t)},
			wantGot: false,
			wantErr: true,
		},
		{
			name:    "When the files can't be merged",
			args:    args{files: [][]byte{testdata.UnreadableFile(), testdata.UnreadableFile()}},
			wantGot: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		pdf := &PdfTool{pdfapi: &pdfApiImpl{}}
		t.Run(tt.name, func(t *testing.T) {
			result, err := pdf.Merge(tt.args.files)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantGot, result != nil)
		})
	}
}

func BenchmarkMerge(b *testing.B) {
	pdf := &PdfTool{pdfapi: &pdfApiImpl{}}
	for i := 0; i < b.N; i++ {
		_, _ = pdf.Merge(testdata.LoadFilesWithInvalid(false, nil))
	}
}

func TestEncrypt(t *testing.T) {
	pdfTool := &PdfTool{pdfapi: &pdfApiImpl{}}
	mockPdfApi := new(MockPdfApi)
	type args struct {
		files    []byte
		password string
	}
	tests := []struct {
		name    string
		args    args
		pdfTool *PdfTool
		wantErr bool
		wantGot bool
	}{
		{
			name:    "When the file is valid",
			args:    args{files: testdata.LoadFilesWithInvalid(false, t)[0], password: "test"},
			pdfTool: pdfTool,
			wantErr: false,
			wantGot: true,
		},
		{
			name:    "When the file is invalid",
			args:    args{files: []byte(`not-a-pdf`), password: "test"},
			pdfTool: pdfTool,
			wantErr: true,
			wantGot: false,
		},
		{
			name:    "When the file is not encryptable",
			args:    args{files: testdata.UnreadableFile(), password: "test"},
			pdfTool: &PdfTool{pdfapi: mockPdfApi},
			wantErr: true,
			wantGot: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pdfTool.Encrypt(tt.args.files, tt.args.password)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantGot, got != nil)
		})
	}
}

func BenchmarkEncrypt(b *testing.B) {
	pdf := &PdfTool{pdfapi: &pdfApiImpl{}}
	for i := 0; i < b.N; i++ {
		_, _ = pdf.Encrypt(testdata.LoadFilesWithInvalid(false, nil)[0], "test")
	}
}
