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
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/vipulbhale/gokul/server/util"

	"github.com/sirupsen/logrus"
	"github.com/vipulbhale/gokul/server/config"
	"github.com/vipulbhale/gokul/server/controller"
	"github.com/vipulbhale/gokul/server/routes"
)

type acceptHeaderWithQuality struct {
	mimeTypes string
	quality   float64
}

type byQuality []acceptHeaderWithQuality

func (a byQuality) Len() int           { return len(a) }
func (a byQuality) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byQuality) Less(i, j int) bool { return a[i].quality > a[j].quality }

var (
	httpServer                       *http.Server
	mapControllerNameToControllerObj map[string]reflect.Value
	tempServer                       *server
	log                              *logrus.Logger
)

func init() {
	// Output to stdout instead of the default stderr, could also be a file.
	// log.SetOutput(os.Stdout)
	// Only log the debug severity or above.
	// log.SetLevel(log.DebugLevel)
	log = util.GetLogger()
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

	var maxRequestSize int64
	var err error
	var requestAcceptHeader string
	acceptHeaderFound := false

	if r.Header["Accept"] == nil || len(r.Header["Accept"][0]) == 0 || r.Header["Accept"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Accept Header is not passed"))
	} else {
		requestAcceptHeader = r.Header["Accept"][0]
	}
	var sortedAcceptHeaders byQuality = getSortedAcceptHeader(requestAcceptHeader)

	if maxRequestSize, err = strconv.ParseInt(tempServer.GetConfig()["http.maxrequestsize"], 10, 64); err != nil {
		log.Fatalln("Error parsing the maxrequestsize. Exiting...")
		os.Exit(1)
	}

	if maxRequestSize > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	}

	log.Debug("Request URL is :: " + r.URL.Path)

	if r.URL.Path != "/favicon.ico" {
		if filteredRoute := routes.GetRoute(r.URL.Path, r.Method); filteredRoute != nil {
			log.Debugln("Filtered Route is ", *filteredRoute)
			log.Debugln("Filtered Route Controller is :: ", filteredRoute.GetController())
			log.Debugln("Filter Route Controller's method is :: ", filteredRoute.GetMethod())
			log.Debugln("Filtered Route URL is :: ", filteredRoute.GetURL())
			log.Debugln("Controller type object :: ", mapControllerNameToControllerObj[filteredRoute.GetController()])
			log.Debugln("About to execute the controller method using the reflection...")
			response := mapControllerNameToControllerObj[filteredRoute.GetController()].MethodByName(filteredRoute.GetMethod()).Call([]reflect.Value{})
			// response := reflect.New(reflect.TypeOf(mapControllerNameToControllerObj[filteredRoute.GetController()])).Elem().MethodByName(filteredRoute.GetMethod()).Call([]reflect.Value{})

			if response != nil && len(response) != 2 {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Controller Method is returning more than 2 parameters."))
			} else if response[1].Elem().Interface() == nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Nil ModelAndView Struct"))
			} else if response[0].Type().Name() != "error" {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Controller Method should return first parameter as error interface."))
			} else if response[1].Elem().Type().Name() != "ModelAndView" {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Controller Method should return second parameter as ModelAndView struct."))
			} else if response[0].Interface() != nil && response[0].Interface().(error) != nil {
				err := response[0].Interface().(error)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			} else {
				modelAndView := response[1].Elem().Interface().(controller.ModelAndView)
				log.Debugln("Value of ModelAndView struct received is :: ", modelAndView)

				model := modelAndView.GetModel()
				view := modelAndView.GetView()

				log.Debugln("Model to be passed to template is :: ", model)
				log.Debugln("View to be passed to template is :: ", view)

				for _, acceptHeader := range sortedAcceptHeaders {
					requestAcceptHeader = acceptHeader.mimeTypes
					if (strings.Contains(requestAcceptHeader, "*/*") || strings.Contains(requestAcceptHeader, "text/html")) && view != "" {
						w.Header().Set("Content-Type", "text/html")
						templateFileName := filepath.Join(tempServer.GetConfig()["apps.directory"], "view", view+".html")
						tmpl := template.Must(template.ParseFiles(templateFileName))
						tmpl.Execute(w, model)
						defer recoverFromTemplateExecute()
						acceptHeaderFound = true
						break
					} else if strings.Contains(requestAcceptHeader, "application/json") && view == "" {
						w.Header().Set("Content-Type", requestAcceptHeader)
						jsonEncoder := json.NewEncoder(w)
						jsonEncoder.SetEscapeHTML(true)
						if err := jsonEncoder.Encode(model); err != nil {
							log.Errorln("Error while marshalling json response for the request", filteredRoute.GetURL())
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte(err.Error()))
						}
						acceptHeaderFound = true
						break
					} else if strings.Contains(requestAcceptHeader, "application/xml") && view == "" {
						w.Header().Set("Content-Type", requestAcceptHeader)
						xmlEncoder := xml.NewEncoder(w)
						if err := xmlEncoder.Encode(model); err != nil {
							log.Errorln("Error while marshalling xml response for the request", filteredRoute.GetURL())
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte(err.Error()))
						}
						acceptHeaderFound = true
						break
					}
				}
				if !acceptHeaderFound {
					log.Errorln("Didn't match with content type with Accept header for the given route.")
					w.WriteHeader(http.StatusNotAcceptable)
					w.Write([]byte("Accept header may be not provided"))
				}
			}
		} else {
			log.Errorln("Didn't match the URL inside the application.")
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("404 - Route not found"))
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
	address := s.GetConfig()["server.address"] + ":" + s.GetConfig()["server.port"]
	network = "tcp"

	log.Debugln("Read Timeout for the http connection for the service is :: ", s.GetConfig()["timeout.read"])
	log.Debugln("Write Timeout for the http connection for the service is :: ", s.GetConfig()["timeout.write"])

	if readTimeOut, err = strconv.ParseInt(s.GetConfig()["timeout.read"], 10, 64); err != nil {
		log.Debugln("The value of readtimeout received is ", s.GetConfig()["timeout.read"])
		log.Fatalln("Error parsing the read timeout. Exiting...")
		os.Exit(1)
	}

	if writeTimeOut, err = strconv.ParseInt(s.GetConfig()["timeout.write"], 10, 64); err != nil {
		log.Debugln("The value of writetimeout received is ", s.GetConfig()["timeout.write"])
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

func getSortedAcceptHeader(acceptHeader string) []acceptHeaderWithQuality {
	log.Debugln("The raw accept header received from client in string format is :: ", acceptHeader)
	return sortAcceptHeaderByQuality(parseAcceptHeader(acceptHeader))
}

func parseAcceptHeader(acceptHeader string) []acceptHeaderWithQuality {
	var regexPatternForAccept = "(([a-z+\\*]+\\/[a-z+\\*]+,?\\s?)+(;q=)?[01]?\\.?[0-9]?(,\\s)?)+"
	searchedStrings := regexp.MustCompile(regexPatternForAccept).FindAllString(acceptHeader, -1)
	log.Debugln("The accept header after matching to regex :: ", searchedStrings)
	ahs := make([]acceptHeaderWithQuality, len(searchedStrings))

	for i, header := range searchedStrings {
		tempHeader := strings.Split(header, ";")
		log.Infoln("Header after splitting with token as ;", tempHeader)
		mimeTypes := tempHeader[0]
		ahs[i].mimeTypes = mimeTypes
		if len(tempHeader) > 1 {
			if quality, err := strconv.ParseFloat(tempHeader[1], 64); err == nil {
				ahs[i].quality = quality
			} else if err != nil {
				log.Errorln("Error while parsing the quality :: ", err)
				log.Infoln("Defaulting the quality to accept header to 1")
				ahs[i].quality = 1
			}
		} else {
			ahs[i].quality = 1
		}
		log.Debugln("One of the accept header is :: ", ahs[i])
	}
	log.Debugln("Parsed Accept Header is :: ", ahs)
	return ahs
}

func sortAcceptHeaderByQuality(parsedAcceptHeader []acceptHeaderWithQuality) []acceptHeaderWithQuality {
	sort.Sort(byQuality(parsedAcceptHeader))
	log.Debugln("Sorted Accept Header list based on quality is :: ", parsedAcceptHeader)
	return parsedAcceptHeader
}

func recoverFromTemplateExecute() {
	log.Debugln("Entering the method recoverFromTemplateExecute while executing template.")
	if err := recover(); err != nil {
		log.Errorln("Error while executing template :: ", err)
	}
}
