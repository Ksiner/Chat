package daemon

import (
	"context"
	"net"

	"github.com/Ksiner/Chat/db"
	"github.com/Ksiner/Chat/model"
	"github.com/Ksiner/Chat/ui"
)

type Config struct {
	Listen string
	Db     db.Config
}

func (cfg *Config) Start() (context.Context, error) {
	dbCmds, err := db.InitDbCmds(cfg.Db)
	if err != nil {
		return nil, err
	}
	m := model.New(dbCmds)
	l, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return nil, err
	}
	context, cancelFunc := context.WithCancel(context.Background())
	ui.Run(m, l, cancelFunc)
	return context, nil
}

/// Signal handling in libraries should be avoided except specialized signal handling libraries (ex. ssh)
/// Provide "graceful stop" channel or method in server. Signals should be handled from main package or some sort of package managing particular application
