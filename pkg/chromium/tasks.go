package chromium

import (
	"bufio"
	"bytes"
	"context"
	"fmt"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

func PrintHTMLToPDFTasks(html string, res *bytes.Buffer, options Options) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		fetch.Enable(),
		runtime.Enable(),
		chromedp.Navigate("about:blank"),
		loadHTMLAction(html),
		printHTMLToPDFAction(res, options),
	}
}

func loadHTMLAction(html string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		frameTree, err := page.GetFrameTree().Do(ctx)
		if err != nil {
			return err
		}

		return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
	}
}

func printHTMLToPDFAction(res *bytes.Buffer, options Options) chromedp.ActionFunc {
	return func(ctx context.Context) error {

		printToPdf := page.PrintToPDF().
			WithTransferMode(page.PrintToPDFTransferModeReturnAsStream).
			WithLandscape(options.Landscape).
			WithPrintBackground(options.PrintBackground).
			WithScale(options.Scale).
			WithPaperWidth(options.PaperWidth).
			WithPaperHeight(options.PaperHeight).
			WithMarginTop(options.MarginTop).
			WithMarginBottom(options.MarginBottom).
			WithMarginLeft(options.MarginLeft).
			WithMarginRight(options.MarginRight).
			WithPageRanges(options.PageRanges).
			WithPreferCSSPageSize(options.PreferCssPageSize)

		hasCustomHeaderFooter := options.HeaderTemplate != DefaultOptions().HeaderTemplate ||
			options.FooterTemplate != DefaultOptions().FooterTemplate

		if !hasCustomHeaderFooter {
			printToPdf = printToPdf.WithDisplayHeaderFooter(false)
		} else {
			printToPdf = printToPdf.
				WithDisplayHeaderFooter(true).
				WithHeaderTemplate(options.HeaderTemplate).
				WithFooterTemplate(options.FooterTemplate)
		}

		_, stream, err := printToPdf.Do(ctx)
		if err != nil {
			return fmt.Errorf("print to Pdf: %w", err)
		}

		reader := &streamReader{
			ctx:    ctx,
			handle: stream,
			r:      nil,
			pos:    0,
			eof:    false,
		}

		defer func() {
			if e := reader.Close(); e != nil {
				err = fmt.Errorf("close reader: %w", e)
			}
		}()

		buffer := bufio.NewReader(reader)

		_, err = buffer.WriteTo(res)

		if err != nil {
			return fmt.Errorf("write result to output: %w", err)
		}

		return nil
	}
}

func forceExactColorsActionFunc() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		// See:
		// https://github.com/gotenberg/gotenberg/issues/354
		// https://github.com/puppeteer/puppeteer/issues/2685
		// https://github.com/chromedp/chromedp/issues/520
		script := `
(() => {
	const css = 'html { -webkit-print-color-adjust: exact !important; }';

	const style = document.createElement('style');
	style.type = 'text/css';
	style.appendChild(document.createTextNode(css));
	document.head.appendChild(style);
})();
`

		evaluate := chromedp.Evaluate(script, nil)
		err := evaluate.Do(ctx)

		if err == nil {
			return nil
		}

		return fmt.Errorf("add CSS for exact colors: %w", err)
	}
}

func navigateActionFunc(logger *zap.Logger, url string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		logger.Debug(fmt.Sprintf("navigate to '%s'", url))

		_, _, _, err := page.Navigate(url).Do(ctx)
		if err != nil {
			return fmt.Errorf("navigate to '%s': %w", url, err)
		}

		err = runBatch(
			ctx,
			waitForEventDomContentEventFired(ctx, logger),
			waitForEventLoadEventFired(ctx, logger),
			waitForEventNetworkIdle(ctx, logger),
			waitForEventLoadingFinished(ctx, logger),
		)

		if err == nil {
			return nil
		}

		return fmt.Errorf("wait for events: %w", err)
	}
}
