package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gokul",
	Long:  `All software has versions. This is gokul's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\t" + VERSION + "\n\t" + RootCmd.Short)
	},
}
