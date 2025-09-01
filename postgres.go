package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/lib/pq"

	"github.com/go-waitfor/waitfor"
)

const Scheme = "postgres"

type Postgres struct {
	url *url.URL
}

func Use() waitfor.ResourceConfig {
	return waitfor.ResourceConfig{
		Scheme:  []string{Scheme},
		Factory: New,
	}
}

func New(u *url.URL) (waitfor.Resource, error) {
	if u == nil {
		return nil, fmt.Errorf("%q: %w", "url", waitfor.ErrInvalidArgument)
	}

	return &Postgres{u}, nil
}

func (s *Postgres) Test(ctx context.Context) error {
	// Always use "postgres" as the driver name for sql.Open, regardless of the URL scheme
	// Remove the scheme and "://" from the URL to get the connection string
	connStr := strings.TrimPrefix(s.url.String(), s.url.Scheme+"://")

	db, err := sql.Open(Scheme, connStr)

	if err != nil {
		return err
	}

	defer db.Close()

	return db.PingContext(ctx)
}
