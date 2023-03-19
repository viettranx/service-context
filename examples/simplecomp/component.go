package main

import (
	"flag"
	sctx "github.com/viettranx/service-context"
)

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

func (s *simpleComponent) Activate(_ sctx.ServiceContext) error {
	return nil
}

func (s *simpleComponent) Stop() error {
	return nil
}

func (s *simpleComponent) GetValue() string {
	return s.value
}
