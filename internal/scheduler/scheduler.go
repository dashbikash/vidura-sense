package scheduler

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func Start() {
	s := gocron.NewScheduler(time.Local)
	_, err := s.Every(30).Seconds().Do(func() {

	})
	if err != nil {
		defer fmt.Println("Error Occured in Job")
	}

	s.StartAsync()
}
