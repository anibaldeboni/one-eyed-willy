package pdf

import (
	"bytes"
	"context"
	"io"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
)

func sanitizeHtml(html string) string {
	p := bluemonday.UGCPolicy()

	// The policy can then be used to sanitize lots of input and it is safe to use the policy in multiple goroutines
	return p.Sanitize(html)
}

func GenerateFromHtml(html string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, printHtmlToPDF(sanitizeHtml(html), &buf)); err != nil {
		return nil, err
	}

	return buf, nil

}

func printHtmlToPDF(html string, res *[]byte) chromedp.Tasks {
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

func MergePdfs(files [][]byte) ([]byte, error) {
	pdf := gofpdf.New("P", "pt", "A4", "")

	imp := gofpdi.NewImporter()

	for _, res := range files {
		rs := io.ReadSeeker(bytes.NewReader(res))
		tpl := imp.ImportPageFromStream(pdf, &rs, 1, "/MediaBox")

		pageSizes := imp.GetPageSizes()
		numberOfPages := len(pageSizes)

		for i := 1; i <= numberOfPages; i++ {
			pdf.AddPage()

			if i > 1 {
				tpl = imp.ImportPageFromStream(pdf, &rs, i, "/MediaBox")
			}
			imp.UseImportedTemplate(pdf, tpl, 0, 0, pageSizes[i]["/MediaBox"]["w"], pageSizes[i]["/MediaBox"]["h"])
		}
	}

	buf := bytes.Buffer{}
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
