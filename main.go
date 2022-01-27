package main

import (
	"fmt"
	_map "gomem/map"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"
)

func main() {
	debug.SetGCPercent(10)
	f, err := os.Create("cpuprofile")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()


	f2, err := os.Create("memprofile")
	if err != nil {
		panic(err)
	}
	defer f2.Close()
	defer pprof.WriteHeapProfile(f2)

	//gobuntdb.TestBuntDb()
	//gomem.TestMemDb()
	//memory.TestMemory()

	fmt.Println("GC pause for startup: ", gcPause())
	_map.TestMap()
	fmt.Println("GC pause for warmup: ", gcPause())

	_map.TestMap()
	fmt.Println("GC pause for map: ", gcPause())
}

var previousPause time.Duration

func gcPause() time.Duration {
	runtime.GC()
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	pause := stats.PauseTotal - previousPause
	previousPause = stats.PauseTotal
	return pause
}

