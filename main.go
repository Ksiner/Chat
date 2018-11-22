package main

import (
	"Chat/daemon" /// w/o go mod it doesn't work
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

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

func WaitForSignal(ctx context.Context) {
	<-ctx.Done()
	log.Panicf("Server is shutted down!")
}

func main() {
	cfg, err := parseConfigFile()
	if err != nil {
		log.Panicf("Error in main! %v \n", err.Error())
	}
	context, err := cfg.Start()
	if err != nil {
		log.Panicf("Error in main! %v \n", err.Error())
	}
	WaitForSignal(context)
}
