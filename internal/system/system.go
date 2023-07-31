package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

func init() {
	Config = getConfig()
	Log = getLogger()
	setEnvironment()
	greetMessage()
}

func setEnvironment() {
	os.Setenv("GOMEMLIMIT", "1024MiB")
}

func greetMessage() {
	figure.NewFigure("Vidura Sense", "", true).Print()
	fmt.Println(strings.Repeat("= ", 42) + Config.Application.Version + "\n")
}
