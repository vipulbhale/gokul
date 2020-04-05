package internal

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/vipulbhale/gokul/pkg/server/util"

	"github.com/spf13/cobra"
)

var (
	userLicense string
	Log         *logrus.Logger
	VERSION     string
)

var RootCmd = &cobra.Command{
	Use:   "gokul",
	Short: "gokul is used to generate stubs for web application , deploy the application , run the application.",
	Long:  `A web application stub generator, deployer and runner`,
}

func init() {
	Log = util.GetLogger()
}

// Execute adds all child commands to the root command
func Execute(version string) {
	VERSION = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
