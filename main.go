package main

import (
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/dashbikash/vidura-sense/provider"
	"github.com/dashbikash/vidura-sense/requestor"
	_ "github.com/dashbikash/vidura-sense/restapi"
	_ "github.com/dashbikash/vidura-sense/scheduler"
)

var config = provider.GetConfig()

func main() {
	greet()
	//scheduler.Start()
	//restapi.Start()
	requestor.Request()
}
func greet() {
	figure.NewFigure(config.Application.Name, "", true).Print()
	fmt.Println(strings.Repeat(" = ", 23) + "Version " + config.Application.Version + "\n")
}
