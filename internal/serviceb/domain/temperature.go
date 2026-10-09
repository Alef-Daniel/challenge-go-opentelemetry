package domain

import "math"

type Temperature struct {
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}

func NewTemperatureFromCelsius(celsius float64) Temperature {
	return Temperature{
		Celsius:    round(celsius),
		Fahrenheit: round(celsius*1.8 + 32),
		Kelvin:     round(celsius + 273),
	}
}

func round(value float64) float64 {
	return math.Round(value*100) / 100
}
