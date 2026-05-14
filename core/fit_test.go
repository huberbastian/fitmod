package core

import (
	"testing"

	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
)

func TestRecordRetrieverOnMesg(t *testing.T) {
	rr := new(RecordRetriever)

	recordMsg := proto.Message{
		Num: mesgnum.Record,
		Fields: []proto.Field{
			{
				FieldBase: &proto.FieldBase{
					Name: "timestamp",
				},
				Value: proto.Any(uint32(100)),
			},
		},
	}

	otherMsg := proto.Message{
		Num: 999,
	}

	rr.OnMesg(recordMsg)
	rr.OnMesg(otherMsg)

	if len(rr.Records) != 1 {
		t.Fatalf("got %d records, want 1", len(rr.Records))
	}

	if rr.Records[0].Timestamp != 100 {
		t.Fatalf("unexpected timestamp")
	}
}

func TestToRecord(t *testing.T) {
	msg := proto.Message{
		Num: mesgnum.Record,
		Fields: []proto.Field{
			{
				FieldBase: &proto.FieldBase{
					Name: "timestamp",
				},
				Value: proto.Any(uint32(12345)),
			},
		},
	}

	got := toRecord(msg)

	if got.Timestamp != 12345 {
		t.Fatalf("got timestamp %d, want %d", got.Timestamp, 12345)
	}
}

func TestInjectDistanceValues(t *testing.T) {
	fit := proto.FIT{
		Messages: []proto.Message{
			{
				Num: mesgnum.Record,
				Fields: []proto.Field{
					{
						FieldBase: &proto.FieldBase{
							Name: "distance",
						},
						Value: proto.Any(float64(0)),
					},
				},
			},
			{
				Num: mesgnum.Record,
				Fields: []proto.Field{
					{
						FieldBase: &proto.FieldBase{
							Name: "distance",
						},
						Value: proto.Any(float64(0)),
					},
				},
			},
		},
	}

	distances := []float64{100, 200}

	injectDistanceValues(&fit, distances)

	got0 := fit.Messages[0].Fields[0].Value.Float64()
	got1 := fit.Messages[1].Fields[0].Value.Float64()

	if got0 != 100 {
		t.Fatalf("got %f want 100", got0)
	}

	if got1 != 200 {
		t.Fatalf("got %f want 200", got1)
	}
}
