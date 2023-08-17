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
	pdfcpuAPI "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpuLog "github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
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

	c, cancelCtx := chromedp.NewContext(ctx)

	// ensure cleanup
	var success bool
	defer func() {
		if !success {
			_ = chromedp.Cancel(c)
			cancelCtx()
		}
	}()

	err := chromedp.Run(c)
	if err != nil {
		panic(err)
	}
	success = true

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

func Encrypt(file []byte, password string) (io.Reader, error) {
	pdfcpuConfig.ConfigPath = "disable"
	pdfcpuLog.DisableLoggers()

	if err := ValidatePdf(file); err != nil {
		return nil, err
	}
	conf := model.NewAESConfiguration(password, password, 256)

	buf := bytes.Buffer{}
	if err := pdfcpuAPI.Encrypt(bytes.NewReader(file), &buf, conf); err != nil {
		return nil, errors.New("Could not encrypt pdf")
	}

	return io.Reader(&buf), nil
}
