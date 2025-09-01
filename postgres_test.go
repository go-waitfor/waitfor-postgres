package postgres_test

import (
	"context"
	"github.com/go-waitfor/waitfor"
	"github.com/go-waitfor/waitfor-postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
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

func TestNew(t *testing.T) {
	t.Run("with valid URL", func(t *testing.T) {
		u, err := url.Parse("postgres://localhost:5432/testdb")
		require.NoError(t, err)

		resource, err := postgres.New(u)
		assert.NoError(t, err)
		assert.NotNil(t, resource)
	})

	t.Run("with nil URL", func(t *testing.T) {
		resource, err := postgres.New(nil)
		assert.Error(t, err)
		assert.Nil(t, resource)
		assert.Contains(t, err.Error(), "url")
		assert.Contains(t, err.Error(), "invalid argument")
	})
}

func TestPostgres_Test(t *testing.T) {
	t.Run("with invalid postgres URL", func(t *testing.T) {
		u, err := url.Parse("postgres://nonexistent:5432/testdb")
		require.NoError(t, err)

		resource, err := postgres.New(u)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err = resource.Test(ctx)
		assert.Error(t, err)
	})

	t.Run("with URL that causes connection failure", func(t *testing.T) {
		// Use a valid URL format but with invalid host
		u, err := url.Parse("postgres://user:password@invalid-host:5432/database")
		require.NoError(t, err)

		resource, err := postgres.New(u)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err = resource.Test(ctx)
		assert.Error(t, err)
	})

	t.Run("with postgresql scheme (alternate scheme bug test)", func(t *testing.T) {
		// This tests a potential bug where using 'postgresql://' scheme
		// instead of 'postgres://' could cause sql.Open to fail
		u, err := url.Parse("postgresql://user:password@localhost:5432/database")
		require.NoError(t, err)

		resource, err := postgres.New(u)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// This should work (or at least fail with connection error, not driver error)
		err = resource.Test(ctx)
		// We expect a connection error, not a driver registration error
		assert.Error(t, err)
		// The error should not be about unknown driver
		assert.NotContains(t, err.Error(), "unknown driver")
	})
}
