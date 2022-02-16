package main

import (
	"flag"

	_ "net/http/pprof"

	"github.com/lwch/goredis/code/app"
	"github.com/lwch/runtime"
)

func main() {
	port := flag.Int("port", 6379, "listen port")
	flag.Parse()

	a := app.New()
	runtime.Assert(a.ListenAndServe(uint16(*port)))
}
