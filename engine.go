package main

//Engine - engine information
type Engine struct {
	Active  bool   `json:"active"`
	Locally string `json:"install"`
	Name    string `json:"name"`
	Config  string `json:"config"`
}
