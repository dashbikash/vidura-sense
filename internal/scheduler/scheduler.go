package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

func Start() {
	s := gocron.NewScheduler(time.Local)
	// _, err := s.Every(5).Seconds().Do(func() {
	// 	crawler.Request("https://dummyjson.com/products")
	// })
	// if err != nil {
	// 	defer fmt.Println("Error Occured in Job")
	// }

	s.StartAsync()
}
