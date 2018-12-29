package routes

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vipulbhale/gokul/server/util"

	"github.com/vipulbhale/gokul/server/config"
)

var (
	routeRepo    map[string]*route
	regex        *regexp.Regexp
	patternRoute = "([A-Z])+\\s+([/a-zA-Z0-9])+\\s+\\w+.\\w+"
	log          *logrus.Logger
)

func init() {
	regex, _ = regexp.Compile(patternRoute)
	// Output to stdout instead of the default stderr, could also be a file.
	// log.SetOutput(os.Stdout)

	// // Only log the debug severity or above.
	// log.SetLevel(log.DebugLevel)
	log = util.GetLogger()
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
	appRouteCfgFile := filepath.Join(config.Cfg["apps.directory"], "/config/routes.cfg")

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
			log.Info("HttpMethod :: completeURL :: controllerAndMethod :: ", httpMethod, completeURL, controllerAndMethod)
			if strings.Compare(appURL, completeURL) == 0 {
				if strings.Compare(httpVerb, httpMethod) == 0 {
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
