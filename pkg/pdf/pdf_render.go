package pdf

import (
	"bytes"
	"context"
	"io"

	"github.com/chromedp/chromedp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/one-eyed-willy/pkg/chromium"
	"go.uber.org/zap"
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

func NewRender(logger *zap.Logger) *PdfRender {
	r := &PdfRender{chromeApi: &chromeApiImpl{}}
	ctx, cancel := chromium.New(logger)

	r.Context = ctx
	r.Cancel = cancel
	return r
}

func (p *PdfRender) GenerateFromHTML(html string) (io.Reader, error) {
	ctx, cancel := p.chromeApi.NewContext(p.Context)
	defer cancel()

	buf := bytes.Buffer{}

	options := chromium.DefaultOptions()
	if err := p.chromeApi.Run(ctx, chromium.PrintHTMLToPDFTasks(sanitizeHTML(html), &buf, options)); err != nil {
		return nil, err
	}

	return io.Reader(&buf), nil

}
