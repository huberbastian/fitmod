package processor

import (
	"path/filepath"
	"testing"

	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
)

func TestProcess_DistanceMode(t *testing.T) {
	input := "../testdata/sample.fit"

	output := filepath.Join(
		t.TempDir(),
		"output.fit",
	)

	req := Request{
		InputPath:  input,
		OutputPath: output,
		Mode:       ModeDistance,
		DistanceKm: 54,
	}

	err := Process(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rr := new(RecordRetriever)

	fit, err := decodeFile(output, rr)
	if err != nil {
		t.Fatalf("failed to decode output file: %v", err)
	}

	if len(rr.Records) == 0 {
		t.Fatal("expected records in processed FIT file")
	}

	lastDistance := getLastDistance(fit)

	if !almostEqual(lastDistance, 54000) {
		t.Fatalf(
			"last distance = %f, want %f",
			lastDistance,
			54000.00,
		)
	}
}

func TestProcess_SpeedMode(t *testing.T) {
	input := "../testdata/sample.fit"

	output := filepath.Join(
		t.TempDir(),
		"output.fit",
	)

	req := Request{
		InputPath:   input,
		OutputPath:  output,
		Mode:        ModeSpeed,
		AvgSpeedKmh: 25,
	}

	err := Process(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rr := new(RecordRetriever)

	fit, err := decodeFile(output, rr)
	if err != nil {
		t.Fatalf("failed to decode output file: %v", err)
	}

	if len(rr.Records) == 0 {
		t.Fatal("expected records in processed FIT file")
	}

	lastDistance := getLastDistance(fit)

	if !almostEqual(lastDistance, 20847.22) {
		t.Fatalf(
			"last distance = %f, want %f",
			lastDistance,
			20847.22,
		)
	}
}

func getLastDistance(fit *proto.FIT) float64 {
	var last float64

	for _, msg := range fit.Messages {
		if msg.Num != mesgnum.Record {
			continue
		}

		for _, f := range msg.Fields {
			if f.Name == "distance" {
				// Wee need to apply the scale and offset according to the SDK specification,
				// since the go SDK does not provide an accessor for the scaled value.
				raw := f.Value.Uint32()
				last = (float64(raw) / f.Scale) - f.Offset
			}
		}
	}

	return last
}
