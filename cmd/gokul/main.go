package main

import "github.com/vipulbhale/gokul/internal"

// Version of the application
var VERSION string = "0.0.1"

func init() {
}

func main() {
	internal.Execute(VERSION)
}
