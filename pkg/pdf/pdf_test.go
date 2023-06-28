package pdf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergePdfFilesCaseSuccess(t *testing.T) {
	var err error

	file1, err := os.ReadFile("../../testdata/file1.pdf")
	if err != nil {
		t.Fatal("testdata/file1.pdf not found")
	}

	file2, err := os.ReadFile("../../testdata/file2.pdf")
	if err != nil {
		t.Fatal("testdata/file2.pdf not found")
	}

	pdf, err := MergePdfs([][]byte{file1, file2})
	assert.NoError(t, err)
	assert.Greater(t, len(pdf), 0)
}
