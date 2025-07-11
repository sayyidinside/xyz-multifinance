package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	FindByUUID(ctx context.Context, uuid uuid.UUID) (payment *entity.Payment, err error)
	FindAll(ctx context.Context, query *model.QueryGet) (payments *[]entity.Payment, err error)
	FindAllByTransactionID(ctx context.Context, query *model.QueryGet, transaction_id uint) (payments *[]entity.Payment, err error)
	FindAllByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (payments *[]entity.Payment, err error)
	Count(ctx context.Context, query *model.QueryGet) (total int64)
	CountByTransactionID(ctx context.Context, query *model.QueryGet, transaction_id uint) (total int64)
	CountByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (total int64)
	CountUnscoped(ctx context.Context, query *model.QueryGet) (total int64)
	BulkInsertWithTransaction(ctx context.Context, tx *gorm.DB, payments []entity.Payment) error
	InsertWithTransaction(ctx context.Context, tx *gorm.DB, payment *entity.Payment) error
	UpdateWithTransaction(ctx context.Context, tx *gorm.DB, payment *entity.Payment) error
	DeleteWithTransaction(ctx context.Context, tx *gorm.DB, payment *entity.Payment) error
}

type paymentRepository struct {
	*gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{DB: db}
}

func (r *paymentRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.DB.Begin()
}

func (r *paymentRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (transaction *entity.Payment, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Preload("Transaction").Preload("Transaction.User").Preload("Installment").
		Find(&transaction); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return
}

func (r *paymentRepository) FindAll(ctx context.Context, query *model.QueryGet) (transactions *[]entity.Payment, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Payment{}).
		Preload("Transaction").Preload("Installment")

	var allowedFields = map[string]string{
		"created": "payments.created_at",
		"updated": "payments.updated_at",
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

func (r *paymentRepository) FindAllByTransactionID(ctx context.Context, query *model.QueryGet, transaction_id uint) (transactions *[]entity.Payment, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Payment{}).
		Preload("Transaction").Where("transaction_id = ?", transaction_id)

	var allowedFields = map[string]string{
		"created": "payments.created_at",
		"updated": "payments.updated_at",
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

func (r *paymentRepository) FindAllByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (transactions *[]entity.Payment, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Payment{}).
		Preload("Transaction").Preload("Installment").
		Joins("JOIN transactions on transactions.id = payments.transaction_id").
		Where("transactions.user_id = ?", user_id).Preload("Transaction")

	var allowedFields = map[string]string{
		"created": "payments.created_at",
		"updated": "payments.updated_at",
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

func (r *paymentRepository) Count(ctx context.Context, query *model.QueryGet) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Payment{})

	var allowedFields = map[string]string{
		"created": "payments.created_at",
		"updated": "payments.updated_at",
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

func (r *paymentRepository) CountByTransactionID(ctx context.Context, query *model.QueryGet, transaction_id uint) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Payment{}).Where("transaction_id = ?", transaction_id)

	var allowedFields = map[string]string{
		"created": "payments.created_at",
		"updated": "payments.updated_at",
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

func (r *paymentRepository) CountByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Payment{}).
		Joins("JOIN transactions on transactions.id = payments.transaction_id").
		Where("transactions.user_id = ?", user_id)

	var allowedFields = map[string]string{
		"created": "payments.created_at",
		"updated": "payments.updated_at",
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

func (r *paymentRepository) CountUnscoped(ctx context.Context, query *model.QueryGet) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Payment{}).Unscoped()

	var allowedFields = map[string]string{
		"created": "payments.created_at",
		"updated": "payments.updated_at",
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

func (r *paymentRepository) InsertWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.Payment) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	return tx.WithContext(ctx).Create(transaction).Error
}

func (r *paymentRepository) BulkInsertWithTransaction(ctx context.Context, tx *gorm.DB, transactions []entity.Payment) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	return tx.WithContext(ctx).Create(transactions).Error
}

func (r *paymentRepository) UpdateWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.Payment) error {
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

func (r *paymentRepository) DeleteWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.Payment) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := tx.WithContext(ctx).Where("id = ?", transaction.ID).Delete(transaction).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
		return err
	}

	return nil
}
