# fitmod

A CLI tool for injecting distance data into FIT activity files which lack GPS data.

Currently there are two ways of generating distance information:

- by directly specifying the total distance in kilometers,
- or by specifying the average speed in kilometers per hour.

The tool calculates the distance covered per record and updates each
records `distance` field accordingly.

For details about the FIT protocol, see the
[official Garmin FIT documentation](https://developer.garmin.com/fit/overview/).

## Motivation

Many stationary bikes or treadmills measure distance and speed locally,
but do not transmit this data to the device recording the activity.
As a result, the uploaded FIT activity may contain no usable distance
information in services such as Garmin Connect or Strava.

`fitmod` provides a simple way to add this missing distance data to
existing FIT activity files.

## Installation

### Go
```bash
go install github.com/huberbastian/fitmod@latest
```

## Usage

`fitmod` currently supports two modes. The selected mode depends on
which command line flag is provided.

### Distance Mode

Generate distance data based on a target total distance using the
`--distance` flag.

You can optionally specify an output path using `--output`.
If omitted, the output file will be written as `<input>_modified.fit`.

Example:
```bash
fitmod process <input.fit> --distance <total-distance> --output <output.fit>
```

To inject a total distance of 50km into `Indoor-Ride.fit`:
```bash
fitmod process Indoor-Ride.fit --distance 50
```

### Speed Mode

Generate distance data based on an average speed using the `--speed` flag.

The total distance is calculated from:

- the average speed,
- and the total activity duration derived from the FIT record timestamps.

Example:
```bash
fitmod process Indoor-Ride.fit --speed 25
```

#### If in doubt, run:
```bash
fitmod help
```