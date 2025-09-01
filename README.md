# waitfor-postgres

[![Build Status](https://github.com/go-waitfor/waitfor-postgres/actions/workflows/build.yml/badge.svg)](https://github.com/go-waitfor/waitfor-postgres/actions/workflows/build.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/go-waitfor/waitfor-postgres)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A PostgreSQL resource readiness assertion library built on top of the [waitfor](https://github.com/go-waitfor/waitfor) framework. This library allows you to wait for PostgreSQL databases to become available before proceeding with your application startup, making it ideal for containerized environments, integration tests, and service orchestration.

## Features

- **PostgreSQL connectivity testing**: Ping PostgreSQL databases to verify they're ready
- **Multiple URL schemes**: Supports both `postgres://` and `postgresql://` connection strings  
- **Configurable retry logic**: Customize attempts, intervals, and timeouts
- **Context support**: Full context cancellation and timeout support
- **Integration ready**: Built for Docker Compose, Kubernetes, and CI/CD pipelines

## Installation

```bash
go get github.com/go-waitfor/waitfor-postgres
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-waitfor/waitfor"
	"github.com/go-waitfor/waitfor-postgres"
)

func main() {
	runner := waitfor.New(postgres.Use())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := runner.Test(
		ctx,
		[]string{"postgres://user:password@localhost/mydb"},
		waitfor.WithAttempts(10),
		waitfor.WithInterval(2000), // 2000 milliseconds (2 seconds)
	)

	if err != nil {
		fmt.Printf("PostgreSQL not ready: %v\n", err)
		return
	}

	fmt.Println("PostgreSQL is ready!")
}
```

## Usage Examples

### Basic Connection Test

```go
runner := waitfor.New(postgres.Use())

err := runner.Test(
	context.Background(),
	[]string{"postgres://localhost/mydb"},
)
```

### Multiple Databases

```go
databases := []string{
	"postgres://user:pass@db1:5432/app_db",
	"postgresql://user:pass@db2:5432/cache_db",
}

err := runner.Test(context.Background(), databases)
```

### Custom Configuration

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

err := runner.Test(
	ctx,
	[]string{"postgres://user:pass@localhost/mydb"},
	waitfor.WithAttempts(20),           // Try up to 20 times
	waitfor.WithInterval(5000),         // Wait 5000 milliseconds (5 seconds) between attempts
)
```

### Error Handling

```go
err := runner.Test(ctx, []string{"postgres://localhost/mydb"})
if err != nil {
	// Handle different types of errors
	fmt.Printf("Database connection failed: %v\n", err)
	
	// Check if it's a timeout
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Timed out waiting for database")
	}
}
```

## Supported URL Formats

The library supports both PostgreSQL URL schemes:

- `postgres://user:password@host:port/dbname?sslmode=disable`
- `postgresql://user:password@host:port/dbname?sslmode=require`

### URL Components

```
postgresql://username:password@hostname:port/database?param1=value1&param2=value2
```

- **username**: Database username (optional)
- **password**: Database password (optional)  
- **hostname**: Database server hostname or IP
- **port**: Database server port (default: 5432)
- **database**: Database name (optional)
- **parameters**: Connection parameters like `sslmode`, `connect_timeout`, etc.

## Configuration Options

The library uses the waitfor framework's configuration options:

| Option | Description | Default |
|--------|-------------|---------|
| `WithAttempts(n)` | Maximum number of connection attempts | 30 |
| `WithInterval(ms)` | Time between connection attempts in milliseconds | 1000 (1 second) |

## Context and Timeouts

Always use context with timeouts for production applications:

```go
// Set overall timeout for the entire wait operation
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

// The Test method will respect the context deadline
err := runner.Test(ctx, []string{"postgres://localhost/mydb"}, waitfor.WithAttempts(20))

if errors.Is(err, context.DeadlineExceeded) {
	log.Fatal("Timed out waiting for PostgreSQL to become ready")
}
```

## Integration Examples

### Docker Compose

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password

  app:
    build: .
    depends_on:
      - postgres
    environment:
      DATABASE_URL: postgres://user:password@postgres:5432/myapp
```

### Kubernetes Init Container

```yaml
initContainers:
- name: wait-for-db
  image: my-app:latest
  command: ["/wait-for-postgres"]
  env:
  - name: DATABASE_URL
    value: "postgres://user:pass@postgres-service:5432/mydb"
```

## Troubleshooting

### Common Connection Issues

1. **Connection refused**: Check if PostgreSQL is running and accessible
2. **Authentication failed**: Verify username and password
3. **Database not found**: Ensure the database exists
4. **SSL/TLS issues**: Check `sslmode` parameter in connection string

### Debugging Tips

```go
// Enable detailed logging (if using a logger)
log.Printf("Testing connection to: %s", databaseURL)

err := runner.Test(ctx, []string{databaseURL})
if err != nil {
	log.Printf("Connection failed: %v", err)
}
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.