package core

import (
	"os"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/encoder"
	"github.com/muktihari/fit/profile/untyped/mesgnum"
	"github.com/muktihari/fit/proto"
)

// Internal Record struct to represent a message of type 'record', containing only the fields needed to
// generate the distance series.
type Record struct {
	Timestamp uint32
}

// RecordRetriever implements the decoder.MesgListener interface to retrieve all 'record' messages
// from the FIT file and store them as internal Record structs.
type RecordRetriever struct{ Records []Record }

var _ decoder.MesgListener = (*RecordRetriever)(nil)

func (r *RecordRetriever) OnMesg(mesg proto.Message) {
	if mesg.Num == mesgnum.Record {
		r.Records = append(r.Records, toRecord(mesg))
	}
}

// Decodes the FIT file at the specified input path and returns the decoded proto.FIT struct.
// Stores all 'record' messages in the provided RecordRetriever.
func decodeFile(inputPath string, rr *RecordRetriever) (*proto.FIT, error) {
	in, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	dec := decoder.New(in,
		// Add listener which retrieves only record messages
		decoder.WithMesgListener(rr),
	)

	return dec.Decode()
}

// Encodes the modified FIT file and writes it to the specified output path.
func encodeFile(fit *proto.FIT, outputPath string) error {
	out, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()

	return encoder.New(out).Encode(fit)
}

// Mapping from proto.Message to internal Record struct.
func toRecord(m proto.Message) Record {
	record := Record{}
	for _, f := range m.Fields {
		if f.Name == "timestamp" {
			record.Timestamp = f.Value.Uint32()
		}
	}
	return record
}

// Updates the value of the distance field on each 'record' message in the FIT file
// according to the generated distance series.
func injectDistanceValues(fit *proto.FIT, distances []float64) {
	i := 0 // index for distance series

	for mi := range fit.Messages {
		msg := &fit.Messages[mi]
		if msg.Num != mesgnum.Record {
			continue
		}

		for fi := range msg.Fields {
			f := &msg.Fields[fi]

			if f.Name == "distance" {
				f.Value = proto.Any(distances[i])
			}
		}

		i++
	}
}
