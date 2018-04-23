package main

//Config - stucture of config file
type Config struct {
	Engines []Engine `json:"engines"`
	Ingores []Ingore `json:"ingores"`
}
