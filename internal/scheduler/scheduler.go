package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func startCron() {
	c := cron.New(cron.WithSeconds()) // Important!
	c.AddFunc("*/10 * * * * *", func() {
		fmt.Println("Cron Executed:", time.Now())
	})

	c.Start()
}
