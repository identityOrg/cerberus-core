package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/jinzhu/gorm"
)

type ScopeClaimStoreServiceImpl struct {
	Db *gorm.DB
}

func NewScopeClaimStoreServiceImpl(db *gorm.DB) *ScopeClaimStoreServiceImpl {
	return &ScopeClaimStoreServiceImpl{Db: db}
}

func (s *ScopeClaimStoreServiceImpl) BeginTransaction(ctx context.Context, readOnly bool) context.Context {
	return beginTransaction(ctx, readOnly, s.Db)
}

func (s *ScopeClaimStoreServiceImpl) CommitTransaction(ctx context.Context) context.Context {
	return commitTransaction(ctx)
}

func (s *ScopeClaimStoreServiceImpl) RollbackTransaction(ctx context.Context) context.Context {
	return rollbackTransaction(ctx)
}

func (s *ScopeClaimStoreServiceImpl) CreateScope(ctx context.Context, name string, description string) (id uint, err error) {
	scope := &models.ScopeModel{
		Name:        name,
		Description: description,
	}
	db := getTransaction(ctx)
	saveResult := db.Save(scope)
	return scope.ID, saveResult.Error
}

func (s *ScopeClaimStoreServiceImpl) FindScopeByName(ctx context.Context, name string) (*models.ScopeModel, error) {
	tx := getTransaction(ctx)
	scope := &models.ScopeModel{}
	result := tx.First(scope, "name like ?", name)
	if result.RecordNotFound() {
		return nil, fmt.Errorf("no scope found with name %s", name)
	}
	return scope, result.Error
}

func (s *ScopeClaimStoreServiceImpl) GetScope(ctx context.Context, id uint) (*models.ScopeModel, error) {
	tx := getTransaction(ctx)
	scope := &models.ScopeModel{}
	result := tx.Find(scope, id)
	if result.RecordNotFound() {
		return nil, fmt.Errorf("no scope found with id %d", id)
	}
	return scope, result.Error
}

func (s *ScopeClaimStoreServiceImpl) GetAllScopes(ctx context.Context, page uint, pageSize uint) ([]*models.ScopeModel, uint, error) {
	var total uint
	tx := getTransaction(ctx)
	scopes := make([]*models.ScopeModel, 0)
	query := tx.Model(&models.ServiceProviderModel{})
	err := query.Limit(pageSize).Offset(pageSize * page).Find(&scopes).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return scopes, total, nil
}

func (s *ScopeClaimStoreServiceImpl) UpdateScope(ctx context.Context, id uint, description string) error {
	scope := &models.ScopeModel{}
	scope.ID = id
	db := getTransaction(ctx)
	findResult := db.Find(scope)
	if findResult.RecordNotFound() {
		return errors.New("scope not found")
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	scope.Description = description
	return db.Save(scope).Error
}

func (s *ScopeClaimStoreServiceImpl) DeleteScope(ctx context.Context, id uint) error {
	scope := &models.ScopeModel{}
	scope.ID = id
	db := getTransaction(ctx)
	return db.Delete(scope).Error
}

func (s *ScopeClaimStoreServiceImpl) AddClaimToScope(ctx context.Context, scopeId uint, claimId uint) error {
	db := getTransaction(ctx)
	scope := &models.ScopeModel{}
	scope.ID = scopeId
	findResult := db.Find(scope)
	if findResult.RecordNotFound() {
		return fmt.Errorf("scope with id %d not found", scopeId)
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	claim := &models.ClaimModel{}
	claim.ID = claimId
	findResult = db.Find(claim)
	if findResult.RecordNotFound() {
		return fmt.Errorf("claim with id %d not found", claimId)
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	claimAssociation := db.Model(scope).Association("Claims")
	if claimAssociation.Error != nil {
		return claimAssociation.Error
	}
	return claimAssociation.Append(claim).Error
}

func (s *ScopeClaimStoreServiceImpl) RemoveClaimFromScope(ctx context.Context, scopeId uint, claimId uint) error {
	db := getTransaction(ctx)
	scope := &models.ScopeModel{}
	scope.ID = scopeId
	findResult := db.Find(scope)
	if findResult.RecordNotFound() {
		return fmt.Errorf("scope with id %d not found", scopeId)
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	claim := &models.ClaimModel{}
	claim.ID = claimId
	findResult = db.Find(claim)
	if findResult.RecordNotFound() {
		return fmt.Errorf("claim with id %d not found", claimId)
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	claimAssociation := db.Model(scope).Association("Claims")
	if claimAssociation.Error != nil {
		return claimAssociation.Error
	}
	return claimAssociation.Delete(claim).Error
}

func (s *ScopeClaimStoreServiceImpl) CreateClaim(ctx context.Context, name string, description string) (id uint, err error) {
	claim := &models.ClaimModel{
		Name:        name,
		Description: description,
	}
	db := getTransaction(ctx)
	saveResult := db.Save(claim)
	return claim.ID, saveResult.Error
}

func (s *ScopeClaimStoreServiceImpl) FindClaimByName(ctx context.Context, name string) (*models.ClaimModel, error) {
	tx := getTransaction(ctx)
	claim := &models.ClaimModel{}
	result := tx.First(claim, "name like ?", name)
	if result.RecordNotFound() {
		return nil, fmt.Errorf("no claim found with name %s", name)
	}
	return claim, result.Error
}

func (s *ScopeClaimStoreServiceImpl) GetClaim(ctx context.Context, id uint) (*models.ClaimModel, error) {
	tx := getTransaction(ctx)
	claim := &models.ClaimModel{}
	result := tx.Find(claim, id)
	if result.RecordNotFound() {
		return nil, fmt.Errorf("no claim found with id %d", id)
	}
	return claim, result.Error
}

func (s *ScopeClaimStoreServiceImpl) GetAllClaims(ctx context.Context, page uint, pageSize uint) ([]*models.ClaimModel, uint, error) {
	var total uint
	tx := getTransaction(ctx)
	claims := make([]*models.ClaimModel, 0)
	query := tx.Model(&models.ClaimModel{})
	err := query.Limit(pageSize).Offset(pageSize * page).Find(&claims).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return claims, total, nil
}

func (s *ScopeClaimStoreServiceImpl) UpdateClaim(ctx context.Context, id uint, description string) error {
	claim := &models.ClaimModel{}
	claim.ID = id
	db := getTransaction(ctx)
	findResult := db.Find(claim)
	if findResult.RecordNotFound() {
		return errors.New("claim not found")
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	claim.Description = description
	return db.Save(claim).Error
}

func (s *ScopeClaimStoreServiceImpl) DeleteClaim(ctx context.Context, id uint) error {
	claim := &models.ClaimModel{}
	claim.ID = id
	db := getTransaction(ctx)
	return db.Delete(claim).Error
}
