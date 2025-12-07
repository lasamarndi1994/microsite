package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func StartCron() {
	c := cron.New(cron.WithSeconds()) // Important!
	c.AddFunc("*/1 * * * * *", func() {
		fmt.Println("Cron Executed:", time.Now())
	})

	// Schedule user seeding (e.g., run once every day at midnight)
	// For testing purposes, you might want to run it more frequently or trigger it manually
	c.AddFunc("0 0 0 * * *", func() {
		SeedFakeUsers()
	})

	c.AddFunc("0 0 0 * * *", func() {
		SeedFakeMicrosites()
	})

	// Uncomment the following line to run it immediately on startup for verification
	// go SeedFakeUsers()
	go SeedFakeMicrosites()

	c.Start()
}
