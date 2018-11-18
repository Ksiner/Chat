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

/// Signal handling in libraries should be avoided except specialized signal handling libraries (ex. ssh)
/// Provide "graceful stop" channel or method in server. Signals should be handled from main package or some sort of package managing particular application
func WaitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
