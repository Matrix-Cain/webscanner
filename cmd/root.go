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

func init() {
	rootCmd.AddCommand(versionCmd)

	scanCmd.AddCommand(scanUrlCmd)
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
