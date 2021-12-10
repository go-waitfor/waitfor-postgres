package postgres_test

import (
	"context"
	"github.com/go-waitfor/waitfor"
	"github.com/go-waitfor/waitfor-postgres"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUse(t *testing.T) {
	w := waitfor.New(postgres.Use())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := w.Test(ctx, []string{"postgres://usr:pass@localhost/my-db"})

	assert.Error(t, err)
}
