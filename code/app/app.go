package app

import (
	"github.com/lwch/goredis/code/obj"
)

// App application
type App struct {
	objs *obj.Objs
}

// New new application
func New() *App {
	objs := obj.New()
	app := &App{
		objs: objs,
	}
	return app
}

// ListenAndServe listen and serve
func (app *App) ListenAndServe(port uint16) error {
	return nil
}
