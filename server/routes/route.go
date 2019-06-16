package routes

import (
	"bufio"
	"fmt"
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
	routeRepo    map[string]*route
	regex        *regexp.Regexp
	patternRoute = "([A-Z])+\\s+([/a-zA-Z0-9])+\\s+\\w+.\\w+"
	log          *logrus.Logger
)

func init() {
	regex, _ = regexp.Compile(patternRoute)
	log = util.GetLogger()
}

type route struct {
	Url        string
	Controller string
	Method     string
}

type routeYaml struct {
	UriInfo []struct {
		Uri              string
		ControllerMethod string `yaml:"controller-method"`
		HttpMethod       string `yaml:"http-method"`
	}
}

func (r *route) GetURL() string {
	return r.Url
}

func (r *route) GetController() string {
	return r.Controller
}
func (r *route) GetMethod() string {
	return r.Method
}

func (r *route) SetURL(url string) {
	r.Url = url
}

func (r *route) SetController(controller string) {
	r.Controller = controller
}
func (r *route) SetMethod(method string) {
	r.Method = method
}

func getAppContext(url string) string {
	log.Debug("Entering the getAppContext method")
	splitURL := strings.Split(url, "/")
	appContext := splitURL[1]
	log.Debugln("Appcontext is :: {} ", appContext)
	return appContext
}

func GetRoute(url string, httpVerb string) (r *route) {
	log.Debugln("Entering the GetRoute method.")
	var appURL string

	r = new(route)
	appURL = ""

	splitURL := strings.Split(url, "/")
	appContext := getAppContext(url)
	log.Debugln("The extracted appContext is :: {}", appContext)
	for i := 2; i < len(splitURL); i++ {
		appURL = appURL + "/" + splitURL[i]
	}

	//appRouteCfgFile := filepath.Join(config.Cfg["apps.directory"], "/config/routes.cfg")
	appRouteCfgFile := filepath.Join(config.Cfg["apps.directory"], "config")

	log.Debug("appRouteCfgFile is :: {}", appRouteCfgFile)

	//find the type of the route file
	routeFileType := findTypeOfConfigFile(appRouteCfgFile)
	if strings.Compare(routeFileType, ".yaml") == 0 || strings.Compare(routeFileType, ".yml") == 0 {
		r = GetRouteFromYaml(appURL, httpVerb)
	} else if strings.Compare(routeFileType, ".cfg") == 0 {
		r = getRouteWhenFileIsCfgType(appRouteCfgFile, httpVerb, appURL, r)
	}
	log.Debugln("The route that was extracted based on the URI and VERB is {}", r)
	log.Debugln("Leaving the GetRoute method.")
	return r
}

func getRouteWhenFileIsCfgType(appRouteCfgFile string, httpVerb string, appURL string, r *route) (routeResult *route) {
	compiledPattern := regexp.MustCompile("\\s+")
	appcfgInputFile, cfgInputError := os.Open(appRouteCfgFile)
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
					r.Url = completeURL
					r.Method = strings.Split(controllerAndMethod, ".")[1]
					r.Controller = strings.Split(controllerAndMethod, ".")[0]
				}
			}
		}
		if cfgReaderError == io.EOF {
			break
		}
	}
	routeResult = r
	return routeResult
}

func findTypeOfConfigFile(routeFileDir string) (routeFileType string) {
	filesInfo, err := ioutil.ReadDir(routeFileDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range filesInfo {
		fmt.Println(file.Name())
		if strings.Compare(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())), "routes") == 0 {
			routeFileType = strings.ToLower(filepath.Ext(file.Name()))
			break
		}
	}
	return routeFileType
}

func GetRouteFromYaml(url string, httpVerb string) (r *route) {
	routeForYaml := readRoutesYaml()
	r = new(route)

	for i := 0; i < len(routeForYaml.UriInfo); i++ {
		if strings.Compare(strings.TrimSpace(url), routeForYaml.UriInfo[i].Uri) == 0 && strings.Compare(strings.ToLower(strings.TrimSpace(httpVerb)), strings.ToLower(routeForYaml.UriInfo[i].HttpMethod)) == 0 {
			r.SetController(routeForYaml.UriInfo[i].ControllerMethod)
			r.SetMethod(routeForYaml.UriInfo[i].HttpMethod)
			r.SetURL(routeForYaml.UriInfo[i].Uri)
		}
	}
	return r
}

func readRoutesYaml() (routeForYaml *routeYaml) {
	yamlFile, err := ioutil.ReadFile(filepath.Join(config.Cfg["apps.directory"], "/config/routes.yml"))
	if err != nil {
		log.Fatalln("Error while reading the app config file in yaml format. Please create the file in config directory of app. Exiting... ")
	}
	err = yaml.Unmarshal(yamlFile, &routeForYaml)
	if err != nil {
		log.Fatalln("Error while parsing the routes yaml file. Please try again with correct routes.yml file")
	}
	log.Debugln("Contents of parsed yaml file are {}", routeForYaml)
	return routeForYaml

}
