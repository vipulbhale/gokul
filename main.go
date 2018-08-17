package main

import (
	//"github.com/gokul/server"
	"os"

	log "github.com/sirupsen/logrus"
	//"flag"
	"github.com/gokul/cmd"
)

var (
	// Version of the application
	VERSION = "0.0.1"
)

func init() {

	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)
	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)

}

func main() {
	log.Debugln("After making changes current GOPATH is ", os.Getenv("GOPATH"))
	cmd.Execute(VERSION)
	log.Debugln("After making changes current GOPATH is ", os.Getenv("GOPATH"))
}
