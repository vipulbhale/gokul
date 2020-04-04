package util

import (
	"os"
	"path/filepath"
)

func PrintFile(path string, info os.FileInfo, err error) ([]string, error) {
	var directories []string

	if err != nil {
		return nil, err
	}
	if info.IsDir() && info.Name() == "controller" {
		dir := filepath.Base(path)
		directories = append(directories, dir)
	}
	return directories, nil
}

func CreateDirectory(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			Logger.Errorln("Error while creating the directory", dirName)
		}
	}
}
