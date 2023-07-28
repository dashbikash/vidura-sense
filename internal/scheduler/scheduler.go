package scheduler

import (
	"fmt"
	"time"

	"github.com/dashbikash/vidura-sense/internal/requestor"
	"github.com/dashbikash/vidura-sense/internal/system"
	"github.com/go-co-op/gocron"
)

var (
	log = system.Logger
)

func Start() {
	s := gocron.NewScheduler(time.Local)
	_, err := s.Every(20).Seconds().Do(func() {
		var stime, ftime time.Time

		stime = time.Now()
		requestor.RequestDemo1()
		ftime = time.Now()
		log.Info(fmt.Sprintf("Time elapsed %f", ftime.Sub(stime).Seconds()))
	})
	if err != nil {
		defer fmt.Println("Error Occured in Job")
	}

	s.StartAsync()
}
