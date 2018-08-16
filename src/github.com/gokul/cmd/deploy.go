package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gokul/server/config"
	goreflect "github.com/gokul/server/reflect"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

func deployApp(cmd *cobra.Command, args []string) {
	log.Debugln("Deploying the apps")
	log.Debugln("Scanning all existing apps for controllers")
	if len(CfgFileLocation) > 0 {
		config.LoadConfigFile(CfgFileLocation)
	} else {
		CfgFileLocation = filepath.Join(AppDirName, "src", "github.com", AppName, "config", "server.yml")
		config.LoadConfigFile(CfgFileLocation)
	}

	goPath, err := exec.LookPath("go")
	if err != nil {
		log.Fatalln("Error while getting the path of the go binary", err)
	}
	log.Debug("The gopath is ", goPath)

	command := exec.Command("go", "get", "-u", "github.com/vipulbhale/gokul", "github.com/sirupsen/logrus")
	command.Env = []string{"GOPATH=" + filepath.Join(AppDirName), "PATH=" + os.Getenv("PATH")}
	stderr, err := command.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := command.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := command.Wait(); err != nil {
		log.Fatal(err)
	}
	// start scanning all controllers for the given app or apps directory
	if len(AppName) != 0 {
		goreflect.ScanAppsDirectory(config.Cfg, AppName)
	} else {
		goreflect.ScanAppsDirectory(config.Cfg, "")
	}
	// all application/s are scanned now copy required server files to the apps directory
}
