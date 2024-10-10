package simple

import (
	"flag"
	"fmt"
)

type Config struct {
}

type simplePlugin struct {
	name  string
	value string
}

func NewSimplePlugin(name string) *simplePlugin {
	return &simplePlugin{name: name}
}

func (s *simplePlugin) GetPrefix() string {
	return s.name
}

func (s *simplePlugin) Get() interface{} {
	return s
}

func (s *simplePlugin) Name() string {
	return s.name
}

func (s *simplePlugin) InitFlags() {
	// Add flags for simple plugin
	flag.StringVar(&s.value, fmt.Sprintf("%s-value", s.name), "default value", "Some value of simple plugin")
}

func (s *simplePlugin) Configure() error {
	// Check if the value is empty
	if s.value == "" {
		return fmt.Errorf("simple plugin value is empty")
	}
	return nil
}

func (s *simplePlugin) Run() error {
	// Configure the service
	if err := s.Configure(); err != nil {
		return err
	}

	return nil
}

func (s *simplePlugin) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}

func (s *simplePlugin) GetValue() string {
	return s.value
}
