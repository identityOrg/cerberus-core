package core

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSecretStoreServiceImpl(t *testing.T) {
	secretService := NewSecretStoreServiceImpl(TestDb)
	ctx := secretService.BeginTransaction(context.Background(), true)
	var channelId uint
	var err error
	t.Run("create channel", func(t *testing.T) {
		channelId, err = secretService.CreateChannel(ctx, "channel1", "RS256", "sign", 10)
		if assert.NoError(t, err) {
			channel, err := secretService.GetChannel(ctx, channelId)
			if assert.NoError(t, err) {
				assert.Equal(t, "RS256", channel.Algorithm)
				assert.Equal(t, 1, len(channel.Secrets))
			}
		}
	})
	t.Run("renew secret", func(t *testing.T) {
		err = secretService.RenewSecret(ctx, channelId)
		if assert.NoError(t, err) {
			channel, err := secretService.GetChannel(ctx, channelId)
			if assert.NoError(t, err) {
				assert.Equal(t, "RS256", channel.Algorithm)
				assert.Equal(t, 2, len(channel.Secrets))
			}
		}
	})
	t.Run("jwks", func(t *testing.T) {
		secrets, err := secretService.GetAllSecrets(ctx)
		if assert.NoError(t, err) {
			if assert.NotNil(t, secrets) {
				assert.Equal(t, 2, len(secrets.Keys))
			}
		}
	})
	secretService.RollbackTransaction(ctx)
}
