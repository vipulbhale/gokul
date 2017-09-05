package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"strings"
)

func init() {
	RootCmd.AddCommand(cmdDeploy)
}

var cmdDeploy = &cobra.Command{
	Use:   "deployapp [app to deploy]",
	Short: "Deployapp to server",
	Long:  `Deploy app to the server`,
	Run: deployApp,
}


func deployApp(cmd *cobra.Command, args []string) {
	fmt.Println(strings.Join(args, " "))
}