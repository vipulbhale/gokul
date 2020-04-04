package routes

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/sirupsen/logrus"
	"github.com/vipulbhale/gokul/server/util"

	"github.com/vipulbhale/gokul/server/config"
)

var (
	routeRepo    map[string]*Route
	regex        *regexp.Regexp
	patternRoute = "([A-Z])+\\s+([/a-zA-Z0-9])+\\s+\\w+.\\w+"
	log          *logrus.Logger
	routeParsed  *routeYaml
)

func init() {
	regex, _ = regexp.Compile(patternRoute)
	log = util.GetLogger()
	routeParsed = new(routeYaml)
}

// Route structure holds the information about the URL, Controller and Method for it.
type Route struct {
	URL        string
	Controller string
	Method     string
}

type routeYaml struct {
	URIInfo []struct {
		URI              string `yaml:"uri"`
		ControllerMethod string `yaml:"controller-method"`
		HTTPMethod       string `yaml:"http-method"`
	} `yaml:"uriInfo"`
}

// GetURL method returns the URL in Route Structure
func (r *Route) GetURL() string {
	return r.URL
}

// GetController method returns the Controller in Route Structure
func (r *Route) GetController() string {
	return r.Controller
}

// GetMethod method returns the Method in Route Structure
func (r *Route) GetMethod() string {
	return r.Method
}

//SetURL method sets the URL in Route Structure.
func (r *Route) SetURL(url string) {
	r.URL = url
}

// SetController method sets the controller in Route Structure
func (r *Route) SetController(controller string) {
	r.Controller = controller
}

// SetMethod method sets the method in Route Structure
func (r *Route) SetMethod(method string) {
	r.Method = method
}

//getAppContext method takes the URL as input and returns the application context component of the URL.
func getAppContext(url string) string {
	log.Debug("Entering the getAppContext method")
	splitURL := strings.Split(url, "/")
	appContext := splitURL[1]
	log.Debugln("Appcontext is ::  ", appContext)
	log.Debug("Leaving the getAppContext method")
	return appContext
}

// GetRoute function takes the url and http verb of request as input returns the pointer to route as *Route
// checks the URL in the file. Returns the route structure.
func GetRoute(url string, httpVerb string) (r *Route) {
	log.Debugln("Entering the GetRoute method.")
	var appURL string

	r = new(Route)
	appURL = ""

	splitURL := strings.Split(url, "/")
	appContext := getAppContext(url)
	log.Debugln("The extracted appContext is :: ", appContext)
	for i := 2; i < len(splitURL); i++ {
		appURL = appURL + "/" + splitURL[i]
	}

	//appRouteCfgFile := filepath.Join(config.Cfg["apps.directory"], "/config/routes.cfg")
	appRouteCfgFile := filepath.Join(config.Cfg["apps.directory"], "config")
	log.Debugln("appRouteCfgFile is :: ", appRouteCfgFile)

	//find the type of the route file
	routeFileType := findTypeOfConfigFile(appRouteCfgFile)
	log.Debugln("The route file type is :: ", routeFileType)

	if strings.Compare(routeFileType, ".yaml") == 0 || strings.Compare(routeFileType, ".yml") == 0 {
		log.Debugln("The route file type is yaml")
		log.Debugln("Getting the route corresponding to this URL :: ", url, " and HTTP Verb :: ", httpVerb)
		r = getRouteFromYaml(appURL, httpVerb)
	} else if strings.Compare(routeFileType, ".cfg") == 0 {
		log.Debugln("The route file type is config/cfg")
		r = getRouteWhenFileIsCfgType(appURL, httpVerb)
	}
	log.Debugln("The route that was extracted based on the URI and VERB is :: ", r)
	log.Debugln("Leaving the GetRoute method.")
	return r
}

func findTypeOfConfigFile(routeFileDir string) (routeFileType string) {
	log.Debugln("Entering the findTypeOfConfigFile method.")
	log.Debugln("The routeFileDir as input to this method is ::  ", routeFileDir)
	filesInfo, err := ioutil.ReadDir(routeFileDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range filesInfo {
		if strings.Compare(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())), "routes") == 0 {
			log.Debugln("Found out the file type after the match")
			routeFileType = strings.ToLower(filepath.Ext(file.Name()))
			break
		}
	}
	log.Debugln("The selected route config file type is :: ", routeFileType)
	log.Debugln("Leaving the findTypeOfConfigFile method.")
	return routeFileType
}

func getRouteWhenFileIsCfgType(appURL string, httpVerb string) (r *Route) {
	log.Debugln("Entering the getRouteWhenFileIsCfgType method.")
	r = new(Route)
	compiledPattern := regexp.MustCompile("\\s+")
	appcfgInputFile, cfgInputError := os.Open(filepath.Join(config.Cfg["apps.directory"], "config", "routes.cfg"))
	if cfgInputError != nil {
		log.Fatalln("Error reading the config file for app. Exiting.")
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
					r.URL = completeURL
					r.Controller = strings.Split(controllerAndMethod, ".")[0]
					r.Method = strings.Split(controllerAndMethod, ".")[1]
				}
			}
		}
		if cfgReaderError == io.EOF {
			break
		}
	}
	log.Debugln("Leaving the getRouteWhenFileIsCfgType method.")
	return r
}

func getRouteFromYaml(url, httpVerb string) (r *Route) {
	log.Debugln("Entering the getRouteFromYaml method.")
	log.Debugln("Input to this method are url :: %v :: httpverb is :: %v", url, httpVerb)
	routeForYaml := readRoutesYaml(routeParsed)
	r = new(Route)
	// Iterate through all routes.
	for i := 0; i < len(routeForYaml.URIInfo); i++ {
		log.Debugln("The route to be compared is :: ", routeForYaml.URIInfo[i])
		// Calling the patternMatch
		if patternMatchingTheURL(url, routeForYaml.URIInfo[i].URI) && strings.Compare(strings.ToLower(strings.TrimSpace(httpVerb)), strings.ToLower(routeForYaml.URIInfo[i].HTTPMethod)) == 0 {
			log.Debugln("There is match for url and httpverb")
			r.SetURL(routeForYaml.URIInfo[i].URI)
			r.SetController(strings.Split(routeForYaml.URIInfo[i].ControllerMethod, ".")[0])
			r.SetMethod(strings.Split(routeForYaml.URIInfo[i].ControllerMethod, ".")[1])
			log.Debugln("The selected route is ", r)
			break
		}
	}

	log.Debugln("Leaving the getRouteFromYaml method.")
	return r
}

func patternMatchingTheURL(url, urlPattern string) (patternMatch bool) {
	splitURL := strings.Split(url, "/")
	splitURLPattern := strings.Split(urlPattern, "/")

	if len(splitURL) != len(splitURLPattern) {
		return false
	}
	for i := 0; i < len(splitURL); i++ {
		if strings.HasPrefix(splitURLPattern[i], "{") && strings.HasSuffix(splitURLPattern[i], "}") {
			continue
		} else if strings.Compare(splitURL[i], splitURLPattern[i]) == 0 {
			continue
		} else {
			return false
		}
	}
	return true
}

func readRoutesYaml(routeParsed *routeYaml) (routeForYaml *routeYaml) {
	if routeParsed != nil && len(routeParsed.URIInfo) == 0 {
		yamlFile, err := ioutil.ReadFile(filepath.Join(config.Cfg["apps.directory"], "config", "routes.yml"))
		if err != nil {
			log.Fatalln("Error while reading the app config file in yaml format. Please create the file in config directory of app. Exiting... ")
		}
		log.Debugln("The yaml file for routes is :: ", yamlFile)
		err = yaml.Unmarshal(yamlFile, &routeParsed)
		if err != nil {
			log.Fatalln("Error while parsing the routes yaml file. Please try again with correct routes.yml file")
		}
	}

	log.Debugln("Contents of parsed yaml file are :: ", &routeParsed)
	return routeParsed

}
