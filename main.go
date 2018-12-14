package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vipulbhale/gokul/cmd"
)

var (
	// Version of the application
	VERSION = "0.0.1"
)

func init() {
	cobra.OnInitialize(initConfig)

}

func main() {
	cmd.Execute(VERSION)
}

func initConfig() {
	viper.SetConfigName("gokul") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/gokul/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.gokul") // call multiple times to add many search paths
	viper.AddConfigPath(".")            // optionally look for config in the working directory
	viper.AutomaticEnv()                // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Can't read config, %v. Creating a template configfile in current directory.\n", err)
		cmdConfigContent := []byte("logging:\n  level: debug\n  destination: stdout\n")
		err := ioutil.WriteFile("./gokul.yaml", cmdConfigContent, 0644)
		if err != nil {
			fmt.Println("Not able to create the config file in current directory. Exiting...", err)
			os.Exit(1)
		}
	}
	loggingLevel := viper.Get("logging.level")
	loggingDestination := viper.Get("logging.destination")
	fmt.Printf("Logging level is %v\n", viper.AllKeys())
	if loggingDestination != nil {
		switch loggingDestination.(string) {
		case "stdout":
			log.SetOutput(os.Stdout)
		default:
			log.SetOutput(os.Stdout)
		}
	}
	if loggingLevel != nil {
		// Only log the debug severity or above.
		logLevel, err := log.ParseLevel(loggingLevel.(string))
		if err != nil {
			fmt.Println("Loglevel not set correctly in config file.")
		}
		log.SetLevel(logLevel)
	}
}
