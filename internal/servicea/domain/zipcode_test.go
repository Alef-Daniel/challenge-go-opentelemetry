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
		{name: "valid", value: "29902555"},
		{name: "too short", value: "2990255", wantErr: ErrInvalidZipcode},
		{name: "too long", value: "299025555", wantErr: ErrInvalidZipcode},
		{name: "with letters", value: "2990255a", wantErr: ErrInvalidZipcode},
		{name: "empty", value: "", wantErr: ErrInvalidZipcode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewZipcode(tt.value); !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
