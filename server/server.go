package server

import (
	"net/http"
	//"time"

	"net"
	"os"
	"reflect"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vipulbhale/gokul/server/config"
	"github.com/vipulbhale/gokul/server/controller"
	goreflect "github.com/vipulbhale/gokul/server/reflect"
	"github.com/vipulbhale/gokul/server/routes"
)

var (
	httpServer                       *http.Server
	baseControllers                  *controller.BaseController
	mapControllerNameToControllerObj map[string]reflect.Value
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
func NewServer(cfgFileLocation string, mapOfControllerNameToControllerObj map[string]reflect.Value) *server {
	log.Debugln("Entering the NewServer constructor.")
	if len(cfgFileLocation) > 0 {
		config.LoadConfigFile(cfgFileLocation)
	}
	mapControllerNameToControllerObj = mapOfControllerNameToControllerObj
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
			log.Debugln("Controller type o:: ", mapControllerNameToControllerObj[filteredRoute.GetController()])
			mapControllerNameToControllerObj[filteredRoute.GetController()].MethodByName(filteredRoute.GetMethod()).Call([]reflect.Value{})
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
