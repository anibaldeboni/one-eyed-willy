package pdf

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePdf(t *testing.T) {
	errInvalidPdf := errors.New("This is not a pdf file")
	type args struct {
		pdf []byte
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "When pdf is valid",
			args: args{pdf: []byte("%PDF-1.5")},
			want: nil,
		},
		{
			name: "When pdf is invalid",
			args: args{pdf: []byte("invalid-pdf")},
			want: errInvalidPdf,
		},
		{
			name: "When pdf terminator is invalid for v1.3",
			args: args{pdf: []byte("%PDF-1.3###EOF")},
			want: errors.New("Invalid file terminator pdf v1.3"),
		},
		{
			name: "When pdf terminator is invalid for v1.4",
			args: args{pdf: []byte("%PDF-1.4###EOF")},
			want: errors.New("Invalid file terminator for pdf v1.4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePdf(tt.args.pdf)
			assert.IsType(t, tt.want, got)
		})
	}
}
