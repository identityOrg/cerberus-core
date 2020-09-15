package core

import (
	"context"
	"fmt"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/jinzhu/gorm"
)

type SPStoreServiceImpl struct {
	Db *gorm.DB
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
	panic("implement me")
}

func (S *SPStoreServiceImpl) UpdateSP(ctx context.Context, id uint, metadata *models.ServiceProviderMetadata) (err error) {
	panic("implement me")
}

func (S *SPStoreServiceImpl) PatchSP(ctx context.Context, id uint, metadata *models.ServiceProviderMetadata) (err error) {
	panic("implement me")
}

func (S *SPStoreServiceImpl) DeleteSP(ctx context.Context, id uint) (err error) {
	panic("implement me")
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
	panic("implement me")
}

func (S *SPStoreServiceImpl) ValidateClientCredentials(ctx context.Context, clientId, clientSecret string) (id uint, err error) {
	panic("implement me")
}

func (S *SPStoreServiceImpl) ValidateSecretSignature(ctx context.Context, token string) (id uint, err error) {
	panic("implement me")
}

func (S *SPStoreServiceImpl) ValidatePrivateKeySignature(ctx context.Context, token string) (id uint, err error) {
	panic("implement me")
}

func (S *SPStoreServiceImpl) GetSP(ctx context.Context, id uint) (sp *models.ServiceProviderModel, err error) {
	tx := getTransaction(ctx)
	result := tx.Find(sp, id)
	if result.RecordNotFound() {
		return nil, fmt.Errorf("no SP found with id %d", id)
	}
	err = result.Error
	return
}

func (S *SPStoreServiceImpl) FindSPByClientId(ctx context.Context, clientId string) (sp *models.ServiceProviderModel, err error) {
	tx := getTransaction(ctx)
	result := tx.Find(sp, "client_id = ?", clientId)
	if result.RecordNotFound() {
		return nil, fmt.Errorf("no SP found with client_id %s", clientId)
	}
	err = result.Error
	return
}

func (S *SPStoreServiceImpl) FindSPByName(ctx context.Context, name string) (sp *models.ServiceProviderModel, err error) {
	tx := getTransaction(ctx)
	result := tx.First(sp, "name = ?", name)
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

func NewSPStoreServiceImpl(db *gorm.DB) ISPStoreService {
	return &SPStoreServiceImpl{Db: db}
}
