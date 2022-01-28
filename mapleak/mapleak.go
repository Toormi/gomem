package main

import (
	"fmt"
	"runtime"
)

func main() {
	v := struct{}{}

	a := make(map[int]struct{})

	for i := 0; i < 10000; i++ {
		a[i] = v
	}

	runtime.GC()
	fmt.Printf("After Map Add 100000 %d\n", len(a))
	printMemStats("After Map Add 100000")

	for i := 0; i < 10000; i++ {
		delete(a, i)
	}

	runtime.GC()
	fmt.Printf("After Map Delete 10000 len%d\n", len(a))
	printMemStats("After Map Delete 10000")
	if len(a) <= 0 {
		a = make(map[int]struct{})
		runtime.GC()
		fmt.Printf("After New Map len %d\n", len(a))
		printMemStats("After New Map")
	}

	for i := 0; i < 10000-1; i++ {
		a[i] = v
	}

	runtime.GC()
	fmt.Printf("After Map Add 9999 again %d\n", len(a))
	printMemStats("After Map Add 9999 again")

	a = nil
	runtime.GC()
	fmt.Printf("After Map Set nil %d\n", len(a))
	printMemStats("After Map Set nil")
}

func printMemStats(mag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%v：memory = %vKB, GC Times = %v\n", mag, m.Alloc/1024, m.NumGC)
}
