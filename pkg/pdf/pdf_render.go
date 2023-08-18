package pdf

import (
	"bytes"
	"context"
	"io"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/one-eyed-willy/pkg/logger"
)

func sanitizeHTML(html string) string {
	p := bluemonday.UGCPolicy()

	// The policy can then be used to sanitize lots of input and it is safe to use the policy in multiple goroutines
	return p.Sanitize(html)
}

type PdfRender struct {
	Context   context.Context
	Cancel    context.CancelFunc
	chromeApi chromeApi
}
type chromeApi interface {
	Run(ctx context.Context, actions ...chromedp.Action) error
	NewContext(parent context.Context, opts ...func(*chromedp.Context)) (context.Context, context.CancelFunc)
}
type chromeApiImpl struct{}

func (c *chromeApiImpl) Run(ctx context.Context, actions ...chromedp.Action) error {
	return chromedp.Run(ctx, actions...)
}

func (c *chromeApiImpl) NewContext(parent context.Context, opts ...func(*chromedp.Context)) (context.Context, context.CancelFunc) {
	return chromedp.NewContext(parent, opts...)
}

func NewRender() *PdfRender {
	r := &PdfRender{chromeApi: &chromeApiImpl{}}
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

	r.Context = ctx
	r.Cancel = cancel
	return r
}

func (p *PdfRender) GenerateFromHTML(html string) (io.Reader, error) {
	ctx, cancel := p.chromeApi.NewContext(p.Context, chromedp.WithLogf(logger.Log().Infof))
	defer cancel()
	defer func() {
		_ = chromedp.Cancel(ctx)
	}()

	buf := bytes.Buffer{}
	if err := p.chromeApi.Run(ctx, printHTMLToPDF(sanitizeHTML(html), &buf)); err != nil {
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
