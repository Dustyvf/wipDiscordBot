package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

// Functions for reading the config files and stuff

// Location of global config (might add a way to set this via a command flag or something in the future, but for now the config location will be hardcoded in)
var configLocation = "./config/config.toml"

// structs to define config layout
// may seperate parts of the config into seperate things to allow for more flexability and server specific config
type configStruct struct {
	Auth   confAuth
	Config confConfig
	Debug  confDebug
}

// auth stuff
type confAuth struct {
	Token string
	DB    confDB
}

// database stuff (to be implimented)
type confDB struct {
	DBType string
	User   string
	Pass   string
}

// general config stuff
type confConfig struct {
	Prefix []string
}

// debug config
type confDebug struct {
	Loglevel string
}

// global variable for global bot config
var gconf configStruct

// load conf
func initconf(file string) {
	var config configStruct
	if _, err := toml.DecodeFile(file, &config); err != nil {
		panic(err)
	}
	fmt.Println(config)
	gconf = config
}

// load config on init
func init() {
	initconf(configLocation)
}
