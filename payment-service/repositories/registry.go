package repositories

import (
	repositories "github.com/anddriii/kita-futsal/payment-service/repositories/payment"
	repositories2 "github.com/anddriii/kita-futsal/payment-service/repositories/paymenthistory"
	"gorm.io/gorm"
)

// Registry adalah implementasi dari IRepositoryRegistry.
// Struct ini berfungsi sebagai dependency container untuk repositori,
// mempermudah pengelolaan dan injeksi dependency database.
type Registry struct {
	db *gorm.DB // koneksi database utama yang digunakan oleh semua repository
}

// IRepositoryRegistry adalah interface yang mendefinisikan akses ke berbagai repository
// serta memberikan akses ke objek transaksi database (gorm.DB).
type IRepositoryRegistry interface {
	GetPayment() repositories.IPaymentRepository
	GetPaymentHistory() repositories2.IPaymentHistoryRepository
	GetTx() *gorm.DB
}

// NewRepositoryRegistry menginisialisasi dan mengembalikan instance Registry sebagai IRepositoryRegistry.
// Parameter:
//   - db: objek koneksi database GORM yang dibagikan ke semua repository.
func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

// GetPayment mengembalikan instance dari PaymentRepository.
// Ini memungkinkan service mengakses fungsi-fungsi terkait entitas pembayaran.
func (r *Registry) GetPayment() repositories.IPaymentRepository {
	return repositories.NewPaymentRepository(r.db)
}

// GetPaymentHistory mengembalikan instance dari PaymentHistoryRepository.
// Ini digunakan untuk mengakses histori pembayaran.
func (r *Registry) GetPaymentHistory() repositories2.IPaymentHistoryRepository {
	return repositories2.NewPaymentHistoryRepository(r.db)
}

// GetTx mengembalikan objek koneksi database GORM untuk kebutuhan transaksi manual.
// Biasanya digunakan ketika service ingin menjalankan operasi DB dalam satu transaksi.
func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
