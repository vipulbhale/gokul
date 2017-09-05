package cmd

import (
	"github.com/gokul/server"

	"github.com/spf13/cobra"
	"fmt"
	"strings"
)

func init() {
	RootCmd.AddCommand(cmdRun)
}

var cmdRun = &cobra.Command{
	Use:   "run [app to run]",
	Short: "Run on server",
	Long:  `Run app to the server`,
	Run: runApp,
}


func runApp(cmd *cobra.Command, args []string) {
	fmt.Println(strings.Join(args, " "))

	server := gokul.NewServer()
	server.ScanAppsForControllers()
	gokul.Run(server)
}