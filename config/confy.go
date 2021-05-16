package config

import (
	"github.com/beardnick/mynvim/myfile"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SshServer struct {
	Account string
	Password string
}

type Config struct {
	Servers []SshServer
}

var Conf Config
var conf *toml.Tree

func defaultConfigFile() string  {
	return filepath.Join(os.Getenv("HOME"),".mynvim/config.toml")
}

func String(name string) string {
	return conf.Get(name).(string)
}

func Int(name string) int {
	return conf.Get(name).(int)
}

func Get(name string) interface{}  {
	return conf.Get(name)
}

//func Arrays(name string) []interface{}  {
//	return conf.Get(name).([]interface{})
//}

func DefaultLoad() (err error)  {
	err = myfile.EnsureFile(defaultConfigFile())
	if err != nil {
		return
	}
	data, err := ioutil.ReadFile(defaultConfigFile())
	if err != nil {
		return
	}
	conf, err = toml.Load(string(data))
	if err != nil {
		return
	}
	err = toml.Unmarshal(data,&Conf)
	return
}