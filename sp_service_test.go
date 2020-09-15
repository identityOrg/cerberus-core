package core

import (
	"context"
	"github.com/identityOrg/cerberus-core/models"
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

func TestSPStoreServiceImpl_CreateSP(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb)
	ctx := spService.BeginTransaction(context.Background(), true)
	t.Run("create", func(t *testing.T) {
		metadata := &models.ServiceProviderMetadata{
			ApplicationType: "web",
		}
		spId, err := spService.CreateSP(ctx, "test create 1", "A SP created during test", metadata)
		if assert.NoError(t, err) {
			if assert.NotEqual(t, 0, spId) {
				sp, err := spService.GetSP(ctx, spId)
				if assert.NoError(t, err) {
					assert.Equal(t, "web", sp.Metadata.ApplicationType)
				}
			}
		}
	})
	spService.RollbackTransaction(ctx)
}

func TestSPStoreServiceImpl_ActivateSP(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb)
	ctx := spService.BeginTransaction(context.Background(), true)
	t.Run("activate", func(t *testing.T) {
		err := spService.ActivateSP(ctx, TestSP.ID)
		assert.NoError(t, err)
	})
	t.Run("deactivate", func(t *testing.T) {
		err := spService.DeactivateSP(ctx, TestSP.ID)
		assert.NoError(t, err)
	})
	spService.RollbackTransaction(ctx)
}

func TestSPStoreServiceImpl_DeleteSP(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb)
	ctx := spService.BeginTransaction(context.Background(), true)
	t.Run("existing", func(t *testing.T) {
		err := spService.DeleteSP(ctx, TestSP.ID)
		assert.NoError(t, err)
	})
	t.Run("non existing", func(t *testing.T) {
		err := spService.DeleteSP(ctx, TestSP.ID)
		assert.NoError(t, err)
	})
	spService.RollbackTransaction(ctx)
}

func TestSPStoreServiceImpl_FindSPByClientId(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb)
	ctx := spService.BeginTransaction(context.Background(), true)
	t.Run("existing", func(t *testing.T) {
		sp, err := spService.FindSPByClientId(ctx, TestSP.ClientID)
		if assert.NoError(t, err) {
			assert.Equal(t, TestSP.ID, sp.ID)
		}
	})
	t.Run("non existing", func(t *testing.T) {
		sp, err := spService.FindSPByClientId(ctx, TestSP.ClientID+"20000")
		assert.Error(t, err)
		assert.Nil(t, sp)
	})
	spService.RollbackTransaction(ctx)
}
