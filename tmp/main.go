package main

import (
	//	"fmt"

	//	"github.com/sinakeshmiri/goraz/packages/securitytrails"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/sinakeshmiri/goraz/packages/securitytrails"
	"github.com/sinakeshmiri/goraz/packages/shodan"
	"gopkg.in/yaml.v2"
)

type ipTables map[string]bool

type IPRange struct {
	Range string `yaml:"range"`
}

type Config struct {
	Host              string `yaml:"host"`
	SecuritytrailsKey string `yaml:"securitytrailsKey"`
	ShodanKey         string `yaml:"shodanKey"`
}

func main() {

	yamlFile, _ := filepath.Abs("./config")
	yamlData, err := ioutil.ReadFile(yamlFile)
	config := Config{}
	err = yaml.Unmarshal(yamlData, &config)

	if err != nil {
		fmt.Printf("error: %v", err)
	}
	ips, err := securitytrails.Find(config.Host, config.SecuritytrailsKey)
	if err != nil {
		fmt.Println(err)
	}
	ips = append(ips, shodan.Find(config.Host, config.ShodanKey)...)
	ips = securitytrails.RemoveDuplicateStr(ips)
	fmt.Println(ips)

}
