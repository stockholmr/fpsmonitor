package auth

func (c *controller) initLog() {

}

func (c *controller) Fatal(err interface{}) {
	c.log.Fatalf("FATAL %s", err)
}

func (c *controller) Error(err interface{}) {
	c.log.Printf("ERROR %s", err)
}

func (c *controller) Warning(err interface{}) {
	c.log.Printf("WARN %s", err)
}

func (c *controller) Info(err interface{}) {
	c.log.Printf("INFO %s", err)
}

func (c *controller) Debug(err interface{}) {
	c.log.Printf("DEBUG %s", err)
}
