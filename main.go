package main

func main() {
	conf := InitConfigs()

	if conf.Payload != "" {
		if err := Publish(conf); err != nil {
			Logger.Errorf("Error while publishing %s", err)
		}
		return
	}
	if conf.Consume {
		if err := Consume(conf); err != nil {
			Logger.Errorf("Error while consuming %s", err)
		}
		return
	}
	if conf.Empty {
		if err := Empty(conf); err != nil {
			Logger.Errorf("Error while empty topic %s", err)
		}
		return
	}

	if conf.Delete {
		if err := Delete(conf); err != nil {
			Logger.Errorf("Error while empty topic %s", err)
		}
		return
	}
}
