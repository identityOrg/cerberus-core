package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/jinzhu/gorm"
)

type SPStoreServiceImpl struct {
	Db      *gorm.DB
	TextEnc ITextEncrypts
	TextDec ITextDecrypts
}

func (S *SPStoreServiceImpl) BeginTransaction(ctx context.Context, readOnly bool) context.Context {
	return beginTransaction(ctx, readOnly, S.Db)
}

func (S *SPStoreServiceImpl) CommitTransaction(ctx context.Context) context.Context {
	return commitTransaction(ctx)
}

func (S *SPStoreServiceImpl) RollbackTransaction(ctx context.Context) context.Context {
	return rollbackTransaction(ctx)
}

func (S *SPStoreServiceImpl) CreateSP(ctx context.Context, clientName string, description string, metadata *models.ServiceProviderMetadata) (id uint, err error) {
	user := &models.ServiceProviderModel{
		Name:        clientName,
		Description: description,
		Metadata:    metadata,
		Public:      true,
		Active:      true,
	}
	db := getTransaction(ctx)
	saveResult := db.Save(user)
	return user.ID, saveResult.Error
}

func (S *SPStoreServiceImpl) UpdateSP(ctx context.Context, id uint, metadata *models.ServiceProviderMetadata) (err error) {
	user := &models.ServiceProviderModel{}
	user.ID = id
	db := getTransaction(ctx)
	findResult := db.Find(user)
	if findResult.RecordNotFound() {
		return errors.New("service provider not found")
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	user.Metadata = metadata
	return db.Save(user).Error
}

func (S *SPStoreServiceImpl) PatchSP(ctx context.Context, id uint, metadata *models.ServiceProviderMetadata) (err error) {
	user := &models.ServiceProviderModel{}
	user.ID = id
	db := getTransaction(ctx)
	findResult := db.Find(user)
	if findResult.RecordNotFound() {
		return errors.New("service provider not found")
	}
	if findResult.Error != nil {
		return findResult.Error
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

func (S *SPStoreServiceImpl) DeleteSP(ctx context.Context, id uint) (err error) {
	user := &models.ServiceProviderModel{}
	user.ID = id
	db := getTransaction(ctx)
	return db.Delete(user).Error
}

func (S *SPStoreServiceImpl) ActivateSP(ctx context.Context, id uint) error {
	return S.updateStatus(ctx, id, true)
}

func (S *SPStoreServiceImpl) DeactivateSP(ctx context.Context, id uint) error {
	return S.updateStatus(ctx, id, false)
}

func (S *SPStoreServiceImpl) updateStatus(ctx context.Context, id uint, active bool) error {
	user := &models.ServiceProviderModel{}
	user.ID = id
	db := getTransaction(ctx)
	updateResult := db.Model(user).Update("active", active)
	if updateResult.Error != nil {
		return updateResult.Error
	} else if updateResult.RowsAffected != 1 {
		return fmt.Errorf("no SP found with id %d", id)
	}
	return nil
}

func (S *SPStoreServiceImpl) ResetClientCredentials(ctx context.Context, id uint) (clientId, clientSecret string, err error) {
	sp, err := S.GetSP(ctx, id)
	if err != nil {
		return "", "", err
	}
	if sp.Public {
		return "", "", fmt.Errorf("service provider not private")
	}
	encrypted, err := S.TextEnc.EncryptText(ctx, uuid.New().String())
	if err != nil {
		return "", "", err
	}
	sp.ClientSecret = encrypted
	sp.ClientID = uuid.New().String()
	db := getTransaction(ctx)
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

func (S *SPStoreServiceImpl) ValidateClientCredentials(ctx context.Context, clientId, clientSecret string) (id uint, err error) {
	sp, err := S.FindSPByClientId(ctx, clientId)
	if err != nil {
		return 0, err
	}
	if !sp.Active {
		return 0, fmt.Errorf("service provider is inactive")
	}
	if sp.Public {
		return sp.ID, nil
	}
	decryptedSecret, err := S.TextDec.DecryptText(ctx, sp.ClientSecret)
	if err != nil {
		return 0, fmt.Errorf("failed to decrypt sp secret - %v", err)
	}

	if clientSecret == decryptedSecret {
		return sp.ID, nil
	}
	return 0, fmt.Errorf("invalid client secret")
}

func (S *SPStoreServiceImpl) ValidateSecretSignature(ctx context.Context, token string) (id uint, err error) {
	panic("implement me")
}

func (S *SPStoreServiceImpl) ValidatePrivateKeySignature(ctx context.Context, token string) (id uint, err error) {
	panic("implement me")
}

func (S *SPStoreServiceImpl) GetSP(ctx context.Context, id uint) (sp *models.ServiceProviderModel, err error) {
	tx := getTransaction(ctx)
	sp = &models.ServiceProviderModel{}
	result := tx.Find(sp, id)
	if result.RecordNotFound() {
		return nil, fmt.Errorf("no SP found with id %d", id)
	}
	err = result.Error
	return
}

func (S *SPStoreServiceImpl) FindSPByClientId(ctx context.Context, clientId string) (sp *models.ServiceProviderModel, err error) {
	tx := getTransaction(ctx)
	sp = &models.ServiceProviderModel{}
	result := tx.Find(sp, "client_id = ?", clientId)
	if result.RecordNotFound() {
		return nil, fmt.Errorf("no SP found with client_id %s", clientId)
	}
	err = result.Error
	return
}

func (S *SPStoreServiceImpl) FindSPByName(ctx context.Context, name string) (sp *models.ServiceProviderModel, err error) {
	tx := getTransaction(ctx)
	sp = &models.ServiceProviderModel{}
	result := tx.First(sp, "name like ?", name)
	if result.RecordNotFound() {
		return nil, fmt.Errorf("no SP found with name %s", name)
	}
	err = result.Error
	return
}

func (S *SPStoreServiceImpl) FindAllSP(ctx context.Context, page uint, pageSize uint) (sps []models.ServiceProviderModel, count uint, err error) {
	var total uint
	tx := getTransaction(ctx)
	query := tx.Select([]string{"id", "name", "description", "client_id", "active"}).Model(&models.ServiceProviderModel{})
	err = query.Limit(pageSize).Offset(pageSize * page).Find(&sps).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return sps, total, nil
}

func NewSPStoreServiceImpl(db *gorm.DB, dec ITextDecrypts, enc ITextEncrypts) ISPStoreService {
	return &SPStoreServiceImpl{Db: db, TextEnc: enc, TextDec: dec}
}
