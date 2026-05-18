package cmd

import (
	"testing"

	"github.com/huberbastian/fitmod/internal/processor"
	"github.com/spf13/cobra"
)

func TestDefaultOutputPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"activity.fit", "activity_modified.fit"},
		{"path/to/activity.fit", "path/to/activity_modified.fit"},
	}

	for _, tt := range tests {
		got := defaultOutputPath(tt.input)

		if got != tt.want {
			t.Errorf("defaultOutputPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		name       string
		inputPath  string
		outputPath string
		distance   float64
		speed      float64
		want       processor.Request
	}{
		{
			name:      "distance mode with default output",
			inputPath: "activity.fit",
			distance:  31,

			want: processor.Request{
				InputPath:   "activity.fit",
				OutputPath:  "activity_modified.fit",
				Mode:        processor.ModeDistance,
				DistanceKm:  31,
				AvgSpeedKmh: 0,
			},
		},
		{
			name:      "speed mode with default output",
			inputPath: "activity.fit",
			speed:     25,

			want: processor.Request{
				InputPath:   "activity.fit",
				OutputPath:  "activity_modified.fit",
				Mode:        processor.ModeSpeed,
				DistanceKm:  0,
				AvgSpeedKmh: 25,
			},
		},
		{
			name:       "custom output path",
			inputPath:  "activity.fit",
			outputPath: "custom.fit",
			distance:   20,

			want: processor.Request{
				InputPath:   "activity.fit",
				OutputPath:  "custom.fit",
				Mode:        processor.ModeDistance,
				DistanceKm:  20,
				AvgSpeedKmh: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			cmd.Flags().Float64("distance", tt.distance, "")
			cmd.Flags().Float64("speed", tt.speed, "")
			cmd.Flags().String("output", tt.outputPath, "")

			got := buildRequest(cmd, []string{tt.inputPath})

			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name        string
		distanceSet bool
		speedSet    bool
		distance    float64
		speed       float64
		wantErr     bool
	}{
		{
			name:    "missing both flags",
			wantErr: true,
		},
		{
			name:        "distance only",
			distanceSet: true,
			distance:    57,
			wantErr:     false,
		},
		{
			name:     "speed only",
			speedSet: true,
			speed:    36,
			wantErr:  false,
		},
		{
			name:        "both flags set",
			distanceSet: true,
			speedSet:    true,
			distance:    108,
			speed:       23,
			wantErr:     true,
		},
		{
			name:        "negative distance",
			distanceSet: true,
			distance:    -5,
			wantErr:     true,
		},
		{
			name:     "negative speed",
			speedSet: true,
			speed:    -16,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			cmd.Flags().Float64("distance", 0, "")
			cmd.Flags().Float64("speed", 0, "")

			// reset globals
			distanceKm = 0
			speedKmh = 0

			// Explicitly set the flags to make sure that cmd.Flags().Changed() returns true
			// for the respective flags during validation.
			if tt.distanceSet {
				distanceKm = tt.distance
				cmd.Flags().Set("distance", "1")
			}

			if tt.speedSet {
				speedKmh = tt.speed
				cmd.Flags().Set("speed", "1")
			}

			err := validateFlags(cmd, nil)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
