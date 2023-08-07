package scheduler

import (
	"fmt"
	"time"

	"github.com/dashbikash/vidura-sense/internal/scheduler/job"
	"github.com/dashbikash/vidura-sense/internal/system"
	"github.com/go-co-op/gocron"
)

func Start() {
	s := gocron.NewScheduler(time.Local)
	_, err := s.Every(30).Seconds().Do(func() {
		var stime, ftime time.Time

		stime = time.Now()
		job.CrawlNewHtmlPages()
		ftime = time.Now()
		system.Log.Info(fmt.Sprintf("Time elapsed %f", ftime.Sub(stime).Seconds()))
	})
	if err != nil {
		defer fmt.Println("Error Occured in Job")
	}

	s.StartAsync()
}
