package goconfig

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

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
	gc.addFlag("show-config-debug", false, "prints the configuration with resolution steps")
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

// GetString gets the value as a string.
func (c *GoConfig) GetString(key string) string {
	c.mustExists(key)
	return viper.GetString(key)
}

// GetInt gets the value as a int.
func (c *GoConfig) GetInt(key string) int {
	c.mustExists(key)
	return viper.GetInt(key)
}

// GetInt64 gets the value as int64.
func (c *GoConfig) GetInt64(key string) int64 {
	c.mustExists(key)
	return viper.GetInt64(key)
}

// GetBool gets the value as bool.
func (c *GoConfig) GetBool(key string) bool {
	c.mustExists(key)
	return viper.GetBool(key)
}

// GetFloat64 gets the value as float64.
func (c *GoConfig) GetFloat64(key string) float64 {
	c.mustExists(key)
	return viper.GetFloat64(key)
}

func (c *GoConfig) mustExists(key string) {
	if !viper.IsSet(key) {
		panic(fmt.Errorf("key %s has not been added", key))
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
func (c *GoConfig) ReadConfig() {
	if c.flagsBinded {
		if viper.GetString("config-file") != "" {
			viper.SetConfigFile(viper.GetString("config-file"))
		}
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	c.executeActionFlagsIfAny()
}

func (c *GoConfig) executeActionFlagsIfAny() {
	if viper.GetBool("show-config") {
		c.PrintConfig()
		os.Exit(1)
	}

	if viper.GetBool("show-config-debug") {
		c.PrintDebugConfig()
		os.Exit(1)
	}
}

func (c *GoConfig) addFlag(key string, defaultValue interface{}, usage string) {
	switch defaultValue.(type) {
	case int:
		flag.Int(key, defaultValue.(int), usage)
	case int64:
		flag.Int64(key, defaultValue.(int64), usage)
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
