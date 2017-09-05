package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "gokul",
	Short: "gokul is used to generate stubs for web application , deploy the application , run the application.",
	Long: `A web application stub generator , deployer and runner`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
