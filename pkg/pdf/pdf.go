package pdf

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/microcosm-cc/bluemonday"
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
	for _, res := range files {
		rs = append(rs, io.ReadSeeker(bytes.NewReader(res)))
	}
	buf := bytes.Buffer{}

	if err = pdfcpuAPI.MergeRaw(rs, &buf, conf); err != nil {
		return nil, errors.New("Could not merge pdfs. Some files can't be read")
	}
	return buf.Bytes(), nil
}
