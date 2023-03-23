# Service-Context

Service Context runs as heart of our services, helps us manage components (such as Consumers, DB Connections, Configuration).

It offers:

- Logger component (using [logrus](https://github.com/sirupsen/logrus)).
- Manage ENV and Flag variables dynamically (with [flagenv](github.com/facebookgo/flagenv), [godotenv](github.com/joho/godotenv)).
- Output env and flag with .env formatted. 
- Easily plug-and-play component (as plugin).

## Examples
- [Demo with simple GIN-HTTP](./examples/ginhttp)
- [Custom component](./examples/simplecomp)
- [Simple CLI](./examples/simplecli) (run service and output all env command)

## How to use it

### 1. Install it with this cmd:
```shell
go get -u github.com/viettranx/service-context
```

### 2. Define your component:
Component can be anything but implements this interface:

```go
type Component interface {
	ID() string
	InitFlags()
	Activate(ServiceContext) error
	Stop() error
}
```

Simple demo custom component:

```go
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
```

### 3. Use the component with Service-Context:

```go
package main

import (
	"github.com/viettranx/service-context"
	"log"
)

func main() {
	const compId = "foo" // identity your component

	// Init service-context, you can put components as much as you can
	serviceCtx := sctx.NewServiceContext(
		sctx.WithComponent(NewSimpleComponent(compId)),
	)

	// Load() will iterate registered components
	// It does parse flags and make some configurations if you defined
	// in Activate() method of the components
	if err := serviceCtx.Load(); err != nil {
		log.Fatal(err)
	}

	type CanGetValue interface {
		GetValue() string
	}

	// Get the component from ServiceContext by its ID
	comp := serviceCtx.MustGet(compId).(CanGetValue)

	log.Println(comp.GetValue())

	_ = serviceCtx.Stop() // can be omitted
}
```

### 4. Run your code with ENV

```shell
go build -o app
SIMPLE_VALUE="Hello World" ./app
```

You will see `Hello World` on your console.

### Why I should use the ServiceContext?

When your service growths, a lot of components are required. More components more configurations and more ENVs to manage. That's why ServiceContext plays its role.

> Service-Context also offers a helpful feature: output all ENV with dot-env (.env) formatted with: `serviceContext.Outenv()`.
> Please check [this example](./examples/simplecli).

Hope it helps and enjoy coding!
