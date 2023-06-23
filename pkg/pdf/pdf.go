package pdf

import (
	"strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func GenerateFromHtml(html string) ([]byte, error) {
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(html)))

	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
