package repo

import (
	"CRM/internal/deal_service/contracts"
	"CRM/internal/deal_service/model"
	"context"
	"gorm.io/gorm"
)

type DealRepository struct {
	db *gorm.DB
}

func NewDealRepository(db *gorm.DB) *DealRepository {
	return &DealRepository{db: db}
}

func (r *DealRepository) ExecTxn(ctx context.Context, fn func(repo *DealRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewDealRepository(tx))
	})
}

func (r *DealRepository) CreateDeal(deal *model.Deal) error {
	return r.db.Create(deal).Error
}

func (r *DealRepository) GetByIDAndOwner(id, ownerID string) (*model.Deal, error) {
	var deal model.Deal
	err := r.db.First(&deal, "id = ? AND owner_id = ?", id, ownerID).Error
	return &deal, err
}

func (r *DealRepository) ListByOwner(ownerID string, req contracts.ListDealsRequest) ([]model.Deal, int64, error) {
	var deals []model.Deal
	var totalCount int64

	query := r.db.Model(&model.Deal{}).Where("owner_id = ?", ownerID)
	if req.Stage != "" {
		query = query.Where("stage = ?", req.Stage)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	err := query.Find(&deals).Error
	return deals, totalCount, err
}

func (r *DealRepository) UpdateDeal(deal *model.Deal) error {
	return r.db.Save(deal).Error
}

func (r *DealRepository) DeleteByIDAndOwner(id, ownerID string) (int64, error) {
	result := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.Deal{})
	return result.RowsAffected, result.Error
}

// ToDealApiResponse converts a database model into an API response contract.
func (r *DealRepository) ToDealApiResponse(deal *model.Deal) *contracts.DealResponse {
	return &contracts.DealResponse{
		ID:                deal.ID,
		Name:              deal.Name,
		Stage:             deal.Stage,
		Value:             deal.Value,
		ExpectedCloseDate: deal.ExpectedCloseDate.Format("2006-01-02"),
		ContactID:         deal.ContactID,
		OwnerID:           deal.OwnerID,
		CreatedAt:         deal.CreatedAt,
		UpdatedAt:         deal.UpdatedAt,
	}
}
