package cmd

import (
	"os"
	"strings"

	"github.com/gokul/server/appTemplates"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmdNew.Flags().StringVarP(&AppDirName, "dir", "d", "", "Directory under which app needs to be created.")
	CmdApp.AddCommand(cmdNew)
}

var cmdNew = &cobra.Command{
	Use:   "new",
	Short: "Make a new Application",
	Long:  `Command Used to create application`,
	Run:   createNewApplication,
}

func createNewApplication(cmd *cobra.Command, args []string) {
	log.Debugln("Inside the new application command")
	log.Debugln("Scanning all existing apps for controllers")
	// First check AppDirName is provided
	if len(AppDirName) == 0 {
		panic("Application Directory is not provided")
	}
	//Add the directory to the goPath
	//First check in GOPATH
	if val := isPresentInGoPath(AppDirName); val {
		log.Debugln("Is AppDirName present in GOPATH", val)
	} else {
		log.Debugln("Adding the AppsDirectory to GOPATH")
		addToGoPath(AppDirName)
		log.Debugln("After making changes current GOPATH is ", os.Getenv("GOPATH"))
	}
	// all application/s are scanned now copy required server files to the apps directory
	appTemplates.CreateTemplates(AppDirName, AppName, CfgFileLocation)
}

/**
Check whether the AppDirectory is present in the the GOPATH
*/
func isPresentInGoPath(appdirname string) bool {
	//check if present in GOPATH
	gopath := os.Getenv("GOPATH")
	log.Debugln("Current GoPath is ", gopath)
	if strings.Contains(gopath, appdirname) {
		return true
	}
	return false
}

func addToGoPath(appDirName string) {
	gopath := os.Getenv("GOPATH")
	log.Debugln("Current GoPath is ", gopath)
	newGoPath := gopath + ":" + appDirName
	os.Setenv("GOPATH", newGoPath)
	log.Debugln("After making changes current GOPATH is ", os.Getenv("GOPATH"))
}
