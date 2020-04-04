package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Logger *logrus.Logger

func init() {
	initConfig()
}

func SetLogger(logDestination io.Writer, loggingLevel logrus.Level) {
	Logger = logrus.New()
	Logger.SetLevel(loggingLevel)
	Logger.SetOutput(logDestination)
}

func GetLogger() *logrus.Logger {
	return Logger
}

func initConfig() {
	v, err := setupViper()
	if err != nil {
		fmt.Printf("Can't read config, %v. Creating a template configfile in current directory.\n", err)
		cmdConfigContent := []byte("logging:\n  level: info\n  destination: stdout\n")
		err := ioutil.WriteFile("./gokul.yaml", cmdConfigContent, 0644)
		if err != nil {
			fmt.Println("Not able to create the config file in current directory. Exiting...", err)
			os.Exit(1)
		}
		v.ReadInConfig()
	}
	loggingLevel := v.Get("logging.level")
	loggingDestination := v.Get("logging.destination")
	fmt.Printf("Logging Level is %v \n", loggingLevel)
	fmt.Printf("Logging destintation is %v \n", loggingDestination)
	setupLogging(loggingDestination.(string), loggingLevel.(string))
}

func setupViper() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("gokul") // name of config file (without extension)
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/gokul/")  // path to look for the config file in
	v.AddConfigPath("$HOME/.gokul") // call multiple times to add many search paths
	v.AddConfigPath(".")            // optionally look for config in the working directory
	v.AutomaticEnv()                // read in environment variables that match
	err := v.ReadInConfig()
	return v, err
}

func setupLogging(loggingDestFromConfig string, loggingLevel string) {
	var logFileName string = ""
	var loggingDestination io.Writer
	var logLevel logrus.Level
	var err error

	// if len(loggingDestFromConfig) != 0 {
	if loggingDestFromConfig == "stdout" {
		loggingDestination = os.Stdout
	} else if strings.HasPrefix(loggingDestFromConfig, "file://") {
		logFileName = strings.Split(loggingDestFromConfig, "file://")[0]
		file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			fmt.Printf("Error while creating the error file %v \n", err)
			fmt.Printf("Using the default location /var/log/gokul/cmd.log. \n")
			CreateDirectory("/var/log/gokul")
			file, err = os.OpenFile("/var/log/gokul/cmd.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		}
		loggingDestination = file
	} else {
		CreateDirectory("/var/log/gokul")
		file, err := os.OpenFile("/var/log/gokul/cmd.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			fmt.Printf("Error while creating the error file :: %v .\n", err)
			fmt.Printf("Using the stdout as the log destination.\n")
			loggingDestination = os.Stdout
		} else {
			loggingDestination = file
		}
	}
	// }

	if len(loggingLevel) != 0 {
		// Only log the debug severity or above.
		logLevel, err = logrus.ParseLevel(loggingLevel)
		if err != nil {
			fmt.Println("Loglevel not set correctly in config file.")
			logLevel = logrus.InfoLevel
		}
	} else {
		logLevel = logrus.InfoLevel
	}
	fmt.Println("Set the logger with correct destination and logLevel")
	SetLogger(loggingDestination, logLevel)
}
