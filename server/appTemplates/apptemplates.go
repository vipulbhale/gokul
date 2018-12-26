package appTemplates

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/vipulbhale/gokul/server/util"
)

var tpl bytes.Buffer
var log *logrus.Logger

type ApplicationForTemplate struct {
	AppNameForTemplate string
	ParentAppDirectory string
	CfgFileLocation    string
}

const CONFIG_FILE_TEMPLATE = `server:
  port: 9000
  address: "0.0.0.0"
logging:
  level: debug
http:
  read.timeout: 0
  write.timeout: 0
  maxrequestsize: 999999
apps:
  directory: {{.ParentAppDirectory}}
`

const ROUTES_CFG_TEMPLATE = `GET             /demo           DemoController.Demo
GET             /demoxml           DemoController.DemoXML
GET             /demojson           DemoController.DemoJson
`
const MAIN_PACKAGE = `package main
import (
	"fmt" 
	"github.com/vipulbhale/gokul/server"
	"github.com/vipulbhale/gokul/server/config"
	"github.com/{{.AppNameForTemplate}}/controller"
	"github.com/{{.AppNameForTemplate}}/util"
	"github.com/sirupsen/logrus"
)

var Log         *logrus.Logger
const cfgFileLocation = "{{.ParentAppDirectory}}/{{.AppNameForTemplate}}/config/server.yml"

func init(){
	Log = util.GetLogger()
}

func main(){
	fmt.Println("Starting the Server")
	Log.Debug("Config File Location is ", cfgFileLocation)
	if len(cfgFileLocation) > 0 {
		config.LoadConfigFile(cfgFileLocation)
	}

	mapOfControllerNameToControllerObj := controller.RegisterControllers()
	appServer := server.NewServer(cfgFileLocation,mapOfControllerNameToControllerObj)
	Log.Debug("Run the server")
	server.Run(appServer)
}
`
const CONTROLLER_TEMPLATE = `package controller

import (
	controller2 "github.com/vipulbhale/gokul/server/controller"
	"github.com/{{.AppNameForTemplate}}/util"
	"github.com/{{.AppNameForTemplate}}/service"
	"github.com/sirupsen/logrus"
)
var Log         *logrus.Logger

func init(){
	Log = util.GetLogger()
}

type DemoController struct {
	*controller2.BaseController
}

//Demo method for GET /demo endpoint
func (d *DemoController) Demo() (error, *controller2.ModelAndView) {
	Log.Debugln("Inside the Demo method of DemoController")
	//Create the instance of ModelAndView
	modelAndView := new(controller2.ModelAndView)
	// Create the instance of the actual model struct
	person := service.GetPerson() 
	// Set the Model and View in the ModelAndView Struct required by Gokul
	modelAndView.SetModel(person)
	modelAndView.SetView("view")
	modelAndView.SetResponseType("text/html")
	//d.Render()
	return nil, modelAndView
}

//Demo method for GET /demo endpoint
func (d *DemoController) DemoXML() (error, *controller2.ModelAndView) {
	Log.Debugln("Inside the Demo method of DemoController")
	//Create the instance of ModelAndView
	modelAndView := new(controller2.ModelAndView)
	// Create the instance of the actual model struct
	person := service.GetPerson() 
	// Set the Model and View in the ModelAndView Struct required by Gokul
	modelAndView.SetModel(person)
	modelAndView.SetResponseType("application/xml")
	//d.Render()
	return nil, modelAndView
}

//Demo method for GET /demo endpoint
func (d *DemoController) DemoJson() (error, *controller2.ModelAndView) {
	Log.Debugln("Inside the Demo method of DemoController")
	//Create the instance of ModelAndView
	modelAndView := new(controller2.ModelAndView)
	// Create the instance of the actual model struct
	person := service.GetPerson()
	// Set the Model and View in the ModelAndView Struct required by Gokul
	modelAndView.SetModel(person)
	modelAndView.SetResponseType("application/json")
	//d.Render()
	return nil, modelAndView
}
`
const SERVICE_TEMPLATE = `package service
import ( 
	"github.com/{{.AppNameForTemplate}}/model"
	"github.com/{{.AppNameForTemplate}}/util"
	"github.com/sirupsen/logrus"
)

var 	Log         *logrus.Logger

func init(){
	Log = util.GetLogger()	
}

func GetPerson() *model.Person { 
	person := new(model.Person)
	person.Name = "Gokul"
	person.Age = 39
	return person
}
`

const MODEL_TEMPLATE = `package model

type Person struct {
	Name string
	Age int
}
`
const VIEW_TEMPLATE = `<html>
	<body>
		hi there
		<h1>{{.Name}}</h1>
		<h1>{{.Age}}</h1>
	</body>
	</html>
`

const UTIL_LOGGER_TEMPLATE = `package util

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
		cmdConfigContent := []byte("logging:\n  level: info\n  destination: \"\"\n")
		err := ioutil.WriteFile("./application.yaml", cmdConfigContent, 0644)
		if err != nil {
			fmt.Println("Not able to create the config file in current directory. Exiting...", err)
			os.Exit(1)
		}
		v.ReadInConfig()
	}
	loggingLevel := v.Get("logging.level")
	loggingDestination := v.Get("logging.destination")
	fmt.Printf("Logging Level is :: %v.\n", loggingLevel)
	fmt.Printf("Logging destintation is :: %v.\n", loggingDestination)
	setupLogging(loggingDestination.(string), loggingLevel.(string))
}

func setupViper() (*viper.Viper, error){
	v := viper.New()
	v.SetConfigName("application") // name of config file (without extension)
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/{{.AppNameForTemplate}}/")  // path to look for the config file in
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

	if len(loggingDestFromConfig) != 0 {
		if loggingDestFromConfig == "stdout" {
			loggingDestination = os.Stdout
		} else if strings.HasPrefix(loggingDestFromConfig, "file://") {
			logFileName = strings.Split(loggingDestFromConfig, "file://")[0]
			file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
			if err != nil {
				fmt.Printf("Error while creating the error file :: %v.\n", err)
				fmt.Printf("Using the default location /var/log/{{.AppNameForTemplate}}/app.log.\n")
				file, err = os.OpenFile("/var/log/{{.AppNameForTemplate}}/app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
			}
			loggingDestination = file
		} else {
			file, err := os.OpenFile("/var/log/{{.AppNameForTemplate}}/app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
			if err != nil {
				fmt.Printf("Error while creating the error file :: %v.\n", err)
				fmt.Printf("Using the stdout as the log destination.\n")
				loggingDestination = os.Stdout
			}
			loggingDestination = file
		}
	}

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
`

func init() {
	log = util.GetLogger()
}

//CreateTemplates Method creates the appName and templates in the directory mentioned.
func CreateTemplates(dirname, appName, cfgFileLocation string) {
	log.Debugln("Entering the CreateTemplates method.")

	// writeToFile(filepath.Join(dirname, "src", "github.com", appName, "routes"), "route.go", appName, cfgFileLocation, ROUTES_TEMPLATE)
	// writeToFile(filepath.Join(dirname, "src", "github.com", appName, "server"), "server.go", appName, cfgFileLocation, SERVER_TEMPLATE)
	// writeToFile(filepath.Join(dirname, "src", "github.com", appName, "config"), "config.go", appName, cfgFileLocation, SERVER_CONFIG_TEMPLATE)
	// writeToFileReflect(filepath.Join(dirname, "src", "github.com", appName, "reflect", "reflect.go"), GOKUL_REFLECT_TEMPLATE)

	writeToFile(filepath.Join(dirname, "src", "github.com", appName), "main.go", appName, cfgFileLocation, MAIN_PACKAGE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "controller"), "controller.go", appName, cfgFileLocation, CONTROLLER_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "service"), "service.go", appName, cfgFileLocation, SERVICE_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "model"), "model.go", appName, cfgFileLocation, MODEL_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "view"), "view.html", appName, cfgFileLocation, VIEW_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "util"), "logger.go", appName, cfgFileLocation, UTIL_LOGGER_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "config"), "server.yml", appName, cfgFileLocation, CONFIG_FILE_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "config"), "routes.cfg", appName, cfgFileLocation, ROUTES_CFG_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName), "variables.env", appName, cfgFileLocation, "GOPATH="+os.Getenv("GOPATH"))

	createBinAndPackageDirectory(filepath.Join(dirname, "bin"))
	createBinAndPackageDirectory(filepath.Join(dirname, "pkg"))
	createBinAndPackageDirectory(filepath.Join("var", "log", appName))
}

// func writeToFileReflect(fileName, content string) {
// 	log.Debugln("For the reflect directory is ", filepath.Dir(fileName))
// 	if _, err := os.Stat(filepath.Dir(fileName)); os.IsNotExist(err) {
// 		err = os.MkdirAll(filepath.Dir(fileName), 0755)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	d1 := []byte(content)
// 	err := ioutil.WriteFile(fileName, d1, 0644)
// 	if err != nil {
// 		log.Fatalln("Error while copying the reflect.go to apps directory")
// 	}
// }

func writeToFile(dirName, fileName, appName, cfgFileLocation, content string) {
	log.Debugln("Entering the writeToFile method with parameters", dirName, fileName, appName, cfgFileLocation)
	appTemplate := ApplicationForTemplate{AppNameForTemplate: appName, ParentAppDirectory: filepath.Dir(dirName), CfgFileLocation: cfgFileLocation}

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			panic(err)
		}
	}
	outputFile, outputError := os.Create(filepath.Join(dirName, fileName))

	if outputError != nil {
		log.Fatalln("An error occurred with file creation")
		panic(outputError)
	}

	defer outputFile.Close()
	outputString := content

	t := template.New("AppTemplates")
	t, _ = t.Parse(outputString)
	if err := t.Execute(outputFile, appTemplate); err != nil {
		log.Fatalln("There was an error while template execution:", err.Error())
		panic(err)

	}
	log.Debugln("Done writing a template.")
}

func createBinAndPackageDirectory(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			panic(err)
		}
	}
}
