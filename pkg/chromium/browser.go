package chromium

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

func New(logger *zap.Logger) (context.Context, context.CancelFunc) {

	debug := &debugLogger{logger: logger}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.CombinedOutput(debug),
		chromedp.DisableGPU,
		chromedp.Headless,
		chromedp.Flag("blink-settings", "scriptEnabled=false"),
		chromedp.NoSandbox,
	)
	allocatorContext, allocatorCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, ctxCancel := chromedp.NewContext(allocatorContext)

	err := chromedp.Run(ctx)
	if err != nil {
		fmt.Println("could not start chromium: %w", err)
		ctxCancel()
		allocatorCancel()
		panic(err)
	}

	cancel := func() {
		allocatorCancel()
		ctxCancel()
	}

	return ctx, cancel
}
