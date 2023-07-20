package main

import (
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	_ "github.com/dashbikash/vidura-sense/datastore/mongo"
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
	requestor.RequestDemo()
	//mongo.QueryData()
}
func greet() {
	figure.NewFigure("Vidura Sense", "", true).Print()
	fmt.Println(strings.Repeat("= ", 42) + config.Application.Version + "\n")

}
