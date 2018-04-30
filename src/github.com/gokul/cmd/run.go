package cmd

import (
	"github.com/gokul/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFileLocation string

func init() {
	cmdRun.Flags().StringVarP(&cfgFileLocation, "config", "c", "", "Config file location.")
	viper.BindPFlag("config", cmdRun.PersistentFlags().Lookup("config"))
	RootCmd.AddCommand(cmdRun)
}

var cmdRun = &cobra.Command{
	Use:   "run [app to run]",
	Short: "Run on server",
	Long:  `Run app to the server`,
	Run:   runApp,
}

func runApp(cmd *cobra.Command, args []string) {
	log.Debug("Config File Location is ", cfgFileLocation)
	log.Debug("Starting the Server...")
	log.Debug("Creating the new server.")
	server := gokul.NewServer(cfgFileLocation)
	log.Debug("Scanning all apps for controllers")
	server.ScanAppsForControllers()
	log.Debug("Run the server")
	gokul.Run(server)
}
