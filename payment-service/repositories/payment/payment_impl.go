package repositories

import (
	"context"
	"errors"
	"fmt"

	errWrap "github.com/anddriii/kita-futsal/payment-service/common/error"
	"github.com/anddriii/kita-futsal/payment-service/constants"
	errConst "github.com/anddriii/kita-futsal/payment-service/constants/error"
	errPayment "github.com/anddriii/kita-futsal/payment-service/constants/error/payment"
	"github.com/anddriii/kita-futsal/payment-service/domains/dto"
	"github.com/anddriii/kita-futsal/payment-service/domains/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentRepository adalah implementasi dari IPaymentRepository
// yang bertanggung jawab mengelola operasi database terkait entitas Payment.
type PaymentRepository struct {
	db *gorm.DB // koneksi database yang digunakan
}

// NewPaymentRepository mengembalikan instance baru PaymentRepository.
// Parameter:
//   - db: koneksi GORM ke database
func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

// Create menyimpan data pembayaran baru ke database.
// Parameter:
//   - ctx: context untuk kontrol lifecycle
//   - db: instance transaksi database (biasanya dari service)
//   - req: data request pembayaran (PaymentRequest)
//
// Return:
//   - models.Payment: objek Payment yang disimpan
//   - error: jika gagal menyimpan
func (p *PaymentRepository) Create(ctx context.Context, db *gorm.DB, req *dto.PaymentRequest) (*models.Payment, error) {
	status := constants.Initial
	orderID := uuid.MustParse(req.OrderID) // pastikan orderID valid sebagai UUID

	Payment := models.Payment{
		UUID:        uuid.New(), // generate UUID baru untuk Payment
		OrderID:     orderID,
		Amount:      req.Amount,
		PaymentLink: req.PaymentLink,
		ExpiredAt:   &req.ExpiredAt,
		Description: req.Description,
		Status:      &status,
	}

	err := db.WithContext(ctx).Create(&Payment).Error
	if err != nil {
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return &Payment, nil
}

// FindAllWithPagination mengambil daftar pembayaran dari database berdasarkan parameter paginasi dan sort.
// Parameter:
//   - ctx: context untuk lifecycle
//   - param: parameter pencarian dan paginasi (limit, page, sort)
//
// Return:
//   - []models.Payment: daftar data pembayaran
//   - int64: total data (sebelum paginasi)
//   - error: jika terjadi kesalahan database
func (p *PaymentRepository) FindAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) ([]models.Payment, int64, error) {
	var (
		fields []models.Payment
		sort   string
		total  int64
	)

	// Atur sorting berdasarkan parameter, default: created_at desc
	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit

	// Ambil data dengan paginasi
	err := p.db.
		WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&fields).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConst.ErrSQLError)
	}

	// Hitung total data (tanpa paginasi)
	err = p.db.
		WithContext(ctx).
		Model(&fields).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConst.ErrSQLError)
	}

	return fields, total, nil
}

// FindByOrderID mencari data Payment berdasarkan Order ID.
// Parameter:
//   - ctx: context
//   - orderID: ID unik dari transaksi/order
//
// Return:
//   - *models.Payment: data pembayaran jika ditemukan
//   - error: jika tidak ditemukan atau terjadi kesalahan DB
func (p *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	var payment models.Payment

	err := p.db.
		WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&payment).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return &payment, nil
}

// FindByUUID mencari data Payment berdasarkan UUID.
// Parameter:
//   - ctx: context
//   - uuid: UUID unik dari Payment
//
// Return:
//   - *models.Payment: data pembayaran jika ditemukan
//   - error: jika tidak ditemukan atau gagal query
func (p *PaymentRepository) FindByUUID(ctx context.Context, uuid string) (*models.Payment, error) {
	var payment models.Payment

	err := p.db.
		WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&payment).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return &payment, nil
}

// Update memperbarui data Payment berdasarkan order ID.
// Parameter:
//   - ctx: context
//   - db: database transaction instance
//   - orderID: ID unik transaksi
//   - req: data yang akan diupdate (status, paid_at, va_number, dsb.)
//
// Return:
//   - *models.Payment: data setelah diupdate (namun tidak di-reload dari DB)
//   - error: jika gagal update
func (p *PaymentRepository) Update(ctx context.Context, db *gorm.DB, orderID string, req *dto.UpdatePaymentRequest) (*models.Payment, error) {
	payment := models.Payment{
		Status:        req.Status,
		TransactionID: req.TransactionID,
		InvoiceLink:   req.InvoiceLink,
		PaidAt:        req.PaidAt,
		VANumber:      req.VANumber,
		Bank:          req.Bank,
		Acquirer:      req.Acquirer,
	}

	err := db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Updates(&payment).Error

	if err != nil {
		return nil, errWrap.WrapError(errConst.ErrSQLError)
	}

	return &payment, nil
}
