package main

import (
	"runtime/pprof"

	"github.com/dashbikash/vidura-sense/internal/apiserver"
)

var threadProfile = pprof.Lookup("threadcreate")

func main() {

	// Start profiling
	// f, err := os.Create("target/myprogram.prof")
	// if err != nil {

	// 	fmt.Println(err)
	// 	return

	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()
	// system.Log.Info(fmt.Sprintf("threads in starting: %d\n", threadProfile.Count()))

	// fmt.Printf(("threads after LookupHost: %d\n"), threadProfile.Count())
	apiserver.Start()
}
