package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// CnxConfig is the struct representing the structure of the environement yaml file describing how to connect
// to database and some other paramters
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type AppConfig struct {
	Catsdb DBConfig `yaml:"db"`
}

func readYamlCnxFile(filename string) (AppConfig, error) {
	var config AppConfig

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		return config, err
	}
	return config, err
}
