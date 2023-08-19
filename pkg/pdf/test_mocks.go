package pdf

import (
	"context"
	"errors"
	"io"

	"github.com/chromedp/chromedp"
	pdfcpuConfig "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/stretchr/testify/mock"
)

type MockPdfApi struct {
	mock.Mock
}

func (m *MockPdfApi) Merge(rs []io.ReadSeeker, w io.Writer, conf *pdfcpuConfig.Configuration) error {
	return errors.New("Could not merge pdfs. Some files can't be read")
}
func (m *MockPdfApi) Encrypt(rs io.ReadSeeker, w io.Writer, conf *pdfcpuConfig.Configuration) error {
	return errors.New("Could not encrypt pdf")
}

func NewMockPdfRender() *PdfRender {
	return &PdfRender{chromeApi: new(MockChromeApi)}
}

type MockChromeApi struct {
	mock.Mock
}

func (m *MockChromeApi) Run(ctx context.Context, actions ...chromedp.Action) error {
	return errors.New("Could not run chrome")
}

func (m *MockChromeApi) NewContext(parent context.Context, opts ...func(*chromedp.Context)) (context.Context, context.CancelFunc) {
	return context.Background(), func() {}
}
