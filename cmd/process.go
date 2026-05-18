package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/huberbastian/fitmod/internal/processor"
	"github.com/spf13/cobra"
)

var (
	distanceKm float64
	speedKmh   float64
	outputPath string

	processCmd = &cobra.Command{
		Use:     "process <input.fit>",
		Short:   "Inject distance data into the specified FIT file",
		Aliases: []string{"p"},
		Args:    cobra.ExactArgs(1),
		PreRunE: validateFlags,
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]

			if outputPath == "" {
				outputPath = defaultOutputPath(inputPath)
			}

			if distanceKm > 0 {
				fmt.Printf("Injecting distance: %.2fkm\n", distanceKm)
			}
			if speedKmh > 0 {
				fmt.Printf("Injecting distance based on average speed: %.2fkm/h\n", speedKmh)
			}

			req := buildRequest(cmd, args)
			if err := processor.Process(req); err != nil {
				return fmt.Errorf("error processing FIT file: %w", err)
			}

			fmt.Printf("File processed successfully. Output written to '%s'\n", outputPath)
			return nil
		},
	}
)

func init() {
	processCmd.Flags().Float64VarP(&distanceKm, "distance", "d", 0, "target distance in kilometers")
	processCmd.Flags().Float64VarP(&speedKmh, "speed", "s", 0, "average speed in km/h")
	processCmd.Flags().StringVarP(&outputPath, "output", "o", "", "(optional) output path")
}

func validateFlags(cmd *cobra.Command, args []string) error {
	distanceSet := cmd.Flags().Changed("distance")
	speedSet := cmd.Flags().Changed("speed")

	if !distanceSet && !speedSet {
		return fmt.Errorf("at least one of --distance or --speed must be provided")
	}

	if distanceSet && speedSet {
		return fmt.Errorf("only one of --distance or --speed can be provided")
	}

	if distanceSet && distanceKm < 0 {
		return fmt.Errorf("--distance must be a positive value")
	}

	if speedSet && speedKmh < 0 {
		return fmt.Errorf("--speed must be a positive value")
	}
	return nil
}

func buildRequest(cmd *cobra.Command, args []string) processor.Request {
	input := args[0]

	output, _ := cmd.Flags().GetString("output")
	if output == "" {
		output = defaultOutputPath(input)
	}

	distance, _ := cmd.Flags().GetFloat64("distance")
	speed, _ := cmd.Flags().GetFloat64("speed")

	mode := processor.ModeDistance
	if speed > 0 {
		mode = processor.ModeSpeed
	}

	return processor.Request{
		InputPath:   input,
		OutputPath:  output,
		Mode:        mode,
		DistanceKm:  distance,
		AvgSpeedKmh: speed,
	}
}

func defaultOutputPath(input string) string {
	ext := filepath.Ext(input)
	base := strings.TrimSuffix(input, ext)

	return base + "_modified.fit"
}
