package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	pdfcpuAPI "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpuLog "github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	pdfcpuConfig "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PdfTool struct {
	pdfapi pdfApi
}

type pdfApi interface {
	Encrypt(rs io.ReadSeeker, w io.Writer, conf *pdfcpuConfig.Configuration) error
	Merge(rs []io.ReadSeeker, w io.Writer, conf *pdfcpuConfig.Configuration) error
}

type pdfApiImpl struct{}

func (p *pdfApiImpl) Encrypt(rs io.ReadSeeker, w io.Writer, conf *pdfcpuConfig.Configuration) error {
	return pdfcpuAPI.Encrypt(rs, w, conf)
}

func (p *pdfApiImpl) Merge(rs []io.ReadSeeker, w io.Writer, conf *pdfcpuConfig.Configuration) error {
	return pdfcpuAPI.MergeRaw(rs, w, conf)
}

func NewPdfTool() *PdfTool {
	return &PdfTool{pdfapi: &pdfApiImpl{}}
}

func (p *PdfTool) Merge(files [][]byte) (file io.Reader, err error) {
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

	if err = p.pdfapi.Merge(rs, &buf, conf); err != nil {
		return nil, errors.New("Could not merge pdfs. Some files can't be read")
	}

	return io.Reader(&buf), nil
}

func (p *PdfTool) Encrypt(file []byte, password string) (io.Reader, error) {
	pdfcpuConfig.ConfigPath = "disable"
	pdfcpuLog.DisableLoggers()

	if err := ValidatePdf(file); err != nil {
		return nil, err
	}
	conf := model.NewAESConfiguration(password, password, 256)

	buf := bytes.Buffer{}
	if err := p.pdfapi.Encrypt(bytes.NewReader(file), &buf, conf); err != nil {
		return nil, errors.New("Could not encrypt pdf")
	}

	return io.Reader(&buf), nil
}
