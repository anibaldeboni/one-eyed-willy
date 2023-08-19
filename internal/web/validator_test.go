package web

import (
	"testing"
)

type validation struct {
	Field string `json:"field" validate:"required,gt=0"`
}

func TestValidator(t *testing.T) {
	v := NewValidator()
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "When the value is empty",
			value:   "",
			wantErr: true,
		},
		{
			name:    "When the value is not empty",
			value:   "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(validation{Field: tt.value})
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
