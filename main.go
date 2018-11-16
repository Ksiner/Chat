package main

import (
	"Chat/daemon"
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq"
)

type Config struct {
	Listen string `json:"listen"`
	Dbconn string `json:"dbconn"`
}

func parseConfigFile() (*daemon.Config, error) {
	var cfg Config
	if byteCfg, err := ioutil.ReadFile("conf.json"); err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal(byteCfg, &cfg); err != nil {
			return nil, err
		}
		daemonCfg := &daemon.Config{Listen: cfg.Listen}
		daemonCfg.Db.CnString = cfg.Dbconn
		return daemonCfg, nil
	}
}

func main() {

	cfg, err := parseConfigFile()
	if err != nil {
		fmt.Printf("Error in main! %v \n", err.Error())
		return
	}
	if err := cfg.Start(); err != nil {
		fmt.Printf("Error in main! %v \n", err.Error())
		return
	}
}
