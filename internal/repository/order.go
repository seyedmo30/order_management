package repository

import (
	"context"

	"github.com/seyedmo30/order_management/internal/dto"
	"github.com/seyedmo30/order_management/pkg"
)

func (r *orderManagementRepository) CreateOrder(ctx context.Context, params dto.CreatOrderRepositoryRequest) (err error) {

	if err = db.WithContext(ctx).Table("orders").Create(&params).Error; err != nil {
		return
	}
	return nil
}

func (r *orderManagementRepository) GetOrderByID(ctx context.Context, orderID string) (res dto.GetOrderByIDRepositoryResponse, err error) {
	err = db.WithContext(ctx).
		Table("orders").
		Where("id = ?", orderID).
		First(&res).Error
	return
}

func (r *orderManagementRepository) UpdateOrderByID(ctx context.Context, params dto.UpdateOrderByIDRepositoryRequest) (err error) {
	err = db.WithContext(ctx).
		Table("orders").
		Updates(&params).Error
	return
}

func (r *orderManagementRepository) LockOrderOptimistic(ctx context.Context, params dto.UpdateOrderByIDRepositoryRequest) (err error) {
	result := db.WithContext(ctx).
		Table("orders").
		Where("lock = ?", false).
		Where("order_id = ?", params.OrderID).
		Update("lock", true)

	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {

		err = &pkg.ErrorCustom{
			Code:    404,
			Message: pkg.NotFoundRepositoryMessage,
		}

		return err
	}

	return
}

func (r *orderManagementRepository) GetNextHighPriorityReadyOrder(ctx context.Context) (res dto.GetNextHighPriorityReadyOrderRepositoryResponse, err error) {
	err = db.WithContext(ctx).
		Table("orders").
		Where("status = ? AND priority = ? AND lock = ?", pkg.StatusOrderManagementPending, pkg.StatusOrderManagementHigh, false).
		Order("priority DESC, created_at ASC").
		First(&res).Error
	return
}

func (r *orderManagementRepository) ListAggregateOrderReport(ctx context.Context) (counts map[string]int, err error) {
	counts = make(map[string]int)
	rows, err := db.WithContext(ctx).
		Table("orders").
		Select("status, count(*) as count").
		Group("status").Rows()

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return
}
