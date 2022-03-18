package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"webscanner/utility"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for both one url or multiple urls from file",
	Long:  "Scan for individual url or import urls from file",
}

var scanUrlCmd = &cobra.Command{
	Use:   "url",
	Short: "Scan for target url",
	Long:  "Scan for individual url",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scanning for single url..")
		targetURL := args[0]
		utility.LoadUrl(false, targetURL, 1, saveResult, saveFileName)
	},
}

var scanFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Scan for multiple urls",
	Long:  "Scan for multiple urls by importing urls from file",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]
		utility.LoadUrl(true, fileName, threadNum, saveResult, saveFileName)
	},
}
