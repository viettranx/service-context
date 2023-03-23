This is a simple cli application demo with ServiceContext.

## 1. Build the code
```shell
go build -o app
```

## 2. Output all ENV
```shell
./app outenv
```

You would see like this:
```
## Env for service. Ex: dev | stg | prd (-app-env)                                                                                                                                                INT ✘  11:09:23 
#APP_ENV="dev"

## gin mode (debug | release). Default debug (-gin-mode)
#GIN_MODE="debug"

## gin server port. Default 3000 (-gin-port)
#GIN_PORT=3000

## Log level: panic | fatal | error | warn | info | debug | trace (-log-level)
#LOG_LEVEL="trace"
```

These ENVs are generated from your components in ServiceContext.

### 3. Use .env file
```shell
./app outenv > .env
```

Open and edit `.env` file as you want. For example, I change `gin server port` to `3001`.
```
GIN_PORT=3001
```


### 4. Start the service (GIN HTTP)

```shell
./app
```

The result:

```
...
[GIN-debug] Listening and serving HTTP on :3001
```

