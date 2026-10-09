package domain

import (
	"errors"
	"testing"
)

func TestNewZipcode(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{name: "valid", value: "01001000"},
		{name: "too short", value: "0100100", wantErr: ErrInvalidZipcode},
		{name: "too long", value: "010010000", wantErr: ErrInvalidZipcode},
		{name: "with letters", value: "0100100a", wantErr: ErrInvalidZipcode},
		{name: "with dash", value: "01001-00", wantErr: ErrInvalidZipcode},
		{name: "empty", value: "", wantErr: ErrInvalidZipcode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zipcode, err := NewZipcode(tt.value)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil && zipcode.String() != tt.value {
				t.Fatalf("expected %q, got %q", tt.value, zipcode)
			}
		})
	}
}
