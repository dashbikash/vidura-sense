package main

import (
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/dashbikash/vidura-sense/internal/requestor"
	"github.com/dashbikash/vidura-sense/internal/system"
)

var config = system.GetConfig()

func main() {
	greet()
	//scheduler.Start()
	//restapi.Start()
	//requestor.RequestDemo()
	//mongo.QueryData()

	requestor.RequestDemo2()
}
func greet() {
	figure.NewFigure("Vidura Sense", "", true).Print()
	fmt.Println(strings.Repeat("= ", 42) + config.Application.Version + "\n")
}
