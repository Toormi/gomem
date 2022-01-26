package main

import (
	_ "go.uber.org/automaxprocs"
	_map "gomem/map"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	//gobuntdb.TestBuntDb()
	//gomem.TestMemDb()
	//memory.TestMemory()
	go http.ListenAndServe(":8080", nil)
	_map.TestMap()
}
