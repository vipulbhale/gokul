package main

import (
	"github.com/vipulbhale/gokul/cmd"
)

var (
	// Version of the application
	VERSION string = "0.0.1"
)

func init() {
}

func main() {
	cmd.Execute(VERSION)
}
