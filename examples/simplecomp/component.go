package main

import "flag"

type simpleComponent struct {
	id    string
	value string
}

func NewSimpleComponent(id string) *simpleComponent {
	return &simpleComponent{id: id}
}

func (s *simpleComponent) ID() string {
	return s.id
}

func (s *simpleComponent) InitFlags() {
	flag.StringVar(&s.value, "simple-value", "demo", "Value in string")
}

func (s *simpleComponent) Activate() error {
	return nil
}

func (s *simpleComponent) Stop() error {
	return nil
}

func (s *simpleComponent) GetValue() string {
	return s.value
}
