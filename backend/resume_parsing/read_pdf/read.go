package readpdf

import (
	"bytes"

	"github.com/dslipak/pdf"
)

func ReadPdf(path string) (string, error) {
	r, err := pdf.Open(path)

	if err != nil {
		return "", err

	}

	var buf bytes.Buffer

	b, err := r.GetPlainText()

	if err != nil {
		return "", err
	}

	buf.ReadFrom(b)

	return buf.String(), nil
}
