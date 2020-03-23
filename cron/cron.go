package cron

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
)

//Cron main type
type Cron struct {
	Cron *cron.Cron
	TZ   string
	Time string
}

// DB instance
var DB *gorm.DB

// New cron
func (c *Cron) New() {
	cron := cron.New(
		cron.WithLogger(cron.DefaultLogger),
	)
	format := fmt.Sprintf("CRON_TZ=%s %s", c.TZ, c.Time)
	cron.AddFunc(format, SendNotifications())
	log.Printf("CRON # : %s\n", format)

	c.Cron = cron
}

// Start cron
func (c *Cron) Start() {
	c.Cron.Start()
}

// Stop cron
func (c *Cron) Stop() {
	c.Cron.Stop()
}

// Entries returns cron entries
func (c *Cron) Entries() []cron.Entry {
	return c.Cron.Entries()
}
