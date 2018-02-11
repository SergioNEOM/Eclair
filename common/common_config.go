package common

import (
	"fmt"
	"os"
)

type MainConfigure struct {
	Port         string `json:"port"`
	LogFile      string `json:"logfile"`
	DebugLogFile string `json:"debuglogfile"`
	UsersFile    string `json:"usersfile"`
	CoursesFile  string `json:"coursesfile"`
	WorkDir      string `json:"workdir"`
}

var ConfigFileName string
var MainConf MainConfigure

func (m *MainConfigure) SetDefaultConf() {
	m.Port = ":9090"
	m.LogFile = APP_NAME + ".log"
	m.DebugLogFile = APP_NAME + ".dbg"
	m.UsersFile = APP_NAME + ".usr"
	m.CoursesFile = APP_NAME + ".crs"
	// m.WorkDir = ""
}

func (m *MainConfigure) LoadConf() bool {
	err := LoadJSON(ConfigFileName, m)
	if err == nil {
		return true
	}
	//todo:  log : fmt.Fprintf(os.Stderr, "*** Error on read config file: %s\n", err)
	fmt.Fprintf(os.Stderr, "*** Error on read config file: %s\n", err)
	//todo: надо ли заполнять данными по умолчанию?
	m.SetDefaultConf()
	fmt.Fprintln(os.Stderr, "*** Loaded default configure values!")
	return false
}

func (m *MainConfigure) SaveConf(filename string, permiss os.FileMode) bool {
	err := SaveJSON(m, filename, permiss)
	if err == nil {
		return true // success
	}
	fmt.Fprintf(os.Stderr, "*** Error on save config file: %s\n", err)
	//	Logs  Fprintf(os.Stderr, "*** Error on read config file: %s\n*** Loaded default configure values!", err)
	//todo: выводить в журнал
	return false
}
