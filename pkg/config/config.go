package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	viperInstance *viper.Viper
	mux           sync.Mutex
)

func init() {
	viperInstance = viper.New()
	viperInstance.SetConfigName("config")
	viperInstance.SetConfigType("yaml")
	viperInstance.AddConfigPath("./config/plantd") // Development
	viperInstance.AddConfigPath("/etc/plantd")     // Production
	if err := viperInstance.ReadInConfig(); err != nil {
		fmt.Printf("Cannot read config file: %s\n", err)
	}

	// Load the load generator scripts
	lgScriptPath := "./apps/loadgen"
	lgScriptFileInfo, err := os.ReadDir(lgScriptPath)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		panic(fmt.Errorf("Cannot read load generator script directory: %s\n", err))
	}
	for _, file := range lgScriptFileInfo {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		filenameNoExt := filename[:len(filename)-len(filepath.Ext(filename))]
		content, err := os.ReadFile(filepath.Join(lgScriptPath, filename))
		if err != nil {
			panic(fmt.Errorf("Cannot read load generator script file: %s\n", err))
		}
		viperInstance.Set(fmt.Sprintf("loadGenerator.script.%s", filenameNoExt), string(content))
		fmt.Printf("Load generator script %s loaded\n", filenameNoExt)
	}
}

func GetViper() *viper.Viper {
	mux.Lock()
	defer mux.Unlock()
	return viperInstance
}
