package pdf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePdf(t *testing.T) {
	type args struct {
		pdf []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "When pdf is valid",
			args:    args{pdf: []byte("%PDF-1.5")},
			wantErr: false,
		},
		{
			name:    "When pdf is invalid",
			args:    args{pdf: []byte("invalid-pdf")},
			wantErr: true,
		},
		{
			name:    "When pdf terminator is invalid for v1.3",
			args:    args{pdf: []byte("%PDF-1.3###EOF")},
			wantErr: true,
		},
		{
			name:    "When pdf terminator is valid for v1.3",
			args:    args{pdf: []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x33, 0x25, 0x25, 0x45, 0x4F, 0x46, 0x20, 0x0A}},
			wantErr: false,
		},
		{
			name:    "When pdf terminator is invalid for v1.4",
			args:    args{pdf: []byte("%PDF-1.4###EOF")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePdf(tt.args.pdf)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
