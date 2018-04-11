package goconfig

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// GoConfig hold the configuration.
type GoConfig struct {
	flagsBinded bool
}

// New creates a new Config.
func New() *GoConfig {

	// always add the current working directory
	viper.AddConfigPath(".")

	gc := &GoConfig{}

	// add default flags
	gc.addFlag("show-config", false, "prints the configuration")
	gc.addFlag("show-config", false, "prints the configuration with resolution steps")
	gc.addFlag("config-file", "", "configuration file to use")

	return gc
}

// SetConfigName sets the configuration filename.
func (c *GoConfig) SetConfigName(filename string) {
	viper.SetConfigName(filename)
}

// AddConfigurationPaths adds the paths where to find the configuration file.
func (c *GoConfig) AddConfigurationPaths(paths ...string) {
	for _, p := range paths {
		viper.AddConfigPath(p)
	}

}

// Add adds a configuration key with default values and usage information.
func (c *GoConfig) Add(key string, defaultValue interface{}, usage string) {
	// set defaults
	viper.SetDefault(key, defaultValue)

	// set flag
	c.addFlag(key, defaultValue, usage)
}

func (c *GoConfig) addFlag(key string, defaultValue interface{}, usage string) {
	switch defaultValue.(type) {
	case int:
		flag.Int(key, defaultValue.(int), usage)
	case int64:
		flag.Int64(key, defaultValue.(int64), usage)
	case uint:
		flag.Uint(key, defaultValue.(uint), usage)
	case uint64:
		flag.Uint64(key, defaultValue.(uint64), usage)
	case string:
		flag.String(key, defaultValue.(string), usage)
	case float64:
		flag.Float64(key, defaultValue.(float64), usage)
	case bool:
		flag.Bool(key, defaultValue.(bool), usage)
	default:
		panic(fmt.Errorf("unknow type:%T for value:%#v", defaultValue, defaultValue))
	}
}

// PrintConfig prints the configuration.
func (c *GoConfig) PrintConfig() {
	encoded, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	fmt.Println(string(encoded))
}

// PrintDebugConfig prints configuration with resolution steps.
func (c *GoConfig) PrintDebugConfig() {
	viper.Debug()
	fmt.Printf("Merged Config: \n%#v\n", viper.AllSettings())
}

// BindFlags binds configuration elements to flags.
func (c *GoConfig) BindFlags() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	c.flagsBinded = true
}

// ReadConfig reads the configuration.
func (c *GoConfig) ReadConfig() error {
	if c.flagsBinded {
		if viper.GetString("config-file") != "" {
			viper.SetConfigFile(viper.GetString("config-file"))
		}
	}
	return viper.ReadInConfig()
}
