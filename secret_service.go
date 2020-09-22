package core

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/google/uuid"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/square/go-jose.v2"
	"time"
)

type SecretStoreServiceImpl struct {
	Db *gorm.DB
}

func NewSecretStoreServiceImpl(db *gorm.DB) *SecretStoreServiceImpl {
	return &SecretStoreServiceImpl{Db: db}
}

func (s *SecretStoreServiceImpl) BeginTransaction(ctx context.Context, readOnly bool) context.Context {
	return beginTransaction(ctx, readOnly, s.Db)
}

func (s *SecretStoreServiceImpl) CommitTransaction(ctx context.Context) context.Context {
	return commitTransaction(ctx)
}

func (s *SecretStoreServiceImpl) RollbackTransaction(ctx context.Context) context.Context {
	return rollbackTransaction(ctx)
}

func (s *SecretStoreServiceImpl) GetAllSecrets(ctx context.Context) (*jose.JSONWebKeySet, error) {
	db := getTransaction(ctx)
	secrets := make([]models.SecretModel, 0)
	db.Find(&secrets)
	keySet := &jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 0),
	}
	for _, secret := range secrets {
		key, err := x509.ParsePKCS8PrivateKey(secret.Value)
		if err != nil {
			continue
		}
		jwk := jose.JSONWebKey{
			Key:       key,
			KeyID:     secret.KeyId,
			Algorithm: secret.Algorithm,
			Use:       secret.Use,
		}
		keySet.Keys = append(keySet.Keys, jwk)
	}
	return keySet, nil
}

func (s *SecretStoreServiceImpl) CreateChannel(ctx context.Context, name string, algorithm string, use string, validityDay uint) (uint, error) {
	channel := &models.SecretChannelModel{
		Name:        name,
		Algorithm:   algorithm,
		Use:         use,
		ValidityDay: validityDay,
	}
	secret := &models.SecretModel{
		IssuedAt:  time.Now(),
		KeyId:     uuid.New().String(),
		Algorithm: algorithm,
		Use:       use,
	}
	var err error
	secret.Value, err = s.createSecret(algorithm)
	if err != nil {
		return 0, err
	}
	validityHour := time.Duration(validityDay) * time.Duration(24)
	secret.ExpiresAt = time.Now().Add(validityHour * time.Hour)
	channel.Secrets = append(channel.Secrets, secret)

	db := getTransaction(ctx)
	result := db.Save(channel)
	return channel.ID, result.Error
}

func (s *SecretStoreServiceImpl) GetAllChannels(ctx context.Context) ([]*models.SecretChannelModel, error) {
	db := getTransaction(ctx)
	channels := make([]*models.SecretChannelModel, 0)
	findResult := db.Find(channels)
	return channels, findResult.Error
}

func (s *SecretStoreServiceImpl) GetChannel(ctx context.Context, channelId uint) (*models.SecretChannelModel, error) {
	db := getTransaction(ctx)
	channels := &models.SecretChannelModel{}
	findResult := db.Preload("Secrets").Find(channels, channelId)
	return channels, findResult.Error
}

func (s *SecretStoreServiceImpl) DeleteChannel(ctx context.Context, channelId uint) error {
	db := getTransaction(ctx)
	return db.Delete(&models.SecretChannelModel{}, channelId).Error
}

func (s *SecretStoreServiceImpl) RenewSecret(ctx context.Context, channelId uint) error {
	db := getTransaction(ctx)
	channel := &models.SecretChannelModel{}
	channel.ID = channelId
	channelResult := db.Preload("Secrets").Find(channel)
	if channelResult.Error != nil {
		return channelResult.Error
	}
	currentTime := time.Now()
	expiry := time.Duration(channel.ValidityDay) * time.Duration(24) * time.Hour
	for _, secret := range channel.Secrets {
		if secret.ExpiresAt.After(currentTime) {
			secret.ExpiresAt = currentTime
			replace := db.Save(secret)
			if replace.Error != nil {
				return replace.Error
			}
		}
	}
	newSecret := &models.SecretModel{
		KeyId:     uuid.New().String(),
		IssuedAt:  currentTime,
		ExpiresAt: currentTime.Add(expiry),
		ChannelId: channelId,
		Algorithm: channel.Algorithm,
		Use:       channel.Use,
	}
	var err error
	newSecret.Value, err = s.createSecret(channel.Algorithm)
	if err != nil {
		return err
	}
	return db.Save(newSecret).Error
}

func (s *SecretStoreServiceImpl) createSecret(algorithm string) ([]byte, error) {
	var key interface{}
	var err error

	switch algorithm {
	case string(jose.RS256):
		key, err = rsa.GenerateKey(rand.Reader, 1024)
	case string(jose.RS384):
		key, err = rsa.GenerateKey(rand.Reader, 2048)
	case string(jose.RS512):
		key, err = rsa.GenerateKey(rand.Reader, 4096)
	case string(jose.PS256):
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case string(jose.PS384):
		key, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case string(jose.PS512):
		key, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return nil, fmt.Errorf("algorithm %s is not supported", algorithm)
	}
	if err != nil {
		return nil, err
	}
	var data []byte
	data, err = x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}
