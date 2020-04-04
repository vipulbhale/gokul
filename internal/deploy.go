package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vipulbhale/gokul/server/config"
	goreflect "github.com/vipulbhale/gokul/server/reflect"
)

func init() {
	// Get the app location with implicit app name parameter.
	// App name is also equal to context root.
	// Get the config location with server.yml location
	CmdApp.AddCommand(cmdDeploy)
}

var cmdDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy app to server",
	Long:  `Deploy given app to the server.`,
	Run:   deployApp,
}

func executeOSCommand(commandString string, parameters ...string) {
	goPath, err := exec.LookPath(commandString)
	if err != nil {
		Log.Fatalln("Error while getting the path of the go binary", err)
	}
	Log.Debug("The gopath is ", goPath)

	command := exec.Command(goPath, parameters...)
	command.Env = []string{"GOPATH=" + filepath.Join(AppDirName), "PATH=" + os.Getenv("PATH")}
	stderr, err := command.StderrPipe()
	if err != nil {
		Log.Fatal(err)
	}

	if err := command.Start(); err != nil {
		Log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := command.Wait(); err != nil {
		Log.Fatal(err)
	}
}

func deployApp(cmd *cobra.Command, args []string) {
	Log.Debugln("Deploying the apps")
	Log.Debugln("Scanning all existing apps for controllers")
	if len(CfgFileLocation) > 0 {
		config.LoadConfigFile(CfgFileLocation)
	} else {
		CfgFileLocation = filepath.Join(AppDirName, "src", "github.com", AppName, "config", "server.yml")
		config.LoadConfigFile(CfgFileLocation)
	}

	executeOSCommand("go", "get", "-u", "-d", "github.com/vipulbhale/gokul/server")
	// start scanning all controllers for the given app or apps directory
	if len(AppName) != 0 {
		goreflect.ScanAppsDirectory(config.Cfg, AppName)
	} else {
		goreflect.ScanAppsDirectory(config.Cfg, "")
	}
	// all application/s are scanned now copy required server files to the apps directory
}
