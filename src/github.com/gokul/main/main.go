package main

import (
	"github.com/gokul/server"
	log "github.com/logrus"
	"os"
	"flag"
	"github.com/gokul/cmd"
)


func init() {
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Debug("Starting the Server...")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 || args[0] == "help" {
		log.Fatalln("No arguments are supplied")
		// call the help

	}else if len(args) > 1 {
		for _,cmd := range cmd.Commands{
			if cmd.Name() == args[0]{

			}
		}


	}

	log.Debugln("Arguments are :: ",args)
	server := gokul.NewServer()
	server.ScanAppsForControllers()
	gokul.Run(server)
}
