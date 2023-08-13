package pdf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/one-eyed-willy/pkg/logger"
	"github.com/one-eyed-willy/pkg/utils"
	pdfcpuAPI "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpuLog "github.com/pdfcpu/pdfcpu/pkg/log"
	pdfcpuConfig "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func sanitizeHTML(html string) string {
	p := bluemonday.UGCPolicy()

	// The policy can then be used to sanitize lots of input and it is safe to use the policy in multiple goroutines
	return p.Sanitize(html)
}

type PdfRender struct {
	Context context.Context
	Cancel  context.CancelFunc
}

func NewRender() *PdfRender {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("blink-settings", "scriptEnabled=false"),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	return &PdfRender{Context: ctx, Cancel: cancel}
}

func GenerateFromHTML(ctx context.Context, html string) (io.Reader, error) {
	ctx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(logger.Log().Infof))
	defer cancel()
	defer func() {
		_ = chromedp.Cancel(ctx)
	}()

	buf := bytes.Buffer{}
	if err := chromedp.Run(ctx, printHTMLToPDF(sanitizeHTML(html), &buf)); err != nil {
		return nil, err
	}

	return io.Reader(&buf), nil

}

func printHTMLToPDF(html string, res *bytes.Buffer) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = *bytes.NewBuffer(buf)
			return nil
		}),
	}
}

func Merge(files [][]byte) (file io.Reader, err error) {
	pdfcpuConfig.ConfigPath = "disable"
	pdfcpuLog.DisableLoggers()

	conf := pdfcpuConfig.NewDefaultConfiguration()
	var rs []io.ReadSeeker
	for i, res := range files {
		if err := ValidatePdf(res); err != nil {
			return nil, fmt.Errorf("File number %d is an invalid pdf: %s", i+1, err.Error())
		}
		rs = append(rs, io.ReadSeeker(bytes.NewReader(res)))
	}
	buf := bytes.Buffer{}

	if err = pdfcpuAPI.MergeRaw(rs, &buf, conf); err != nil {
		return nil, errors.New("Could not merge pdfs. Some files can't be read")
	}

	return io.Reader(&buf), nil
}

func ValidatePdf(data []byte) error {
	isInitialPdfBytes := utils.IsSubSlice(data[:5], []byte{0x25, 0x50, 0x44, 0x46, 0x2D})

	if !isInitialPdfBytes {
		return errors.New("This is not a pdf file")
	}

	if isVersion1dot3(data) {
		eof1dot3 := []byte{
			0x25, // %
			0x25, // %
			0x45, // E
			0x4F, // O
			0x46, // F
			0x20, // SPACE
			0x0A, // EOL
		}

		if utils.IsSubSlice(data[len(data)-7:], eof1dot3) {
			return nil
		}
		return errors.New("Invalid file terminator pdf v1.3")
	}

	if isVersion1dot4(data) {
		eof1dot4 := []byte{
			0x25, // %
			0x25, // %
			0x45, // E
			0x4F, // O
			0x46, // F
			0x0A, // EOL
		}

		if utils.IsSubSlice(data[len(data)-6:], eof1dot4) {
			return nil
		}
		return errors.New("Invalid file terminator for pdf v1.4")
	}

	return nil
}

func isVersion1dot3(data []byte) bool {
	return utils.IsSubSlice(data[5:8], []byte{0x31, 0x2E, 0x33})
}

func isVersion1dot4(data []byte) bool {
	return utils.IsSubSlice(data[5:8], []byte{0x31, 0x2E, 0x34})
}
