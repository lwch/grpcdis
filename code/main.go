package main

import (
	"flag"

	"net/http"
	_ "net/http/pprof"

	"github.com/lwch/goredis/code/app"
	"github.com/lwch/runtime"
)

func main() {
	port := flag.Int("port", 6379, "listen port")
	flag.Parse()

	// for pprof
	go http.ListenAndServe(":9999", nil)

	a := app.New()
	runtime.Assert(a.ListenAndServe(uint16(*port)))
}
