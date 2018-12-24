package server

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vipulbhale/gokul/server/config"
	"github.com/vipulbhale/gokul/server/controller"
	"github.com/vipulbhale/gokul/server/routes"
)

var (
	httpServer                       *http.Server
	baseControllers                  *controller.BaseController
	mapControllerNameToControllerObj map[string]reflect.Value
	tempServer                       *server
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

//NewServer ... Creates the new server
func NewServer(cfgFileLocation string, mapOfControllerNameToControllerObj map[string]reflect.Value) *server {
	log.Debugln("Entering the NewServer constructor.")
	if len(cfgFileLocation) > 0 {
		config.LoadConfigFile(cfgFileLocation)
	}
	mapControllerNameToControllerObj = mapOfControllerNameToControllerObj
	tempServer = &server{cfg: config.Cfg}
	return tempServer
}

// This method handles all requests.
func handle(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Inside the handle method for the request")
	log.Debugln("Accept header in the request is :: ", r.Header["Accept"][0])
	var requestAcceptHeader = r.Header["Accept"][0]
	var maxRequestSize int64
	var err error

	if maxRequestSize, err = strconv.ParseInt(config.Cfg["http.maxrequestsize"], 10, 64); err != nil {
		log.Fatalln("Error parsing the maxrequestsize. Exiting...")
		os.Exit(1)
	}

	if maxRequestSize > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	}

	log.Debug("Request URL is :: " + r.URL.Path)

	if r.URL.Path != "/favicon.ico" {
		if filteredRoute := routes.GetRoute(r.URL.Path, r.Method); filteredRoute != nil {
			log.Debugln("Filtered route is ", *filteredRoute)
			log.Debugln(filteredRoute.GetController())
			log.Debugln(filteredRoute.GetMethod())
			log.Debugln(filteredRoute.GetURL())
			log.Debugln(reflect.ValueOf(filteredRoute.GetController()))
			log.Debugln("Controller type object :: ", mapControllerNameToControllerObj[filteredRoute.GetController()])
			log.Debugln("About to execute the controller method using the reflection...")
			response := mapControllerNameToControllerObj[filteredRoute.GetController()].MethodByName(filteredRoute.GetMethod()).Call([]reflect.Value{})

			if response != nil && len(response) != 2 {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Controller Method is returning more than 2 parameters."))
			}

			if response[1].Elem().Interface() == nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Nil ModelAndView Struct"))
			}

			log.Debugln("First Arg Type name is :: ", response[0].Type().Name())
			log.Debugln("Second Arg Type name is ::", response[1].Elem().Type().Name())

			if response[0].Type().Name() != "error" {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Controller Method should return first parameter as error interface."))
			}

			if response[1].Elem().Type().Name() != "ModelAndView" {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Controller Method should return second parameter as ModelAndView struct."))
			}
			modelAndView := response[1].Elem().Interface().(controller.ModelAndView)
			log.Debugln("Value of ModelAndView struct received is :: ", modelAndView)
			model := modelAndView.GetModel()
			view := modelAndView.GetView()
			responseType := modelAndView.GetResponseType()

			log.Debugln("Model to be passed to template is :: ", model)
			log.Debugln("View to be passed to template is :: ", view)
			log.Debugln("ResponseType of the response is :: ", responseType)

			if r.Header["Accept"] != nil && len(r.Header["Accept"]) != 0 && (strings.Contains(requestAcceptHeader, "*/*") || strings.Contains(requestAcceptHeader, "text/html") && strings.Contains(requestAcceptHeader, responseType)) && len(view) != 0 {
				w.Header().Set("Content-Type", responseType)
				templateFileName := filepath.Join(tempServer.GetConfig()["apps.directory"], "view", view+".html")
				tmpl := template.Must(template.ParseFiles(templateFileName))
				tmpl.Execute(w, model)
				defer recoverFromTemplateExecute()
			} else if strings.Contains(requestAcceptHeader, "application/json") && strings.Contains(requestAcceptHeader, responseType) {
				w.Header().Set("Content-Type", responseType)
				jsonEncoder := json.NewEncoder(w)
				jsonEncoder.SetEscapeHTML(true)
				if err := jsonEncoder.Encode(model); err != nil {
					log.Errorln("Error while marshalling json response for the request", filteredRoute.GetURL())
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
				}
			} else if strings.Contains(requestAcceptHeader, "application/xml") && strings.Contains(requestAcceptHeader, responseType) {
				w.Header().Set("Content-Type", responseType)
				xmlEncoder := xml.NewEncoder(w)
				if err := xmlEncoder.Encode(model); err != nil {
					log.Errorln("Error while marshalling xml response for the request", filteredRoute.GetURL())
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
				}
			}

		}

	}
}

// Run ... Runs the server
func Run(s *server) {
	log.Debugln("Entering the Server's Run method.")

	var readTimeOut int64
	var writeTimeOut int64
	var network string
	var err error

	//	templateFileLocation = filepath.Join(s.cfg["apps.directory"], "views")
	address := s.cfg["server.address"] + ":" + s.cfg["server.port"]
	network = "tcp"

	log.Debugln("Read Timeout for the http connection for the service is :: ", config.Cfg["timeout.read"])
	log.Debugln("Write Timeout for the http connection for the service is :: ", config.Cfg["timeout.write"])

	if readTimeOut, err = strconv.ParseInt(config.Cfg["timeout.read"], 10, 64); err != nil {
		log.Debugln("The value of readtimeout received is ", config.Cfg["timeout.read"])
		log.Fatalln("Error parsing the read timeout. Exiting...")
		os.Exit(1)
	}

	if writeTimeOut, err = strconv.ParseInt(config.Cfg["timeout.write"], 10, 64); err != nil {
		log.Debugln("The value of writetimeout received is ", config.Cfg["timeout.write"])
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
		log.Fatalln("Failed to listen :: ", err)
	}

	log.Fatalln("Failed to serve :: ", httpServer.Serve(listener))
}

func recoverFromTemplateExecute() {
	log.Debugln("Entering the method recoverFromTemplateExecute while executing template.")
	if err := recover(); err != nil {
		log.Errorln("Error while executing template :: ", err)
	}
}
