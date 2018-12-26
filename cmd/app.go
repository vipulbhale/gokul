package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//AppName is the name of the app
	AppName string
	//CfgFileLocation is the location of the file
	CfgFileLocation string
	//AppDirName is the name of the directory where application needs to be run/copied etc.
	AppDirName string
)

func init() {
	// Get the app location with implicit app name parameter.
	// App name is also equal to context root.
	// Get the config location with server.yml location

	CmdApp.PersistentFlags().StringVarP(&AppDirName, "dir", "d", "", "Directory under which app needs to be run.")
	CmdApp.PersistentFlags().StringVarP(&CfgFileLocation, "config", "c", "", "Config file location(default is $HOME/.gokul.yaml)")
	CmdApp.PersistentFlags().StringVarP(&AppName, "name", "n", "", "App name. Will search for the same name in apps directory.")
	viper.BindPFlag("config", CmdApp.PersistentFlags().Lookup("config"))
	RootCmd.AddCommand(CmdApp)
}

//CmdApp is the command corresponding to app.
var CmdApp = &cobra.Command{
	Use:   "app",
	Short: "Application command",
	Long:  `Command Used to create, deploy, run application`,
	Run:   application,
}

func application(cmd *cobra.Command, args []string) {
	Log.Debugln("Inside the main application command")
}
