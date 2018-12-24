package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	CmdApp.AddCommand(cmdRun)
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run opp on server",
	Long:  `Run app to the server`,
	Run:   runApp,
}

func runApp(cmd *cobra.Command, args []string) {
	//log.Debug("Config File Location is ", CfgFileLocation)
	// if len(CfgFileLocation) > 0 {
	// 	config.LoadConfigFile(CfgFileLocation)
	// }

	// server := gokul.NewServer(CfgFileLocation)
	// log.Debug("Scanning an app for controllers")
	// server.ScanAppsForControllers(AppName)
	// log.Debug("Run the server")
	//gokul.Run(server)
	log.Debugln("Location of variable.env is ::  {}", filepath.Join(AppDirName, "src", "github.com", AppName, "variables.env"))
	// command := exec.Command("bash", "-c", "source", filepath.Join(AppDirName, "src", "github.com", AppName, "variables.env"))
	goPath, err := exec.LookPath("go")
	if err != nil {
		log.Fatalln("Error while getting the path of the go binary", err)
	}
	log.Debug("The gopath is :: ", goPath)
	command := exec.Command(goPath, "run", filepath.Join(AppDirName, "src", "github.com", AppName, "main.go"))
	stdOutReader, errors := command.StdoutPipe()
	done := make(chan struct{})
	log.Debugln("Running the main.go of app :: ", AppName)
	scanner := bufio.NewScanner(stdOutReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}

		done <- struct{}{}
	}()

	if errors != nil {
		log.Fatal(errors)
	}
	if err := command.Start(); err != nil {
		log.Fatal(err)
	}

	<-done
	if err := command.Wait(); err != nil {
		log.Fatal(err)
	}

}
