package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/dashbikash/vidura-sense/internal/requester"
)

func main() {

	// Start profiling
	f, err := os.Create("target/myprogram.prof")
	if err != nil {

		fmt.Println(err)
		return

	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	requester.SimpleRequest()

}
