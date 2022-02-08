package main

import (
	"flag"
	"fmt"

	"net/http"
	_ "net/http/pprof"

	"github.com/lwch/goredis/code/app"
	"github.com/lwch/goredis/code/obj"
	"github.com/lwch/runtime"
)

func main() {
	port := flag.Int("port", 6379, "listen port")
	flag.Parse()

	go func() {
		go http.ListenAndServe(":9999", nil)
		dict := obj.NewDict()
		var i int
		for {
			dict.Set([]byte(fmt.Sprintf("%d", i)), i)
			i++
		}
	}()

	a := app.New()
	runtime.Assert(a.ListenAndServe(uint16(*port)))
}
