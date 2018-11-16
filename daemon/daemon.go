package daemon

import (
	"Chat/db"
	"Chat/model"
	"Chat/ui"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Listen string
	Db     db.Config
}

func (cfg *Config) Start() error {
	dbCmds, err := db.InitDbCmds(cfg.Db)
	if err != nil {
		return err
	}
	m := model.New(dbCmds)
	l, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return err
	}
	ui.Run(m, l)
	WaitForSignal()
	return nil
}

func WaitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
