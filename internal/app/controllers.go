package app

func (a *app) RegisterController(name string, c interface{}) {
	for cName, _ := range a.controllers {
		if cName == name {
			a.Fatal("controller already exists with name:", name)
		}
	}
	a.controllers[name] = c
}

func (a *app) Controller(name string) interface{} {
	for cName, c := range a.controllers {
		if cName == name {
			return c
		}
	}
	a.Fatal("invalid controller name")
	return nil
}
