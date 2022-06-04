package config

func (c *Config) UpdateConfig(_newConfig Config) {
	oldConfig := *c

	//move private fields over
	_newConfig.appCallback = c.appCallback
	_newConfig.altConfigLocation = c.altConfigLocation
	*c = _newConfig

	c.appCallback(oldConfig, *c)
	c.Save()
}
