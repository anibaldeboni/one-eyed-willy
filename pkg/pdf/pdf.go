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

func GenerateFromHTML(html string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, printHTMLToPDF(sanitizeHTML(html), &buf)); err != nil {
		return nil, err
	}

	return buf, nil

}

func printHTMLToPDF(html string, res *[]byte) chromedp.Tasks {
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
			*res = buf
			return nil
		}),
	}
}

func Merge(files [][]byte) (file []byte, err error) {
	pdfcpuConfig.ConfigPath = "disable"
	pdfcpuLog.DisableLoggers()

	conf := pdfcpuConfig.NewDefaultConfiguration()
	var rs []io.ReadSeeker
	for i, res := range files {
		if err := IsPdf(res); err != nil {
			return nil, fmt.Errorf("File number %d is not a pdf", i)
		}
		rs = append(rs, io.ReadSeeker(bytes.NewReader(res)))
	}
	buf := bytes.Buffer{}

	if err = pdfcpuAPI.MergeRaw(rs, &buf, conf); err != nil {
		return nil, errors.New("Could not merge pdfs. Some files can't be read")
	}
	return buf.Bytes(), nil
}

func IsPdf(data []byte) error {
	isInitialPdfBytes := utils.IsByteSubSlice(data[:5], []byte{0x25, 0x50, 0x44, 0x46, 0x2D})

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

		if utils.IsByteSubSlice(data[len(data)-7:], eof1dot3) {
			return nil
		}
		return errors.New("Invalid file terminator")
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

		if utils.IsByteSubSlice(data[len(data)-6:], eof1dot4) {
			return nil
		}
		return errors.New("Invalid file terminator")
	}

	return nil
}

func isVersion1dot3(data []byte) bool {
	return utils.IsByteSubSlice(data[5:8], []byte{0x31, 0x2E, 0x33})
}

func isVersion1dot4(data []byte) bool {
	return utils.IsByteSubSlice(data[5:8], []byte{0x31, 0x2E, 0x34})
}
