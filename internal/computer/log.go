package computer

import (
	"time"
)

func (c *controller) initLog() {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	c.log.SetPrefix(
		timestamp,
	)
}

func (c *controller) Fatal(err error) {
	c.log.Fatalf("FATAL %s", err)
}

func (c *controller) Error(err error) {
	c.log.Fatalf("ERROR %s", err)
}

func (c *controller) Warning(err error) {
	c.log.Fatalf("WARN %s", err)
}

func (c *controller) Info(err error) {
	c.log.Fatalf("INFO %s", err)
}

func (c *controller) Debug(err error) {
	c.log.Fatalf("DEBUG %s", err)
}
