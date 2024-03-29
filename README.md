# Yandex Appmetrica Push API Go Library

Golang client library for Push API. (WIP warning)

For additional information see [Appmetrica Push API documentation](https://appmetrica.yandex.ru/docs/mobile-api/push/about.html)
## Getting Started
To install appmetrica-push-go, use `go get`:

```bash
go get github.com/Fodro/appmetrica-push-go
```
## Sample Usage
You can directly access structs and assemble them by yourself...
```go
package main

import (
	"fmt"
	appmetrica "github.com/Fodro/appmetrica-push-go"
)

func main() {
	client := appmetrica.NewClient("token")
	group, err := client.CreateGroup(&appmetrica.Group{
		AppId:    12345,
		Name:     "name",
		SendRate: 100500,
	})
	
	if err != nil {
		panic(err)
    }
	
	fmt.Println(fmt.Sprintf("%+v\n", group))
}

```
...or use built-in functions to construct minimal structs and modify them by accessing their attributes
```go
package main

import (
	"fmt"
	appmetrica "github.com/Fodro/appmetrica-push-go"
)

func main() {
	client := appmetrica.NewClient("token")
	g := appmetrica.NewCreateGroupRequest(1234, "name")
	g.SendRate = 5000
	group, err := client.CreateGroup(g)

	if err != nil {
		panic(err)
	}
	
	fmt.Println(fmt.Sprintf("%+v\n", group))
}

```
## Plans
* More comfortable error handling
* Extend functionality to all Appmetrica API