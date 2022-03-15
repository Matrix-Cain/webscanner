package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of fofa",
	Long:  "All software has versions. This is fofa's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fofa webapp fingerprint scanner v1.0 -- HEAD")
	},
}
