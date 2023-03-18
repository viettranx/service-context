package main

import (
	"github.com/viettranx/service-context"
	"log"
)

type CanGetValue interface {
	GetValue() string
}

func main() {
	const compId = "foo"

	serviceCtx := sctx.NewServiceContext(
		sctx.WithName("Simple Component"),
		sctx.WithComponent(NewSimpleComponent(compId)),
	)

	if err := serviceCtx.Load(); err != nil {
		log.Fatal(err)
	}

	comp := serviceCtx.MustGet(compId).(CanGetValue)

	log.Println(comp.GetValue())

	_ = serviceCtx.Stop()
}
