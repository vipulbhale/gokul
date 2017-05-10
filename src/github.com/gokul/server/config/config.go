package config

import (
	"bufio"
	"github.com/gokul"
	log "github.com/logrus"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	regex         *regexp.Regexp
	regexToIgnore *regexp.Regexp
	//pattern = "[#].*\\n|\\s+\\n|\\S+[=]|.*\n"
	pattern             = "[\\w.]+\\s+=\\s+[\\w]+"
	patternToIgnoreLine = "^#+.*"
	//srcRoot			string
	Cfg map[string]string
)

func init() {
	regex, _ = regexp.Compile(pattern)
	regexToIgnore, _ = regexp.Compile(patternToIgnoreLine)
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

func LoadConfigFile(cfgFile string) {
	log.Debugln("Input to LoadConfigFile function is ", cfgFile)
	serverConfig := make(map[string]string)

	srcRoot, _ := os.Getwd()
	log.Debugln("srcRoot is ", srcRoot)
	srcRoot = filepath.Join(srcRoot, "gokul/src/", gokul.GOKUL_SRC_ROOT, "/")

	log.Debugln("srcRoot is ", srcRoot)

	inputFile, inputError := os.Open(srcRoot + "/" + cfgFile)
	if inputError != nil {
		log.Fatal("Error reading the config file for server. Exiting.")
		os.Exit(1)
	}
	defer inputFile.Close()
	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readerError := inputReader.ReadString('\n')
		if regexToIgnore.MatchString(inputString) {
			continue
		} else if regex.MatchString(inputString) {
			pair := strings.Split(inputString, "=")
			if len(pair) == 2 {
				serverConfig[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
			}
		}
		if readerError == io.EOF {
			break
		}
	}
	if len(serverConfig) > 0 {
		Cfg = serverConfig
	}
}
