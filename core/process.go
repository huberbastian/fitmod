package core

import (
	"fmt"
)

// Request contains all the necessary information to decode, process and re-encode a FIT file.
type Request struct {
	InputPath   string
	OutputPath  string
	Mode        Mode
	DistanceKm  float64
	AvgSpeedKmh float64
}

// Receives a Request and processes the FIT file accordingly.
func Process(req Request) error {
	rr := new(RecordRetriever)

	fit, err := decodeFile(req.InputPath, rr)
	if err != nil {
		return fmt.Errorf("error while decoding the FIT file '%s': %w", req.InputPath, err)
	}

	// Avoid by zero division during distance series generation.
	if len(rr.Records) == 0 {
		return fmt.Errorf("no record messages found in FIT file '%s'", req.InputPath)
	}

	distances, err := buildDistanceSeries(req.Mode, req.DistanceKm, req.AvgSpeedKmh, rr.Records)
	if err != nil {
		return err
	}
	injectDistanceValues(fit, distances)

	if err := encodeFile(fit, req.OutputPath); err != nil {
		return fmt.Errorf("error while encoding the FIT file '%s': %w", req.OutputPath, err)
	}

	return nil
}
