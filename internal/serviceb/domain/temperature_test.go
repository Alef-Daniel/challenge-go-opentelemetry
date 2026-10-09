package domain

import "testing"

func TestNewTemperatureFromCelsius(t *testing.T) {
	tests := []struct {
		celsius float64
		want    Temperature
	}{
		{celsius: 0, want: Temperature{Celsius: 0, Fahrenheit: 32, Kelvin: 273}},
		{celsius: 28.5, want: Temperature{Celsius: 28.5, Fahrenheit: 83.3, Kelvin: 301.5}},
		{celsius: -10, want: Temperature{Celsius: -10, Fahrenheit: 14, Kelvin: 263}},
		{celsius: 21.37, want: Temperature{Celsius: 21.37, Fahrenheit: 70.47, Kelvin: 294.37}},
	}

	for _, tt := range tests {
		if got := NewTemperatureFromCelsius(tt.celsius); got != tt.want {
			t.Errorf("NewTemperatureFromCelsius(%v) = %+v, want %+v", tt.celsius, got, tt.want)
		}
	}
}
