package processor

import (
	"math"
	"testing"
)

func TestBuildDistanceSeries(t *testing.T) {
	records := []Record{
		{Timestamp: 0},
		{Timestamp: 10},
		{Timestamp: 20},
		{Timestamp: 30},
	}

	tests := []struct {
		name     string
		mode     Mode
		km       float64
		kmh      float64
		wantLast float64
		wantErr  bool
	}{
		{
			name:     "distance mode",
			mode:     ModeDistance,
			km:       76,
			wantLast: 76000,
		},
		{
			name:     "speed mode",
			mode:     ModeSpeed,
			kmh:      36,  // 10 m/s
			wantLast: 300, // 10 m/s for 30 seconds
		},
		{
			name:    "invalid mode",
			mode:    Mode(999),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildDistanceSeries(tt.mode, tt.km, tt.kmh, records)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(records) {
				t.Fatalf("got %d distances, want %d", len(got), len(records))
			}

			last := got[len(got)-1]

			if !almostEqual(last, tt.wantLast) {
				t.Fatalf("last distance = %f, want %f", last, tt.wantLast)
			}
		})
	}
}

func TestFromTotal(t *testing.T) {
	tests := []struct {
		name     string
		distance float64
		n        int
		want     []float64
	}{
		{
			name:     "1000m over 5 records",
			distance: 1000,
			n:        5,
			want:     []float64{0, 250, 500, 750, 1000},
		},
		{
			name:     "12km over 3 records",
			distance: 12000,
			n:        3,
			want:     []float64{0, 6000, 12000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromTotal(tt.distance, tt.n)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got len=%d want len=%d", len(got), len(tt.want))
			}

			for i := range got {
				if !almostEqual(got[i], tt.want[i]) {
					t.Fatalf("index %d = %f, want %f", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestFromSpeed(t *testing.T) {
	records := []Record{
		{Timestamp: 0},
		{Timestamp: 10},
		{Timestamp: 20},
		{Timestamp: 30},
	}

	tests := []struct {
		name     string
		speed    float64
		wantLast float64
	}{
		{
			name:     "10 m/s for 30 seconds",
			speed:    10,
			wantLast: 300,
		},
		{
			name:     "5 m/s for 30 seconds",
			speed:    5,
			wantLast: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromSpeed(tt.speed, records)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			last := got[len(got)-1]

			if !almostEqual(last, tt.wantLast) {
				t.Fatalf("last distance = %f, want %f", last, tt.wantLast)
			}
		})
	}
}

func TestToMeters(t *testing.T) {
	tests := []struct {
		km   float64
		want float64
	}{
		{1, 1000},
		{12, 12000},
		{0.5, 500},
	}

	for _, tt := range tests {
		got := toMeters(tt.km)

		if got != tt.want {
			t.Fatalf("got %f want %f", got, tt.want)
		}
	}
}

func TestToMetersPerSecond(t *testing.T) {
	tests := []struct {
		kmh  float64
		want float64
	}{
		{45, 12.5},
		{18, 5},
		{26.5, 7.361111},
		{22, 6.111111},
	}

	for _, tt := range tests {
		got := toMetersPerSecond(tt.kmh)

		if !almostEqual(got, tt.want) {
			t.Fatalf("got %f want %f", got, tt.want)
		}
	}
}

// tolerance for floating point comparisons
func almostEqual(a, b float64) bool {
	const epsilon = 0.000001
	return math.Abs(a-b) < epsilon
}
