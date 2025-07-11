package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type InstallmentRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	FindByUUID(ctx context.Context, uuid uuid.UUID) (installment *entity.TransactionInstallment, err error)
	FindByID(ctx context.Context, id uint) (installment *entity.TransactionInstallment, err error)
	FindAll(ctx context.Context, query *model.QueryGet) (installments *[]entity.TransactionInstallment, err error)
	FindAllByTransactionID(ctx context.Context, query *model.QueryGet, transaction_id uint) (installments *[]entity.TransactionInstallment, err error)
	FindAllByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (installments *[]entity.TransactionInstallment, err error)
	Count(ctx context.Context, query *model.QueryGet) (total int64)
	CountByTransactionID(ctx context.Context, query *model.QueryGet, transaction_id uint) (total int64)
	CountByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (total int64)
	CountUnscoped(ctx context.Context, query *model.QueryGet) (total int64)
	BulkInsertWithTransaction(ctx context.Context, tx *gorm.DB, installments []entity.TransactionInstallment) error
	InsertWithTransaction(ctx context.Context, tx *gorm.DB, installment *entity.TransactionInstallment) error
	UpdateWithTransaction(ctx context.Context, tx *gorm.DB, installment *entity.TransactionInstallment) error
	DeleteWithTransaction(ctx context.Context, tx *gorm.DB, installment *entity.TransactionInstallment) error
	CancelWithTransaction(ctx context.Context, tx *gorm.DB, transaction_id uint) error
}

type installmentRepository struct {
	*gorm.DB
}

func NewIntallmentRepository(db *gorm.DB) InstallmentRepository {
	return &installmentRepository{DB: db}
}

func (r *installmentRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.DB.Begin()
}

func (r *installmentRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (transaction *entity.TransactionInstallment, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Preload("Transaction").Preload("Transaction.User").Preload("Payments").
		Find(&transaction); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return
}

func (r *installmentRepository) FindByID(ctx context.Context, id uint) (transaction *entity.TransactionInstallment, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if result := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).
		Preload("Transaction").Preload("Transaction.User").Preload("Payments").
		Find(&transaction); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return
}

func (r *installmentRepository) FindAll(ctx context.Context, query *model.QueryGet) (transactions *[]entity.TransactionInstallment, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.TransactionInstallment{}).Preload("Transaction")

	var allowedFields = map[string]string{
		"created": "transaction_installments.created_at",
		"updated": "transaction_installments.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&transactions).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return
}

func (r *installmentRepository) FindAllByTransactionID(ctx context.Context, query *model.QueryGet, transaction_id uint) (transactions *[]entity.TransactionInstallment, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.TransactionInstallment{}).
		Preload("Transaction").Where("transaction_id = ?", transaction_id)

	var allowedFields = map[string]string{
		"created": "transaction_installments.created_at",
		"updated": "transaction_installments.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&transactions).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return
}

func (r *installmentRepository) FindAllByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (transactions *[]entity.TransactionInstallment, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.TransactionInstallment{}).
		Joins("JOIN transactions on transactions.id = transaction_installments.transaction_id").
		Where("transactions.user_id = ?", user_id).Preload("Transaction")

	var allowedFields = map[string]string{
		"created": "transaction_installments.created_at",
		"updated": "transaction_installments.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&transactions).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return
}

func (r *installmentRepository) Count(ctx context.Context, query *model.QueryGet) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.TransactionInstallment{})

	var allowedFields = map[string]string{
		"created": "transaction_installments.created_at",
		"updated": "transaction_installments.updated_at",
	}

	tx = tx.Scopes(
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Count(&total).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
	}

	return
}

func (r *installmentRepository) CountByTransactionID(ctx context.Context, query *model.QueryGet, transaction_id uint) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.TransactionInstallment{}).Where("transaction_id = ?", transaction_id)

	var allowedFields = map[string]string{
		"created": "transaction_installments.created_at",
		"updated": "transaction_installments.updated_at",
	}

	tx = tx.Scopes(
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Count(&total).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
	}

	return
}

func (r *installmentRepository) CountByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.TransactionInstallment{}).
		Joins("JOIN transactions on transactions.id = transaction_installments.transaction_id").
		Where("transactions.user_id = ?", user_id)

	var allowedFields = map[string]string{
		"created": "transaction_installments.created_at",
		"updated": "transaction_installments.updated_at",
	}

	tx = tx.Scopes(
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Count(&total).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
	}

	return
}

func (r *installmentRepository) CountUnscoped(ctx context.Context, query *model.QueryGet) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.TransactionInstallment{}).Unscoped()

	var allowedFields = map[string]string{
		"created": "transaction_installments.created_at",
		"updated": "transaction_installments.updated_at",
	}

	tx = tx.Scopes(
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Count(&total).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
	}

	return
}

func (r *installmentRepository) InsertWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.TransactionInstallment) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	return tx.WithContext(ctx).Create(transaction).Error
}

func (r *installmentRepository) BulkInsertWithTransaction(ctx context.Context, tx *gorm.DB, transactions []entity.TransactionInstallment) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	return tx.WithContext(ctx).Create(transactions).Error
}

func (r *installmentRepository) UpdateWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.TransactionInstallment) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := tx.WithContext(ctx).Where("id = ?", transaction.ID).Updates(transaction).
		Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
		return err
	}

	return nil
}

func (r *installmentRepository) CancelWithTransaction(ctx context.Context, tx *gorm.DB, transaction_id uint) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := tx.WithContext(ctx).Where("transaction_id = ? AND payment_status IN ?", transaction_id, []string{"pending", "partial", "overdue"}).
		Updates(entity.TransactionInstallment{PaymentStatus: entity.PaymentStatusFailed}).
		Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
		return err
	}

	return nil
}

func (r *installmentRepository) DeleteWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.TransactionInstallment) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := tx.WithContext(ctx).Where("id = ?", transaction.ID).Delete(transaction).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
		return err
	}

	return nil
}
