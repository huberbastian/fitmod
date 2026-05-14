package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fitmod",
	Short: "A tool for injecting distance data into FIT files.",
	Long: `A tool for injecting distance data into FIT files.

fitmod allows you to inject distance data into FIT files from activities that are missing that information.
This can be useful for indoor cycling or running activities where the recording device does not receive distance data from
the trainer or treadmill.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(processCmd)
}
