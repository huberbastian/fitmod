package core

import (
	"fmt"
)

// Represents to mode used to generate the distance series.
// Currently supports the following:
// - ModeDistance: generates a distance series from a target distance value.
// - ModeSpeed: generates a distance series based on an average speed value.
type Mode int

const (
	ModeDistance Mode = iota
	ModeSpeed
)

func buildDistanceSeries(mode Mode, km float64, kmh float64, records []Record) ([]float64, error) {
	switch mode {

	case ModeDistance:
		return fromTotal(toMeters(km), len(records))

	case ModeSpeed:
		return fromSpeed(toMetersPerSecond(kmh), records)

	default:
		return nil, fmt.Errorf("unsupported mode: %v", mode)
	}
}

// Generates a distance series from a target total distance value.
// Divides the total distance by the number of records to determine how much distance
// needs to be covered per record.
// Expects a distance in meters.
func fromTotal(distance float64, n int) ([]float64, error) {
	distances := make([]float64, n)
	step := distance / float64(n-1)
	for i := range n {
		distances[i] = step * float64(i)
	}
	distances[n-1] = distance

	return distances, nil
}

// Generates a distance series based on an average speed value.
// Calculates total time using the Record Timestamps and computes the total distance.
// Then calls fromTotal to generate the distance series.
// Expects a speed in meters per second.
func fromSpeed(speed float64, records []Record) ([]float64, error) {
	firstTs := records[0].Timestamp
	lastTs := records[len(records)-1].Timestamp

	time := float64(lastTs - firstTs)
	distance := speed * time

	return fromTotal(distance, len(records))
}

func toMeters(km float64) float64 {
	return km * 1000
}

func toMetersPerSecond(kmh float64) float64 {
	return kmh * 1000 / 3600
}
