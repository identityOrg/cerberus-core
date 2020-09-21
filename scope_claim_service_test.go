package core

import (
	"context"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScopeClaimStoreServiceImpl_Combined(t *testing.T) {
	scopeClaimStore := NewScopeClaimStoreServiceImpl(TestDb)
	ctx := context.Background()
	ctx = scopeClaimStore.BeginTransaction(ctx, true)
	var scopeId, claimId uint
	var err error
	var scope *models.ScopeModel
	t.Run("insert", func(t *testing.T) {
		scopeId, err = scopeClaimStore.CreateScope(ctx, "testscope", "A test scope")
		if assert.NoError(t, err) {
			assert.GreaterOrEqual(t, uint(1), scopeId)
		}
	})
	t.Run("find by id", func(t *testing.T) {
		scope, err = scopeClaimStore.GetScope(ctx, scopeId)
		if assert.NoError(t, err) {
			assert.Equal(t, scopeId, scope.ID)
			assert.Equal(t, "testscope", scope.Name)
		}
	})
	t.Run("create claim", func(t *testing.T) {
		claimId, err = scopeClaimStore.CreateClaim(ctx, "claim1", "claim one")
		if assert.NoError(t, err) {
			assert.GreaterOrEqual(t, uint(1), claimId)
		}
	})
	t.Run("add claim to scope", func(t *testing.T) {
		claimId2, err := scopeClaimStore.CreateClaim(ctx, "claim2", "claim two")
		if assert.NoError(t, err) {
			if assert.GreaterOrEqual(t, uint(1), claimId) {
				err = scopeClaimStore.AddClaimToScope(ctx, scopeId, claimId)
				assert.NoError(t, err)
				err = scopeClaimStore.AddClaimToScope(ctx, scopeId, claimId2)
				if assert.NoError(t, err) {
					err = scopeClaimStore.RemoveClaimFromScope(ctx, scopeId, claimId2)
					assert.NoError(t, err)
				}
			}
		}
	})
	scopeClaimStore.RollbackTransaction(ctx)
}
