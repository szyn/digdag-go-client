# digdag-go-client

[![codecov](https://codecov.io/gh/szyn/digdag-go-client/branch/master/graph/badge.svg?token=XCc0ySTin2)](https://codecov.io/gh/szyn/digdag-go-client)
[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/szyn/digdag-go-client)

## Description

Digdag client library written in Go. (Unofficial)

## Requirement

* digdag server

## Install

```console
$ go get -u github.com/szyn/digdag-go-client
```

## Usage

### e.g. Get projects infomation

```go
package main

import (
	"fmt"

	digdag "github.com/szyn/digdag-go-client"
)

func main() {
	/*
		default endpoint: http://localhost:65432
		If you want to change endpoint, write as the following.
		client, err := digdag.NewClient("http://hostname:5432", false)
	*/
	client, err := digdag.NewClient("", false)
	if err != nil {
		fmt.Println(err)
	}

	projects, err := client.GetProjects()
	if err != nil {
		fmt.Println(err)
	}
	for _, project := range projects {
		fmt.Println(project.ID, project.Name)
	}
}

```

See also: [Godoc](https://godoc.org/github.com/szyn/digdag-go-client)

## Contribution

1. Fork it ( http://github.com/szyn/digdag-go-client )
1. Create your feature branch ( `git checkout -b my-new-feature` )
1. Commit your changes ( `git commit -am 'Add some feature'` )
1. Push to the branch ( `git push origin my-new-feature` )
1. Run test suite ( `go test ./...` ) and confirm that it passes
1. Run gofmt ( `gofmt -s` )
1. Create new Pull Request 😆

## LICENCE

[Apache License 2.0](LICENSE)

## Author

[szyn](https://github.com/szyn)