package core

import (
	"context"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSPStoreServiceImpl_FindAllSP(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb, nil, nil)
	spService.Db = beginTransaction(context.Background(), spService.Db)
	ctx := context.Background()
	t.Run("page 0", func(t *testing.T) {
		sp, count, err := spService.FindAllSP(ctx, 0, 5)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(sp))
			assert.Equal(t, uint(1), count)
		}
	})

	rollbackTransaction(spService.Db)
}

func TestSPStoreServiceImpl_CreateSP(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb, nil, nil)
	spService.Db = beginTransaction(context.Background(), spService.Db)
	ctx := context.Background()
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
	rollbackTransaction(spService.Db)
}

func TestSPStoreServiceImpl_ActivateSP(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb, nil, nil)
	spService.Db = beginTransaction(context.Background(), spService.Db)
	ctx := context.Background()
	t.Run("activate", func(t *testing.T) {
		err := spService.ActivateSP(ctx, TestSP.ID)
		assert.NoError(t, err)
	})
	t.Run("deactivate", func(t *testing.T) {
		err := spService.DeactivateSP(ctx, TestSP.ID)
		assert.NoError(t, err)
	})
	rollbackTransaction(spService.Db)
}

func TestSPStoreServiceImpl_DeleteSP(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb, nil, nil)
	spService.Db = beginTransaction(context.Background(), spService.Db)
	ctx := context.Background()
	t.Run("existing", func(t *testing.T) {
		err := spService.DeleteSP(ctx, TestSP.ID)
		assert.NoError(t, err)
	})
	t.Run("non existing", func(t *testing.T) {
		err := spService.DeleteSP(ctx, TestSP.ID)
		assert.NoError(t, err)
	})
	rollbackTransaction(spService.Db)
}

func TestSPStoreServiceImpl_FindSPByClientId(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb, nil, nil)
	spService.Db = beginTransaction(context.Background(), spService.Db)
	ctx := context.Background()
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
	rollbackTransaction(spService.Db)
}

func TestSPStoreServiceImpl_FindSPByName(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb, nil, nil)
	spService.Db = beginTransaction(context.Background(), spService.Db)
	ctx := context.Background()
	t.Run("existing", func(t *testing.T) {
		sp, err := spService.FindSPByName(ctx, TestSP.Name)
		if assert.NoError(t, err) {
			assert.Equal(t, TestSP.ID, sp.ID)
		}
	})
	t.Run("wild card", func(t *testing.T) {
		sp, err := spService.FindSPByName(ctx, "Test%")
		if assert.NoError(t, err) {
			assert.Equal(t, TestSP.ID, sp.ID)
		}
	})
	t.Run("non existing", func(t *testing.T) {
		sp, err := spService.FindSPByName(ctx, TestSP.Name+"20000")
		assert.Error(t, err)
		assert.Nil(t, sp)
	})
	rollbackTransaction(spService.Db)
}

func TestSPStoreServiceImpl_PatchSP(t *testing.T) {
	spService := NewSPStoreServiceImpl(TestDb, nil, nil)
	spService.Db = beginTransaction(context.Background(), spService.Db)
	ctx := context.Background()
	t.Run("patch existing", func(t *testing.T) {
		TestSP.Metadata.ApplicationType = "native"
		err := spService.PatchSP(ctx, TestSP.ID, TestSP.Metadata)
		if assert.NoError(t, err) {
			sp, err := spService.GetSP(ctx, TestSP.ID)
			if assert.NoError(t, err) {
				assert.Equal(t, "native", sp.Metadata.ApplicationType)
			}
		}
		TestSP.Metadata.ApplicationType = "web"
	})
	t.Run("non existing", func(t *testing.T) {
		TestSP.Metadata.ApplicationType = "native"
		err := spService.PatchSP(ctx, TestSP.ID+2000, TestSP.Metadata)
		assert.Error(t, err)
		TestSP.Metadata.ApplicationType = "web"
	})
	rollbackTransaction(spService.Db)
}

func TestSPStoreServiceImpl_ResetClientCredentials(t *testing.T) {
	encDec := &TextEncryptDecryptMock{}
	spService := NewSPStoreServiceImpl(TestDb, encDec, encDec)
	spService.Db = beginTransaction(context.Background(), spService.Db)
	ctx := context.Background()
	t.Run("reset existing", func(t *testing.T) {
		clientId, clientSecret, err := spService.ResetClientCredentials(ctx, TestSP.ID)
		if assert.NoError(t, err) {
			TestSP.ClientID = clientId
			TestSP.ClientSecret = clientSecret
			sp, err := spService.ValidateClientCredentials(ctx, clientId, clientSecret)
			if assert.NoError(t, err) {
				assert.Equal(t, TestSP.ID, sp)
			}
		}
	})
	t.Run("non existing", func(t *testing.T) {
		sp, err := spService.ValidateClientCredentials(ctx, TestSP.ClientID, TestSP.ClientSecret+"111")
		if assert.Error(t, err) {
			assert.Equal(t, uint(0), sp)
		}
	})
	rollbackTransaction(spService.Db)
}

type TextEncryptDecryptMock struct{}

func (m TextEncryptDecryptMock) DecryptText(_ context.Context, cypherText string) (text string, err error) {
	return cypherText, nil
}

func (TextEncryptDecryptMock) EncryptText(_ context.Context, text string) (cypherText string, err error) {
	return text, nil
}
