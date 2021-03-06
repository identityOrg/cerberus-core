package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/identityOrg/oidcsdk"
	"gorm.io/gorm"
)

type SPStoreServiceImpl struct {
	Db      *gorm.DB
	TextEnc ITextEncrypts
	TextDec ITextDecrypts
}

func (s *SPStoreServiceImpl) GetClient(ctx context.Context, clientID string) (client oidcsdk.IClient, err error) {
	return s.FindSPByClientId(ctx, clientID)
}

func (s *SPStoreServiceImpl) FetchClientProfile(_ context.Context, clientID string) oidcsdk.RequestProfile {
	profile := oidcsdk.RequestProfile{}
	profile.SetUsername(clientID)
	return profile
}

func (s *SPStoreServiceImpl) CreateSP(ctx context.Context, clientName string, description string, metadata *models.ServiceProviderMetadata) (id uint, err error) {
	user := &models.ServiceProviderModel{
		Name:        clientName,
		Description: description,
		Metadata:    metadata,
		Public:      true,
		Active:      true,
	}
	db := s.Db.WithContext(ctx)
	saveResult := db.Save(user)
	return user.ID, saveResult.Error
}

func (s *SPStoreServiceImpl) UpdateSP(ctx context.Context, id uint, public bool, metadata *models.ServiceProviderMetadata) (err error) {
	user := &models.ServiceProviderModel{}
	user.ID = id
	db := s.Db.WithContext(ctx)
	findResult := db.Find(user)
	if findResult.Error != nil {
		return findResult.Error
	}
	if findResult.RowsAffected != 1 {
		return fmt.Errorf("client not found with id %d", id)
	}
	user.Metadata = metadata
	user.Public = public
	return db.Save(user).Error
}

func (s *SPStoreServiceImpl) PatchSP(ctx context.Context, id uint, metadata *models.ServiceProviderMetadata) (err error) {
	user := &models.ServiceProviderModel{}
	user.ID = id
	db := s.Db.WithContext(ctx)
	findResult := db.Find(user)
	if findResult.Error != nil {
		return findResult.Error
	}
	if findResult.RowsAffected != 1 {
		return fmt.Errorf("client not found with id %d", id)
	}
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, user.Metadata)
	if err != nil {
		return err
	}
	return db.Save(user).Error
}

func (s *SPStoreServiceImpl) DeleteSP(ctx context.Context, id uint) (err error) {
	user := &models.ServiceProviderModel{}
	user.ID = id
	db := s.Db.WithContext(ctx)
	return db.Delete(user).Error
}

func (s *SPStoreServiceImpl) ActivateSP(ctx context.Context, id uint) error {
	return s.updateStatus(ctx, id, true)
}

func (s *SPStoreServiceImpl) DeactivateSP(ctx context.Context, id uint) error {
	return s.updateStatus(ctx, id, false)
}

func (s *SPStoreServiceImpl) updateStatus(ctx context.Context, id uint, active bool) error {
	user := &models.ServiceProviderModel{}
	user.ID = id
	db := s.Db.WithContext(ctx)
	updateResult := db.Model(user).Update("active", active)
	if updateResult.Error != nil {
		return updateResult.Error
	} else if updateResult.RowsAffected != 1 {
		return fmt.Errorf("no SP found with id %d", id)
	}
	return nil
}

func (s *SPStoreServiceImpl) ResetClientCredentials(ctx context.Context, id uint) (clientId, clientSecret string, err error) {
	sp, err := s.GetSP(ctx, id)
	if err != nil {
		return "", "", err
	}
	if sp.Public {
		return "", "", fmt.Errorf("service provider not private")
	}
	encrypted, err := s.TextEnc.EncryptText(ctx, uuid.New().String())
	if err != nil {
		return "", "", err
	}
	sp.ClientSecret = encrypted
	sp.ClientID = uuid.New().String()
	db := s.Db.WithContext(ctx)
	result := db.Model(&models.ServiceProviderModel{}).
		Where("id = ?", id).
		UpdateColumns(models.ServiceProviderModel{
			ClientID:     sp.ClientID,
			ClientSecret: sp.ClientSecret,
		})
	if result.Error != nil {
		return "", "", result.Error
	}
	return sp.ClientID, sp.ClientSecret, nil
}

func (s *SPStoreServiceImpl) ValidateClientCredentials(ctx context.Context, clientId, clientSecret string) (id uint, err error) {
	sp, err := s.FindSPByClientId(ctx, clientId)
	if err != nil {
		return 0, err
	}
	if !sp.Active {
		return 0, fmt.Errorf("service provider is inactive")
	}
	if sp.Public {
		return sp.ID, nil
	}
	decryptedSecret, err := s.TextDec.DecryptText(ctx, sp.ClientSecret)
	if err != nil {
		return 0, fmt.Errorf("failed to decrypt sp secret - %v", err)
	}

	if clientSecret == decryptedSecret {
		return sp.ID, nil
	}
	return 0, fmt.Errorf("invalid client secret")
}

func (s *SPStoreServiceImpl) ValidateSecretSignature(ctx context.Context, token string) (id uint, err error) {
	panic("implement me")
}

func (s *SPStoreServiceImpl) ValidatePrivateKeySignature(ctx context.Context, token string) (id uint, err error) {
	panic("implement me")
}

func (s *SPStoreServiceImpl) GetSP(ctx context.Context, id uint) (sp *models.ServiceProviderModel, err error) {
	tx := s.Db.WithContext(ctx)
	sp = &models.ServiceProviderModel{}
	result := tx.Find(sp, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("no SP found with id %d", id)
	}
	return
}

func (s *SPStoreServiceImpl) FindSPByClientId(ctx context.Context, clientId string) (sp *models.ServiceProviderModel, err error) {
	tx := s.Db.WithContext(ctx)
	sp = &models.ServiceProviderModel{}
	result := tx.Find(sp, "client_id = ?", clientId)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("sp with client_id = %s not found", clientId)
	}
	return
}

func (s *SPStoreServiceImpl) FindSPByName(ctx context.Context, name string) (sp *models.ServiceProviderModel, err error) {
	tx := s.Db.WithContext(ctx)
	sp = &models.ServiceProviderModel{}
	result := tx.First(sp, "name like ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("no client found with name like %s", name)
	}
	err = result.Error
	return
}

func (s *SPStoreServiceImpl) FindAllSP(ctx context.Context, page uint, pageSize uint) (sps []models.ServiceProviderModel, count uint, err error) {
	var total int64
	tx := s.Db.WithContext(ctx)
	query := tx.Select([]string{"id", "name", "description", "client_id", "active"}).Model(&models.ServiceProviderModel{})
	err = query.Limit(int(pageSize)).Offset(int(pageSize * page)).Find(&sps).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return sps, uint(total), nil
}

func NewSPStoreServiceImpl(db *gorm.DB, dec ITextDecrypts, enc ITextEncrypts) *SPStoreServiceImpl {
	return &SPStoreServiceImpl{Db: db, TextEnc: enc, TextDec: dec}
}
