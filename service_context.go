// Copyright (c) 2023, Viet Tran, 200Lab Team.

package sctx

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	DevEnv = "dev"
	StgEnv = "stg"
	PrdEnv = "prd"
)

type Component interface {
	ID() string
	InitFlags()
	Activate(ServiceContext) error
	Stop() error
}

type ServiceContext interface {
	Load() error
	MustGet(id string) interface{}
	Get(id string) (interface{}, bool)
	Logger(prefix string) Logger
	EnvName() string
	GetName() string
	Stop() error
	OutEnv()
}

type serviceCtx struct {
	name       string
	env        string
	components []Component
	store      map[string]Component
	cmdLine    *AppFlagSet
	logger     Logger
}

func NewServiceContext(opts ...Option) ServiceContext {
	sv := &serviceCtx{
		store: make(map[string]Component),
	}

	sv.components = []Component{defaultLogger}

	for _, opt := range opts {
		opt(sv)
	}

	sv.initFlags()

	sv.cmdLine = newFlagSet(sv.name, flag.CommandLine)
	sv.parseFlags()

	sv.logger = defaultLogger.GetLogger(sv.name)

	return sv
}

func (s *serviceCtx) initFlags() {
	flag.StringVar(&s.env, "app-env", DevEnv, "Env for service. Ex: dev | stg | prd")

	for _, c := range s.components {
		c.InitFlags()
	}
}

func (s *serviceCtx) Get(id string) (interface{}, bool) {
	c, ok := s.store[id]

	if !ok {
		return nil, false
	}

	return c, true
}

func (s *serviceCtx) MustGet(id string) interface{} {
	c, ok := s.Get(id)

	if !ok {
		panic(fmt.Sprintf("can not get %s\n", id))
	}

	return c
}

func (s *serviceCtx) Load() error {
	s.logger.Infoln("Service context is loading...")

	for _, c := range s.components {
		if err := c.Activate(s); err != nil {
			return err
		}
	}

	return nil
}

func (s *serviceCtx) Logger(prefix string) Logger {
	return defaultLogger.GetLogger(prefix)
}

func (s *serviceCtx) Stop() error {
	s.logger.Infoln("Stopping service context")
	for i := range s.components {
		if err := s.components[i].Stop(); err != nil {
			return err
		}
	}

	s.logger.Infoln("Service context stopped")

	return nil
}

func (s *serviceCtx) GetName() string { return s.name }
func (s *serviceCtx) EnvName() string { return s.env }
func (s *serviceCtx) OutEnv()         { s.cmdLine.GetSampleEnvs() }

type Option func(*serviceCtx)

func WithName(name string) Option {
	return func(s *serviceCtx) { s.name = name }
}

func WithComponent(c Component) Option {
	return func(s *serviceCtx) {
		if _, ok := s.store[c.ID()]; !ok {
			s.components = append(s.components, c)
			s.store[c.ID()] = c
		}
	}
}

func (s *serviceCtx) parseFlags() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if err == nil {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatalf("Loading env(%s): %s", envFile, err.Error())
		}
	} else if envFile != ".env" {
		log.Fatalf("Loading env(%s): %s", envFile, err.Error())
	}

	s.cmdLine.Parse([]string{})
}
