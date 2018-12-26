package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"

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
	Log.Debugln("Location of variable.env is ::  {}", filepath.Join(AppDirName, "src", "github.com", AppName, "variables.env"))
	goPath, err := exec.LookPath("go")
	if err != nil {
		Log.Fatalln("Error while getting the path of the go binary", err)
	}
	Log.Debug("The gopath is :: ", goPath)
	command := exec.Command(goPath, "run", filepath.Join(AppDirName, "src", "github.com", AppName, "main.go"))
	stdOutReader, errors := command.StdoutPipe()
	done := make(chan struct{})
	Log.Debugln("Running the main.go of app :: ", AppName)
	scanner := bufio.NewScanner(stdOutReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}

		done <- struct{}{}
	}()

	if errors != nil {
		Log.Fatal(errors)
	}
	if err := command.Start(); err != nil {
		Log.Fatal(err)
	}

	<-done
	if err := command.Wait(); err != nil {
		Log.Fatal(err)
	}

}
