package admin

func (c *controller) initLog() {

}

func (c *controller) Fatal(err error) {
	c.logg.Fatalf("FATAL %s", err)
}

func (c *controller) Error(err error) {
	c.logg.Fatalf("ERROR %s", err)
}

func (c *controller) Warning(err error) {
	c.logg.Fatalf("WARN %s", err)
}

func (c *controller) Info(err error) {
	c.logg.Fatalf("INFO %s", err)
}

func (c *controller) Debug(err error) {
	c.logg.Fatalf("DEBUG %s", err)
}
