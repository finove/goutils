package vconfig

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Viper viper instance
var Viper *viper.Viper

var (
	pwd, appPath string
)

// LoadConfigFile load config from fileName
func LoadConfigFile(fileName string, create ...bool) (err error) {
	name, ext, cfgPath, paths := getConfigNameAndPath(fileName)
	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		viper.SetConfigName(name)
		viper.SetConfigType(ext)
		for _, ph := range paths {
			viper.AddConfigPath(ph)
		}
	}
	err = viper.ReadInConfig()
	if errors.As(err, &viper.ConfigFileNotFoundError{}) && len(create) > 0 && create[0] == true {
		if cfgPath == "" && len(paths) > 0 {
			viper.SetConfigFile(filepath.Join(paths[len(paths)-1], name))
		}
		err = viper.WriteConfig()
	}
	return
}

func getConfigNameAndPath(fileName string) (name, ext, cfgPath string, paths []string) {
	var basePath string
	if fileName != "" && filepath.IsAbs(fileName) {
		cfgPath = fileName
		return
	}
	if fileName != "" {
		basePath, name = filepath.Split(fileName)
		ext = filepath.Ext(name)
		if ext != "" {
			ext = ext[1:]
		} else {
			ext = "yaml"
		}
	} else {
		name = filepath.Base(appPath) + ".yaml"
		ext = "yaml"
	}
	paths = append(paths, filepath.Join(pwd, basePath))
	if pwd != filepath.Dir(appPath) {
		paths = append(paths, filepath.Join(filepath.Dir(appPath), basePath))
	}
	return
}

func init() {
	pwd, _ = os.Getwd()
	appPath, _ = filepath.Abs(os.Args[0])
	Viper = viper.GetViper()
}
