package main

import (
	//"github.com/gokul/server"
	log "github.com/logrus"
	"os"
	//"flag"
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

	cmd.RootCmd.Execute()
	//
	//server := gokul.NewServer()
	//server.ScanAppsForControllers()
	//gokul.Run(server)
}
