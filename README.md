<div align = "center">
<p>
    <img width="160" src="https://github.com/elliotxx/healthcheck/blob/main/logo.png?sanitize=true">
</p>
<h2>Low-dependency, High-efficiency, Kubernetes-style healthcheck for Gin Framework</h2>
<a title="Go Reference" target="_blank" href="https://pkg.go.dev/github.com/elliotxx/healthcheck"><img src="https://pkg.go.dev/badge/github.com/elliotxx/healthcheck.svg"></a>
<a title="Go Report Card" target="_blank" href="https://goreportcard.com/report/github.com/elliotxx/healthcheck"><img src="https://goreportcard.com/badge/github.com/elliotxx/healthcheck?style=flat-square"></a>
<a title="Coverage Status" target="_blank" href="https://coveralls.io/github/elliotxx/healthcheck?branch=main"><img src="https://img.shields.io/coveralls/github/elliotxx/healthcheck/main"></a>
<a title="Code Size" target="_blank" href="https://github.com/elliotxx/healthcheck"><img src="https://img.shields.io/github/languages/code-size/elliotxx/healthcheck.svg?style=flat-square"></a>
<br>
<a title="GitHub release" target="_blank" href="https://github.com/elliotxx/healthcheck/releases"><img src="https://img.shields.io/github/release/elliotxx/healthcheck.svg"></a>
<a title="License" target="_blank" href="https://github.com/elliotxx/healthcheck/blob/main/LICENSE"><img src="https://img.shields.io/github/license/elliotxx/healthcheck"></a>
</p>
</div>


This module will create a [**kubernetes-style** endpoints](https://kubernetes.io/docs/reference/using-api/health-checks/) for Gin framework. Inspired by [tavsec/gin-healthcheck](https://github.com/tavsec/gin-healthcheck).

## 📜 Language

[English](https://github.com/elliotxx/healthcheck/blob/main/README.md) | [简体中文](https://github.com/elliotxx/healthcheck/blob/main/README-zh.md)


## ✨ Core Features
* ⚡ Lightweight
* 🌲 Low dependency
* 🔥 High efficiency
* 🔨 Highly customizable
* ⎈  Kubernetes-style


## ⚙️ Usage
```shell
go get github.com/elliotxx/healthcheck
```


## 📖 Examples
### Default Check
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/elliotxx/healthcheck"
    "github.com/elliotxx/healthcheck/checks"
)

func main() {
    r := gin.Default()

    healthcheck.Register(r)
	
    r.Run()
}
```

This will register the default healthcheck endpoint (`/healthz`) to the route. The path can be customized
using `healthcheck.Config` structure.

Or use `NewHandler()` function directly:
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/elliotxx/healthcheck"
    "github.com/elliotxx/healthcheck/checks"
)

func main() {
    r := gin.Default()

    r.GET("livez", NewHandler(NewDefaultHandlerConfig()))

    readyzChecks := []checks.Check{checks.NewPingCheck(), checks.NewEnvCheck("DB_HOST")}
    r.GET("readyz",NewHandler(NewDefaultHandlerConfigFor(readyzChecks)))
	
    r.Run()
}
```

Enjoy it!

```shell
$ curl -k http://localhost/readyz
OK

$ curl -k http://localhost/readyz?verbose
[+] Ping ok
[-] Env-DB_HOST ok
health check failed

$ curl -k http://localhost/readyz?verbose&excludes=Env-DB_HOST
[+] Ping ok
health check passed
```
