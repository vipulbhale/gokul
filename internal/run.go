package internal

import (
	"bufio"
	"fmt"
	"os"
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

func commandExecutor(cmdArgs ...string) {
	Log.Debugln("Entering commandExecutor function.")
	goBinaryPath, err := exec.LookPath("go")
	if err != nil {
		Log.Fatalln("Error while getting the path of the go binary", err)
	}
	Log.Debugln("The path of go binary is :: ", goBinaryPath)
	appRoot := filepath.Join(AppDirName, "src", "github.com", AppName)
	command := exec.Command(goBinaryPath, cmdArgs[:]...)
	command.Dir = appRoot
	command.Env = []string{"GOPATH=" + os.Getenv("GOPATH"), "PATH=" + os.Getenv("PATH"), "HOME=" + os.Getenv("HOME")}
	stdOutReader, errors := command.StdoutPipe()
	stdErrReader, errors := command.StderrPipe()
	done := make(chan struct{})
	Log.Debugln("Running the command :: ", command.String)
	scannerStdOut := bufio.NewScanner(stdOutReader)
	scannerStdErr := bufio.NewScanner(stdErrReader)
	go func() {
		for scannerStdOut.Scan() {
			fmt.Printf("%s\n", scannerStdOut.Text())
		}
		for scannerStdErr.Scan() {
			fmt.Printf("%s\n", scannerStdErr.Text())
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
		Log.Info(err)
	}
}

func runApp(cmd *cobra.Command, args []string) {
	Log.Debugln("Location of variable.env is :: ", filepath.Join(AppDirName, "src", "github.com", AppName, "variables.env"))
	commandExecutor("mod", "init", "github.com/"+AppName)
	commandExecutor("run", filepath.Join(AppDirName, "src", "github.com", AppName, "cmd", AppName, "main.go"))
}
