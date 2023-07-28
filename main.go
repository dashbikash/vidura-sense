package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/common-nighthawk/go-figure"
	"github.com/dashbikash/vidura-sense/internal/scheduler"
	"github.com/dashbikash/vidura-sense/internal/system"
)

var (
	log    = system.Logger
	config = system.Config
)

func main() {
	//greet()
	var wg sync.WaitGroup
	wg.Add(1)
	go scheduler.Start()

	wg.Wait()

	//restapi.Start()

}
func greet() {
	figure.NewFigure("Vidura Sense", "", true).Print()
	fmt.Println(strings.Repeat("= ", 42) + config.Application.Version + "\n")
}
