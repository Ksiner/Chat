package main

import (
	"Chat/daemon" /// w/o go mod it doesn't work
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq" /// it should be a strong need to import package into current one
)

type Config struct {
	Listen string `json:"listen"`
	Dbconn string `json:"dbconn"`
}

func parseConfigFile() (*daemon.Config, error) {
	/// BTW it can be refactored in future
	/*
		if ...; err!=nil{
			return ...
		}
		if ...; err!=nil{
			return ...
		}
		return ...
	*/
	var cfg Config
	if byteCfg, err := ioutil.ReadFile("conf.json"); err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal(byteCfg, &cfg); err != nil {
			return nil, err
		}
		daemonCfg := &daemon.Config{Listen: cfg.Listen}
		daemonCfg.Db.CnString = cfg.Dbconn /// it can be done in declaration statement eralier
		return daemonCfg, nil
	}
}

func main() {
	cfg, err := parseConfigFile()
	if err != nil {
		fmt.Printf("Error in main! %v \n", err.Error())
		return
		/// it can be done with log.Panic() or simply panic()
	}
	if err := cfg.Start(); err != nil {
		fmt.Printf("Error in main! %v \n", err.Error())
		return
	}
}
