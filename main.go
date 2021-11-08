package main

import (
	"github.com/ilhamabdlh/go-restapi/configes"
)

type Configs interface{
	getConfigs()
	getConfig()
	createConfig()
	updateConfigs()
	deleteConfigs()
	Configs()
}

func excecuteConfigs(c Config){
	c.getConfigs()
	c.getConfig()
	c.createConfig()
	c.updateConfigs()
	c.deleteConfigs()
	c.Configs()
}
func main() {
	excecuteConfigs()
}
