package appTemplates

import (
	"bytes"
	"io/ioutil"
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

const MAIN_PACKAGE = `package main
import (
	"fmt" 
	// "github.com/{{.AppNameForTemplate}}/server"
	// "github.com/{{.AppNameForTemplate}}/config"
	"github.com/vipulbhale/gokul/server"
	"github.com/vipulbhale/gokul/config"
	log "github.com/sirupsen/logrus"
)

const cfgFileLocation = "{{.ParentAppDirectory}}/{{.AppNameForTemplate}}/config/server.yml"

func main(){
	fmt.Println("Starting the Server")
	log.Debug("Config File Location is ", cfgFileLocation)
	if len(cfgFileLocation) > 0 {
		config.LoadConfigFile(cfgFileLocation)
	}

	appServer := server.NewServer(cfgFileLocation)
	log.Debug("Scanning an app for controllers")
	appServer.ScanAppsForControllers("")
	log.Debug("Run the server")
	gokul.Run(appServer)
}
`

const CONTROLLER_TEMPLATE = `package basecontroller

import (
	"fmt"
)

type Controller interface {
	Render()
}

type BaseController struct {
	Controller

}

func (baseController *BaseController) Render() {
	fmt.Println("Inside render method")
}
`

const ROUTES_TEMPLATE = `
package routes

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	// "github.com/{{.AppNameForTemplate}}/config"
	"github.com/vipulbhale/config"

	log "github.com/sirupsen/logrus"
)

var (
	routeRepo map[string]*route
	regex     *regexp.Regexp
	//	patternRoute 		= 	"^(\\/)([/a-zA-Z0-9])+\\s+\\w+.\\w+"
	patternRoute = "([A-Z])+\\s+([/a-zA-Z0-9])+\\s+\\w+.\\w+"
)

func init() {
	regex, _ = regexp.Compile(patternRoute)
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

type route struct {
	url        string
	controller string
	method     string
}

//func NewRoute(url string, controller string, method string) (r *route) {
//	return &route{url: url, controller: controller, method: method}
//}
//
//func AddRouteToRepo(r *route) {
//	routeRepo = make(map[string]*route)
//	routeRepo[r.url] = r
//}

func (r *route) GetURL() string {
	return r.url
}

func (r *route) GetController() string {
	return r.controller
}
func (r *route) GetMethod() string {
	return r.method
}

func GetRoute(url string, httpVerb string) (r *route) {
	log.Debug("Entering the GetRoute method.")
	var appURL string

	r = new(route)
	appURL = ""

	compiledPattern := regexp.MustCompile("\\s+")
	splitURL := strings.Split(url, "/")
	appContext := splitURL[1]
	log.Debugln("Appcontext is :: ", appContext)

	for i := 2; i < len(splitURL); i++ {
		appURL = appURL + "/" + splitURL[i]
	}

	// appSrcRoot, _ := os.Getwd()
	// appRouteCfgFile := filepath.Join(appSrcRoot, "gokul", "src", "github.com", gokul.APPS_SRC_ROOT, appContext, "/config/routes.cfg")
	appRouteCfgFile := filepath.Join(config.Cfg["apps.directory"], appContext, "/config/routes.cfg")

	log.Debug("appRouteCfgFile is :: ", appRouteCfgFile)

	appcfgInputFile, cfgInputError := os.Open(appRouteCfgFile)
	if cfgInputError != nil {
		log.Fatal("Error reading the config file for app. Exiting.")
		os.Exit(1)
	}
	defer appcfgInputFile.Close()

	cfgInputReader := bufio.NewReader(appcfgInputFile)
	for {
		routeLineString, cfgReaderError := cfgInputReader.ReadString('\n')
		if regex.MatchString(routeLineString) {
			routeLine := compiledPattern.Split(routeLineString, -1)
			httpMethod := strings.TrimSpace(routeLine[0])
			completeURL := strings.TrimSpace(routeLine[1])
			controllerAndMethod := strings.TrimSpace(routeLine[2])
			log.Info("HttpMethod %s :: completeURL %s :: controllerAndMethod :: %s\n", httpMethod, completeURL, controllerAndMethod)
			if strings.Compare(appURL, completeURL) == 0 {
				if strings.Compare(httpVerb, httpMethod) == 0 {
					//r = new(route)
					r.url = completeURL
					r.method = strings.Split(controllerAndMethod, ".")[1]
					r.controller = strings.Split(controllerAndMethod, ".")[0]
				}

			}
		}
		if cfgReaderError == io.EOF {
			break
		}
	}
	return r
}
`

const SERVER_TEMPLATE = `
package server

import (
	"net/http"
	//"time"

	"net"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/{{.AppNameForTemplate}}/controller"
	goreflect "github.com/{{.AppNameForTemplate}}/reflect"
	"github.com/{{.AppNameForTemplate}}/routes"
	"github.com/{{.AppNameForTemplate}}/config"


	"github.com/vipulbhale/gokul/controller"
	goreflect "github.com/vipulbhale/gokul/reflect"
	"github.com/vipulbhale/gokul/routes"
	"github.com/vipulbhale/gokul/config"

	log "github.com/sirupsen/logrus"
)

var (
	httpServer      *http.Server
	baseControllers *controller.BaseController
)

func init() {
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
	baseControllers = new(controller.BaseController)
}

type server struct {
	cfg map[string]string
}

func (s *server) GetConfig() map[string]string {
	return s.cfg
}

// Creates the new server
func NewServer(cfgFileLocation string) *server {
	log.Debugln("Entering the NewServer constructor.")
	if len(cfgFileLocation) > 0 {
		config.LoadConfigFile(cfgFileLocation)
	}
	return &server{cfg: config.Cfg}
}

func (s *server) ScanAppsForControllers(appName string) {
	log.Debugln("Entering the ScanAppsForController function.")
	if len(appName) != 0 {
		goreflect.ScanAppsDirectory(config.Cfg, appName)
	} else {
		goreflect.ScanAppsDirectory(config.Cfg, "")
	}
}

// This method handles all requests.
func handle(w http.ResponseWriter, r *http.Request) {
	log.Debug("Inside the handle method for the request")
	var maxRequestSize int64
	var err error

	if maxRequestSize, err = strconv.ParseInt(config.Cfg["http.maxrequestsize"], 10, 64); err != nil {
		log.Fatalln("Error parsing the maxrequestsize . Exiting")
		os.Exit(1)
	}

	if maxRequestSize > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	}

	log.Debug("Request URL is :: " + r.URL.Path)

	if r.URL.Path != "/favicon.ico" {
		if filteredRoute := routes.GetRoute(r.URL.Path, r.Method); filteredRoute != nil {
			log.Debugln("Filtered route is %v\n", *filteredRoute)
			//path := strings.Split(r.URL.Path, "/")
			log.Debugln(filteredRoute.GetController())
			log.Debugln(filteredRoute.GetMethod())
			log.Debugln(filteredRoute.GetURL())
			log.Debugln(reflect.ValueOf(filteredRoute.GetController()))
		}
	}
}

func Run(s *server) {
	var readTimeOut int64
	var writeTimeOut int64
	var network string
	var err error

	address := s.cfg["server.address"] + ":" + s.cfg["server.port"]
	network = "tcp"

	if readTimeOut, err = strconv.ParseInt(config.Cfg["timeout.read"], 10, 64); err != nil {
		log.Debugln("The value of readtimeout received is ", config.Cfg["timeout.read"])
		log.Fatalln("Error parsing the read timeout. Exiting...")
		os.Exit(1)
	}

	if writeTimeOut, err = strconv.ParseInt(config.Cfg["timeout.write"], 10, 64); err != nil {
		log.Fatalln("Error parsing the write timeout. Exiting...")
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:         address,
		Handler:      http.HandlerFunc(handle),
		ReadTimeout:  time.Duration(readTimeOut) * time.Second,
		WriteTimeout: time.Duration(writeTimeOut) * time.Second,
	}

	listener, err := net.Listen(network, httpServer.Addr)

	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	log.Fatalln("Failed to serve:", httpServer.Serve(listener))

}
`
const SERVER_CONFIG_TEMPLATE = `package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Cfg map[string]string
)

func init() {
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)
	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

// Load the config file for the server
//
func LoadConfigFile(cfgFile string) {
	log.Debugln("Input to LoadConfigFile function is :: ", cfgFile)
	serverConfig := loadConfig(cfgFile)

	if len(serverConfig) > 0 {
		Cfg = serverConfig
	}
}

func readConfig(filename, dirname string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(dirname)
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}

// Load the config as default or from config file
func loadConfig(cfgFileName string) map[string]string {
	serverConfig := make(map[string]string)

	filename := strings.TrimSuffix(filepath.Base(cfgFileName), filepath.Ext(filepath.Base(cfgFileName)))
	dirname := filepath.Dir(cfgFileName)

	log.Debugln("Directory of config file is :: ", dirname)
	log.Debugln("Filename of config file is :: ", filename)

	v1, err := readConfig(filename, dirname, map[string]interface{}{
		"server": map[string]interface{}{
			"address": "0.0.0.0",
			"port":    9090,
		},
		"http": map[string]interface{}{
			"read.timeout":   0,
			"write.timeout":  0,
			"maxrequestsize": 999999,
		},
		"logging": map[string]interface{}{
			"level": "debug",
		},
		"apps" : map[string]interface{}{
			"directory" : "apps",
		},
	})

	if err != nil {
		panic(fmt.Errorf("Error when reading config: %v\n", err))
	}
	log.Debugln("Configuration is :: ", v1)

	serverConfig["server.port"] = strconv.Itoa(v1.Get("server.port").(int))
	serverConfig["server.address"] = v1.Get("server.address").(string)
	serverConfig["logging.level"] = v1.Get("logging.level").(string)
	serverConfig["timeout.read"] = strconv.Itoa(v1.Get("http.read.timeout").(int))
	serverConfig["timeout.write"] = strconv.Itoa(v1.Get("http.write.timeout").(int))
	serverConfig["http.maxrequestsize"] = strconv.Itoa(v1.Get("http.maxrequestsize").(int))
	serverConfig["apps.directory"] = v1.Get("apps.directory").(string)

	log.Debugln("The config map created is :: ", serverConfig)
	return serverConfig
}
`

const GOKUL_REFLECT_TEMPLATE = `package reflect

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	cntrlSpec      = new(controllerSpec)
	packageNameMap = make(map[string]string)
)

const MAIN = ` + "`" + `// GENERATED CODE - DO NOT EDIT
package controllerwrapper

import (
	"reflect"{{ range $index, $packageName := .PackageName }}
	"{{ $packageName }}"{{ end }}
)

var (
	mapOfControllerNameToControllerObj = make(map[string]reflect.Type)
)

func RegisterControllers(){
	{{range $index, $element := .ControllerName}}
    		{{ $element | ToLower }} := {{ $element }}{}
    		typeOfController := reflect.TypeOf({{ $element | ToLower }})
    		mapOfControllerNameToControllerObj[typeOfController.Name()] = typeOfController
	{{ end }}
}

func New(name string) (interface{}, bool) {
	t, ok := mapOfControllerNameToControllerObj[name]
	if !ok {
		return nil, false
	}
	v := reflect.New(t)
	return v.Interface(), true
}

` + "` " + `
 type controllerSpec struct {
	ControllerName []string
	PackageName    []string
}

func init() {
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)
	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

func ScanAppsDirectory(configuration map[string]string, appName string) {
	log.Debugln("Entering the ScanAppsDirectory.")
	log.Debugln("inputs are :: ", configuration)
	var appsHomeDirPath string
	// srcRoot, _ := os.Getwd()
	// log.Debugln("The srcRoot is :: ", srcRoot)
	if len(appName) != 0 {
		// appsHomeDirPath := filepath.Join(srcRoot, "src", "github.com", "apps", appName)
		appsHomeDirPath = filepath.Join(configuration["apps.directory"], appName)
	} else {
		appsHomeDirPath = filepath.Join(configuration["apps.directory"])
	}

	directoryList := []string{}

	err := filepath.Walk(appsHomeDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "controller" {
			log.Debugln("Path is ", path)
			log.Debugln("Directory is ", path)
			dir := filepath.Dir(path + "/" + info.Name())
			log.Debugln("Dir is ", dir)
			// packagename := strings.Split(path, filepath.Join(srcRoot, "src"))[1]
			packagename := strings.Split(path, appsHomeDirPath)[1]
			packagename = strings.Replace(packagename, "/", "", 1)
			log.Debugln("PackageName is ", packagename)
			packageNameMap[packagename] = packagename
			directoryList = append(directoryList, dir)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Debugln("Found directories are ", directoryList)

	for packageNameKey, _ := range packageNameMap {
		cntrlSpec.PackageName = append(cntrlSpec.PackageName, packageNameKey)
		log.Debugln("Package list is ", cntrlSpec.PackageName)
	}

	//srcRoot = filepath.Join(srcRoot, "gokul", "src", "github.com")
	//log.Debugln("srcRoot is ", srcRoot)

	for _, directoryController := range directoryList {
		makeControllers(directoryController)
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}
	tmpl, err := template.New("test").Funcs(funcMap).Parse(MAIN)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, cntrlSpec)
	if err != nil {
		panic(err)
	}

}

func makeControllers(srcRoot string) {
	log.Debugln("Entering the makeControllers method")
	var allNamed []*types.Object

	structMap := make(map[string]reflect.Type)
	log.Debugln("map is ", structMap)
	kPath := filepath.Join(srcRoot)
	log.Debugln("kpath is ", kPath)

	fset := token.NewFileSet()

	pkgs, e := parser.ParseDir(fset, kPath, func(f os.FileInfo) bool {
		return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go")
	}, 0)
	if e != nil {
		log.Debugln(e)
		return
	}
	log.Debugln("AllNamed is ", allNamed, structMap)

	log.Debugln("parsed package map is :: ", pkgs)
	for _, pkg := range pkgs {
		if pkg.Name == "controller" {
			log.Debugln("package name is ", pkg.Name)
			processPackage(pkg, pkg.Name)

		}
	}

}

func processPackage(pkg *ast.Package, packageName string) {
	log.Debugln("Entering the processPackage function")
	log.Debugln(pkg.Name)
	log.Debugln(pkg.Files)
	printASTVisitor := &PrintASTVisitor{}
	//controllers := make([]string,0)
	//printASTVisitor.cntrlSpec.packageControllers = make(map[string][]string)
	//printASTVisitor.cntrlSpec.ControllerName = make([]string, 0)
	//printASTVisitor.cntrlSpec.PackageName =	packageName
	//printASTVisitor.cntrlSpec.packageControllers[packageName] = controllers
	ast.Walk(printASTVisitor, pkg)

}

type PrintASTVisitor struct {
	info *types.Info
	//cntrlSpec controllerSpec
}

func (v *PrintASTVisitor) Visit(node ast.Node) ast.Visitor {
	// fmt.Println(v.info.Types)
	if node != nil {
		switch kk := node.(type) {

		case *ast.Package:
			{
				fmt.Println(kk.Name)

			}
		case *ast.TypeSpec:
			{
				log.Debugln("Name  of struct is :: " + kk.Name.Name)
				structType := kk.Type.(*ast.StructType)
				log.Debugln("hi there", structType)
				for _, field := range structType.Fields.List {
					log.Debugln(reflect.TypeOf(field.Type), " name is ", field.Names)
					fieldType := field.Type
					pkgName, typeName := func() (string, string) {
						// Drill through any StarExprs.
						for {
							if starExpr, ok := fieldType.(*ast.StarExpr); ok {
								fieldType = starExpr.X
								continue
							}
							break
						}

						// If the embedded type is in the same package, it is an Ident.
						if ident, ok := fieldType.(*ast.Ident); ok {
							return "", ident.Name
						}

						if selectorExpr, ok := fieldType.(*ast.SelectorExpr); ok {
							if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok {
								return pkgIdent.Name, selectorExpr.Sel.Name
							}
						}
						return "", ""
					}()
					if typeName == "BaseController" {
						log.Debugln("I am the man ", pkgName, typeName)
						//v.cntrlSpec.ControllerName = append(v.cntrlSpec.ControllerName, kk.Name.Name)
						cntrlSpec.ControllerName = append(cntrlSpec.ControllerName, kk.Name.Name)

						log.Debugln("ControllerSpec is :: ", cntrlSpec)

					}

				}

			}
		case *ast.GenDecl:
			{
				log.Debugln("Name  of struct is :: ", kk)

			}
		}
	}
	return v
}
`

//CreateTemplates Method creates the appName and templates in the directory mentioned.
func CreateTemplates(dirname, appName, cfgFileLocation string) {
	log.Debugln("Entering the CreateTemplates method.")

	// writeToFile(filepath.Join(dirname, "src", "github.com", appName, "controller"), "controller.go", appName, cfgFileLocation, CONTROLLER_TEMPLATE)
	// writeToFile(filepath.Join(dirname, "src", "github.com", appName, "routes"), "route.go", appName, cfgFileLocation, ROUTES_TEMPLATE)
	// writeToFile(filepath.Join(dirname, "src", "github.com", appName, "server"), "server.go", appName, cfgFileLocation, SERVER_TEMPLATE)
	// writeToFile(filepath.Join(dirname, "src", "github.com", appName, "config"), "config.go", appName, cfgFileLocation, SERVER_CONFIG_TEMPLATE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName, "config"), "server.yml", appName, cfgFileLocation, CONFIG_FILE_TEMPLATE)

	// writeToFileReflect(filepath.Join(dirname, "src", "github.com", appName, "reflect", "reflect.go"), GOKUL_REFLECT_TEMPLATE)

	writeToFile(filepath.Join(dirname, "src", "github.com", appName), "main.go", appName, cfgFileLocation, MAIN_PACKAGE)
	writeToFile(filepath.Join(dirname, "src", "github.com", appName), "variables.env", appName, cfgFileLocation, "GOPATH="+os.Getenv("GOPATH"))
	createBinAndPackageDirectory(filepath.Join(dirname, "bin"))
	createBinAndPackageDirectory(filepath.Join(dirname, "pkg"))

}

func writeToFileReflect(fileName, content string) {
	log.Debugln("For the reflect directory is ", filepath.Dir(fileName))
	if _, err := os.Stat(filepath.Dir(fileName)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(fileName), 0755)
		if err != nil {
			panic(err)
		}
	}

	d1 := []byte(content)
	err := ioutil.WriteFile(fileName, d1, 0644)
	if err != nil {
		log.Fatalln("Error while copying the reflect.go to apps directory")
	}
}

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
		log.Fatalln("An error occurred with file creation\n")
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
}

func createBinAndPackageDirectory(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			panic(err)
		}
	}
}
