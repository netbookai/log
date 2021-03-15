# Log

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/go-coldbrew/log)

Package log provides a minimal interface for structured logging in services.
ColdBrew uses this log package for all logs.

## How To Use

The simplest way to use this package is by calling static log functions to report particular level (error/warning/info/debug)

```go
log.Error(...)
log.Warn(...)
log.Info(...)
log.Debug(...)
```

You can also initialize a new logger by calling 'log.NewLogger' and passing a loggers.BaseLogger implementation (loggers package provides a number of pre built implementations)

```go
logger := log.NewLogger(gokit.NewLogger())
logger.Info(ctx, "key", "value")
```

Note:

```go
Preferred logging output is in either logfmt or json format, so to facilitate these log function arguments should be in pairs of key-value
```

## Contextual Logs

log package uses context.Context to pass additional information to logs, you can use 'loggers.AddToLogContext' function to add additional information to logs. For example in access log from service

```go
{"@timestamp":"2018-07-30T09:58:18.262948679Z","caller":"http/http.go:66","error":null,"grpcMethod":"/AuthSvc.AuthService/Authenticate","level":"info","method":"POST","path":"/2.0/authenticate/","took":"1.356812ms","trace":"15592e1b-93df-11e8-bdfd-0242ac110002","transport":"http"}
```

we pass 'method', 'path', 'grpcMethod' and 'transport' from context, this information gets automatically added to all log calls called inside the service and makes debugging services much easier.
ColdBrew also generates a 'trace' ID per request, this can be used to trace an entire request path in logs.

this package is based on [https://github.com/carousell/Orion/tree/master/utils/log](https://github.com/carousell/Orion/tree/master/utils/log)

## Sub Packages

* [loggers](./loggers): Package loggers provides loggers implementation for log package

* [loggers/gokit](./loggers/gokit): Package gokit provides BaseLogger implementation for go-kit/log

* [loggers/logrus](./loggers/logrus): Package logrus provides a BaseLogger implementation for logrus

* [loggers/stdlog](./loggers/stdlog): Package stdlog provides a BaseLogger implementation for golang "log" package

* [wrap](./wrap): Package wrap provides multiple wrap functions to wrap log implementation of other log packages

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
