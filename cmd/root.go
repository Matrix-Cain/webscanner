package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "fofa",
	Short: `
________       ________        
___  __/______ ___  __/______ _
__  /_  _  __ \__  /_  _  __ '/
_  __/  / /_/ /_  __/  / /_/ / 
/_/     \____/ /_/     \__,_/\\
`,
}
var threadNum int
var saveResult bool
var saveFileName string

func init() {
	rootCmd.AddCommand(versionCmd)

	scanUrlCmd.Flags().StringVarP(&saveFileName, "name", "n", "", "Naming results file")
	scanUrlCmd.Flags().BoolVarP(&saveResult, "save", "s", false, "Saving results to a file")
	scanCmd.AddCommand(scanUrlCmd)

	scanFileCmd.Flags().StringVarP(&saveFileName, "name", "n", "", "Naming results file")
	scanFileCmd.Flags().IntVarP(&threadNum, "thread", "t", 10, "Specify the worker pool size")
	scanFileCmd.Flags().BoolVarP(&saveResult, "save", "s", false, "Saving results to a file")
	scanCmd.AddCommand(scanFileCmd)

	rootCmd.AddCommand(scanCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
