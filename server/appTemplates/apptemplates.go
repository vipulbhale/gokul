package appTemplates

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	log "github.com/sirupsen/logrus"
)

var tpl bytes.Buffer

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

const ROUTES_CFG_TEMPLATE = `
GET             /demo           DemoController.Demo
`
const MAIN_PACKAGE = `package main
import (
	"fmt" 
	"github.com/vipulbhale/gokul/server"
	"github.com/vipulbhale/gokul/server/config"
	"github.com/{{.AppNameForTemplate}}/controller"
	log "github.com/sirupsen/logrus"
)

const cfgFileLocation = "{{.ParentAppDirectory}}/{{.AppNameForTemplate}}/config/server.yml"

func main(){
	fmt.Println("Starting the Server")
	log.Debug("Config File Location is ", cfgFileLocation)
	if len(cfgFileLocation) > 0 {
		config.LoadConfigFile(cfgFileLocation)
	}

	mapOfControllerNameToControllerObj := controller.RegisterControllers()
	appServer := server.NewServer(cfgFileLocation,mapOfControllerNameToControllerObj)
	log.Debug("Run the server")
	server.Run(appServer)
}
`
const CONTROLLER_TEMPLATE = `package controller

import (
	log "github.com/sirupsen/logrus"
	"github.com/tempapp/model"
	controller2 "github.com/vipulbhale/gokul/server/controller"
)

type DemoController struct {
	*controller2.BaseController
}

//Demo method for GET /demo endpoint
func (d *DemoController) Demo() (error, *controller2.ModelAndView) {
	log.Debugln("Inside the Demo method of DemoController")
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
	log.Debugln("Inside the Demo method of DemoController")
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
	log.Debugln("Inside the Demo method of DemoController")
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
const MODEL_TEMPLATE = `package model

type Person struct {
	Name string
	Age int
}
`
const SERVICE_TEMPLATE = `package service
import "github.com/tempapp/model"

func GetPerson() Person{
	person := new(Person)
	person.Name = "Gokul"
	person.Age = 39
	return person
}
`

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
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "config"), "server.yml", appName, cfgFileLocation, CONFIG_FILE_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "config"), "routes.cfg", appName, cfgFileLocation, ROUTES_CFG_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName), "variables.env", appName, cfgFileLocation, "GOPATH="+os.Getenv("GOPATH"))

	createBinAndPackageDirectory(filepath.Join(dirname, "bin"))
	createBinAndPackageDirectory(filepath.Join(dirname, "pkg"))
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
