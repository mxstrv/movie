package main

type config struct {
	API apIConfig `yaml:"api"`
}

type apIConfig struct {
	Port int `yaml:"port"`
}
