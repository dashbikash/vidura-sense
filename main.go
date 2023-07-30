package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/dashbikash/vidura-sense/internal/apiserver"
	"github.com/dashbikash/vidura-sense/internal/requestor"
	"github.com/dashbikash/vidura-sense/internal/system"
)

func main() {
	system.Setup()
	system.GreetMessage()

	// Start profiling
	f, err := os.Create("target/myprogram.prof")
	if err != nil {

		fmt.Println(err)
		return

	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	//scheduler.Start()
	requestor.RequestDemo2()
	apiserver.Start()

	//requestor.RequestDemo1()

}
