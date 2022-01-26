package main

import (
	_map "gomem/map"
	"os"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("cpuprofile")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	//gobuntdb.TestBuntDb()
	//gomem.TestMemDb()
	//memory.TestMemory()

	_map.TestMap()
}
