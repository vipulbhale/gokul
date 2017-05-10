package gokul

import (
	"net/http"
	//"time"
	"fmt"
	goreflect "github.com/gokul/reflect"
	"github.com/gokul/routes"
	"github.com/gokul/server/config"
	log "github.com/logrus"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
)

func init() {
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

var (
	httpServer *http.Server
)

type server struct {
	//ip    string
	//vhost string
	//port  int
	cfg map[string]string
}

func (s *server) GetConfig() map[string]string {
	return s.cfg
}

func NewServer() *server {
	config.LoadConfigFile("server/config/server.cfg")
	return &server{cfg: config.Cfg}
}

func (s *server) ScanAppsForControllers() {
	log.Debugln("Entering the ")
	goreflect.ScanAppsDirectory(s.cfg)
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
	//fmt.Println(r.Method, r.URL)
	if r.URL.Path != "/favicon.ico" {
		if filteredRoute := routes.GetRoute(r.URL.Path, r.Method); filteredRoute != nil {
			fmt.Printf("Filtered route is %v\n", *filteredRoute)
			//path := strings.Split(r.URL.Path, "/")
			log.Debug(filteredRoute.GetController())
			log.Debug(filteredRoute.GetMethod())
			log.Debug(filteredRoute.GetURL())
			log.Debug(reflect.ValueOf(filteredRoute.GetController()))

		}
	}
	log.Debug("Request URL is :: " + r.URL.Path)
}

func Run(s *server) {
	var readTimeOut int64
	var writeTimeOut int64
	var network string
	var err error

	//address := s.ip + ":" + strconv.Itoa(s.port)
	address := s.cfg["ip.address"] + ":" + s.cfg["server.port"]

	network = "tcp"

	if readTimeOut, err = strconv.ParseInt(config.Cfg["timeout.read"], 10, 64); err != nil {
		fmt.Println("Error parsing the read timeout. Exiting")
		os.Exit(1)
	}

	if writeTimeOut, err = strconv.ParseInt(config.Cfg["timeout.write"], 10, 64); err != nil {
		fmt.Println("Error parsing the write timeout. Exiting")
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
