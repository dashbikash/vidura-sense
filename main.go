package main

import (
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/dashbikash/vidura-sense/internal/common"
	"github.com/dashbikash/vidura-sense/internal/requestor"
)

var config = common.GetConfig()

func main() {
	greet()
	//scheduler.Start()
	//restapi.Start()
	//requestor.RequestDemo()
	//mongo.QueryData()
	requestor.GetRobots()

}
func greet() {
	figure.NewFigure("Vidura Sense", "", true).Print()
	fmt.Println(strings.Repeat("= ", 42) + config.Application.Version + "\n")

}
