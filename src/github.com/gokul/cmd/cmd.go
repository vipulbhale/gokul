package cmd

import "strings"

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

var Commands = []*Command{
	cmdDeploy,

}

func parseArgs() {

}

