package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

var srcRoot string

// Command structure cribbed from the genius organization of the "go" command.
type Command struct {
	Execute                    func(args []string)
	UsageLine, Short, Long string
}

func (cmd *Command) Name() string {
	name := cmd.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

var commands = []*Command{
	cmdDeploy,

}

func main() {
	//var err error
	fmt.Println("Hello World!")
	demo, err := build.Import("github.com/gokul", "", build.FindOnly)
	fmt.Println(err, demo)
	workingDir, _ := os.Getwd()
	goPathList := filepath.SplitList(build.Default.GOPATH)
	fmt.Printf("Current Working directory is %s\n", workingDir)
	for _, path := range goPathList {
		fmt.Printf("Gopath is %s\n", path)

		if strings.HasPrefix(strings.ToLower(workingDir), strings.ToLower(path)) {
			srcRoot = path
			fmt.Printf("%s,,%s", srcRoot, path)
			break
		}

	}

	fmt.Printf("filePath base is %s\n", filepath.Base(workingDir))
	fmt.Printf("tseting %s \n", filepath.ToSlash(filepath.Dir("import/src/depot")))
}
