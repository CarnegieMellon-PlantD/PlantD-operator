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
		panic(fmt.Errorf("Cannot read config file: %s\n", err))
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
		fmt.Printf("Added load generator script \"%s\"\n", filenameNoExt)
	}
}

func GetInt(key string) int {
	mux.Lock()
	defer mux.Unlock()
	if !viperInstance.IsSet(key) {
		panic(fmt.Errorf("Key \"%s\" not found in config file\n", key))
	}
	return viperInstance.GetInt(key)
}

func GetInt32(key string) int32 {
	mux.Lock()
	defer mux.Unlock()
	if !viperInstance.IsSet(key) {
		panic(fmt.Errorf("Key \"%s\" not found in config file\n", key))
	}
	return viperInstance.GetInt32(key)
}

func GetInt64(key string) int64 {
	mux.Lock()
	defer mux.Unlock()
	if !viperInstance.IsSet(key) {
		panic(fmt.Errorf("Key \"%s\" not found in config file\n", key))
	}
	return viperInstance.GetInt64(key)
}

func GetString(key string) string {
	mux.Lock()
	defer mux.Unlock()
	if !viperInstance.IsSet(key) {
		panic(fmt.Errorf("Key \"%s\" not found in config file\n", key))
	}
	return viperInstance.GetString(key)
}

func GetStringMapString(key string) map[string]string {
	mux.Lock()
	defer mux.Unlock()
	if !viperInstance.IsSet(key) {
		panic(fmt.Errorf("Key \"%s\" not found in config file\n", key))
	}
	return viperInstance.GetStringMapString(key)
}
