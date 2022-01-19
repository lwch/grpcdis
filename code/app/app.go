package app

import (
	"fmt"
	"net"

	"github.com/lwch/goredis/code/client"
	"github.com/lwch/goredis/code/command/server"
	"github.com/lwch/goredis/code/command/strings"
	"github.com/lwch/logging"
)

// App application
type App struct {
	cmds *server.Command
}

// New new application
func New() *App {
	app := &App{
		cmds: server.NewCommand(),
	}
	// strings
	app.cmds.Add(strings.NewSet())
	return app
}

// ListenAndServe listen and serve
func (app *App) ListenAndServe(port uint16) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			logging.Error("accept error: %v", err)
			continue
		}
		cli := client.New(conn, app.cmds)
		go cli.Run()
	}
}
