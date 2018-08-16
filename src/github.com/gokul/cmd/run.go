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
	// cmdRun.Flags().StringVarP(&AppDirName, "dir", "d", "", "Directory under which app needs to be run.")
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
	fmt.Println("hi there " + filepath.Join(AppDirName, "src", "github.com", AppName, "variables.env"))
	// command := exec.Command("bash", "-c", "source", filepath.Join(AppDirName, "src", "github.com", AppName, "variables.env"))
	goPath, err := exec.LookPath("go")
	if err != nil {
		log.Fatalln("Error while getting the path of the go binary", err)
	}
	log.Debug("The gopath is ", goPath)
	command := exec.Command(goPath, "run", filepath.Join(AppDirName, "src", "github.com", AppName, "main.go"))
	stdOutReader, errors := command.StdoutPipe()
	// command.Env = append(os.Environ(),
	// 	"GOPATH=$GOPATH:"+AppDirName,
	// )
	done := make(chan struct{})

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
