# waitfor-postgres
Postgres resource readiness assertion library

# Quick start

```go
package main

import (
	"context"
	"fmt"
	"github.com/go-waitfor/waitfor"
	"github.com/go-waitfor/waitfor-postgres"
	"os"
)

func main() {
	runner := waitfor.New(postgres.Use())

	err := runner.Test(
		context.Background(),
		[]string{"postgres://localhost:8080/my-db"},
		waitfor.WithAttempts(5),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```