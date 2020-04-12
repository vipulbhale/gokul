package internal

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vipulbhale/gokul/pkg/server/apptempl"
)

var cmdNew = &cobra.Command{
	Use:   "new",
	Short: "Make a new Application",
	Long:  `Command Used to create application`,
	Run:   createNewApplication,
}

func init() {
	cmdNew.Flags().StringVarP(&AppDirName, "dir", "d", "", "Directory under which app needs to be created.")
	CmdApp.AddCommand(cmdNew)
}

func createNewApplication(cmd *cobra.Command, args []string) {
	Log.Debugln("Inside the new application command")
	Log.Debugln("Scanning all existing apps for controllers")

	//Add the directory to the goPath
	//First check in GOPATH
	if val := isPresentAndInGoPath(AppDirName); val {
		Log.Debugln("Is AppDirName present in GOPATH", val)
	} else {
		Log.Debugln("Adding the AppsDirectory to GOPATH")
		addToGoPath(AppDirName)
		Log.Debugln("After making changes current GOPATH is ", os.Getenv("GOPATH"))
	}

		// // First check AppDirName is provided
	if len(AppDirName) == 0 {
		AppDirName = os.Getenv("GOPATH")
	}
	// all application/s are scanned now copy required server files to the apps directory
	apptempl.CreateTemplates(AppDirName, AppName, CfgFileLocation)
}

/**
Check whether the AppDirectory is present in the the GOPATH
*/
func isPresentAndInGoPath(appDirName string) bool {
	//check if present in GOPATH
	gopath := os.Getenv("GOPATH")
	Log.Debugln("Current GoPath is ", gopath)
	if strings.Contains(gopath, appDirName) {
		return true
	}
	return false
}

/*
*
*	Add the directory to gopath
 */
func addToGoPath(appDirName string) {
	gopath := os.Getenv("GOPATH")
	Log.Debugln("Current GoPath is ", gopath)
	newGoPath := gopath + ":" + appDirName
	os.Setenv("GOPATH", newGoPath)
	Log.Debugln("After making changes current GOPATH is ", os.Getenv("GOPATH"))
}
