package core

import (
	"context"
	"fmt"
	"github.com/identityOrg/cerberus-core/models"
	"gorm.io/gorm"
)

type ScopeClaimStoreServiceImpl struct {
	Db *gorm.DB
}

func NewScopeClaimStoreServiceImpl(db *gorm.DB) *ScopeClaimStoreServiceImpl {
	return &ScopeClaimStoreServiceImpl{Db: db}
}

func (s *ScopeClaimStoreServiceImpl) CreateScope(ctx context.Context, name string, description string) (id uint, err error) {
	scope := &models.ScopeModel{
		Name:        name,
		Description: description,
	}
	db := s.Db.WithContext(ctx)
	saveResult := db.Save(scope)
	return scope.ID, saveResult.Error
}

func (s *ScopeClaimStoreServiceImpl) FindScopeByName(ctx context.Context, name string) (*models.ScopeModel, error) {
	tx := s.Db.WithContext(ctx)
	scope := &models.ScopeModel{}
	result := tx.First(scope, "name like ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("scope not found with name %s", name)
	}
	return scope, result.Error
}

func (s *ScopeClaimStoreServiceImpl) GetScope(ctx context.Context, id uint) (*models.ScopeModel, error) {
	tx := s.Db.WithContext(ctx)
	scope := &models.ScopeModel{}
	result := tx.Find(scope, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("scope not found with id %d", id)
	}
	return scope, result.Error
}

func (s *ScopeClaimStoreServiceImpl) GetAllScopes(ctx context.Context, page uint, pageSize uint) ([]*models.ScopeModel, uint, error) {
	var total int64
	tx := s.Db.WithContext(ctx)
	scopes := make([]*models.ScopeModel, 0)
	query := tx.Model(&models.ServiceProviderModel{})
	err := query.Limit(int(pageSize)).Offset(int(pageSize * page)).Find(&scopes).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return scopes, uint(total), nil
}

func (s *ScopeClaimStoreServiceImpl) UpdateScope(ctx context.Context, id uint, description string) error {
	scope := &models.ScopeModel{}
	scope.ID = id
	db := s.Db.WithContext(ctx)
	findResult := db.Find(scope)
	if findResult.Error != nil {
		return findResult.Error
	}
	if findResult.RowsAffected != 1 {
		return fmt.Errorf("scope not found with id %d", id)
	}
	scope.Description = description
	return db.Save(scope).Error
}

func (s *ScopeClaimStoreServiceImpl) DeleteScope(ctx context.Context, id uint) error {
	scope := &models.ScopeModel{}
	scope.ID = id
	db := s.Db.WithContext(ctx)
	return db.Delete(scope).Error
}

func (s *ScopeClaimStoreServiceImpl) AddClaimToScope(ctx context.Context, scopeId uint, claimId uint) error {
	db := s.Db.WithContext(ctx)
	scope := &models.ScopeModel{}
	scope.ID = scopeId
	findResult := db.Find(scope)
	if findResult.Error != nil {
		return findResult.Error
	}
	if findResult.RowsAffected != 1 {
		return fmt.Errorf("scope not found with id %d", scopeId)
	}
	claim := &models.ClaimModel{}
	claim.ID = claimId
	findResult = db.Find(claim)
	if findResult.Error != nil {
		return findResult.Error
	}
	if findResult.RowsAffected != 1 {
		return fmt.Errorf("claim not found with id %d", claimId)
	}
	claimAssociation := db.Model(scope).Association("Claims")
	if claimAssociation.Error != nil {
		return claimAssociation.Error
	}
	return claimAssociation.Append(claim)
}

func (s *ScopeClaimStoreServiceImpl) RemoveClaimFromScope(ctx context.Context, scopeId uint, claimId uint) error {
	db := s.Db.WithContext(ctx)
	scope := &models.ScopeModel{}
	scope.ID = scopeId
	findResult := db.Find(scope)
	if findResult.Error != nil {
		return findResult.Error
	}
	if findResult.RowsAffected != 1 {
		return fmt.Errorf("scope not found with id %d", scopeId)
	}
	claim := &models.ClaimModel{}
	claim.ID = claimId
	findResult = db.Find(claim)
	if findResult.Error != nil {
		return findResult.Error
	}
	if findResult.RowsAffected != 1 {
		return fmt.Errorf("claim not found with id %d", claimId)
	}
	claimAssociation := db.Model(scope).Association("Claims")
	if claimAssociation.Error != nil {
		return claimAssociation.Error
	}
	return claimAssociation.Delete(claim)
}

func (s *ScopeClaimStoreServiceImpl) CreateClaim(ctx context.Context, name string, description string) (id uint, err error) {
	claim := &models.ClaimModel{
		Name:        name,
		Description: description,
	}
	db := s.Db.WithContext(ctx)
	saveResult := db.Save(claim)
	return claim.ID, saveResult.Error
}

func (s *ScopeClaimStoreServiceImpl) FindClaimByName(ctx context.Context, name string) (*models.ClaimModel, error) {
	tx := s.Db.WithContext(ctx)
	claim := &models.ClaimModel{}
	result := tx.First(claim, "name like ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("claim not found with name %s", name)
	}
	return claim, result.Error
}

func (s *ScopeClaimStoreServiceImpl) GetClaim(ctx context.Context, id uint) (*models.ClaimModel, error) {
	tx := s.Db.WithContext(ctx)
	claim := &models.ClaimModel{}
	result := tx.Find(claim, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, fmt.Errorf("claim not found with is %d", id)
	}
	return claim, nil
}

func (s *ScopeClaimStoreServiceImpl) GetAllClaims(ctx context.Context, page uint, pageSize uint) ([]*models.ClaimModel, uint, error) {
	var total int64
	tx := s.Db.WithContext(ctx)
	claims := make([]*models.ClaimModel, 0)
	query := tx.Model(&models.ClaimModel{})
	err := query.Limit(int(pageSize)).Offset(int(pageSize * page)).Find(&claims).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return claims, uint(total), nil
}

func (s *ScopeClaimStoreServiceImpl) UpdateClaim(ctx context.Context, id uint, description string) error {
	claim := &models.ClaimModel{}
	claim.ID = id
	db := s.Db.WithContext(ctx)
	findResult := db.Find(claim)
	if findResult.Error != nil {
		return findResult.Error
	}
	if findResult.RowsAffected != 1 {
		return fmt.Errorf("claim not found with id %d", id)
	}
	claim.Description = description
	return db.Save(claim).Error
}

func (s *ScopeClaimStoreServiceImpl) DeleteClaim(ctx context.Context, id uint) error {
	claim := &models.ClaimModel{}
	claim.ID = id
	db := s.Db.WithContext(ctx)
	return db.Delete(claim).Error
}
