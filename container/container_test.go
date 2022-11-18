package container_test

import (
	"context"
	"testing"

	"github.com/BrosSquad/vokapi/container"
	"github.com/stretchr/testify/require"
)

func TestNewContainer(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	di := container.New(context.Background(), &container.Config{
		BadgedDBPath: "",
	})

	assert.NotNil(di)
}

func TestGetBadgerDB(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})

	di := container.New(ctx, &container.Config{
		BadgedDBPath: t.TempDir(),
	})

	db := di.GetBadgerDB()

	assert.NotNil(db)

	otherDb := di.GetBadgerDB()

	assert.Equal(db, otherDb)
}
