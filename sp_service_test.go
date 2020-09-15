package core

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSPStoreServiceImpl_FindAllSP(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb)
	ctx := spService.BeginTransaction(context.Background(), true)
	t.Run("page 0", func(t *testing.T) {
		sp, count, err := spService.FindAllSP(ctx, 0, 5)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(sp))
			assert.Equal(t, uint(1), count)
		}
	})

	spService.RollbackTransaction(ctx)
}
