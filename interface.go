package main

// generate events
type Generator interface {
	Name() string
	Generate()
}
