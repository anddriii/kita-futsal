package dto

import (
	"time"

	"github.com/anddriii/kita-futsal/payment-service/constants"
	"github.com/google/uuid"
)

// PaymentRequest merepresentasikan data yang dibutuhkan untuk membuat permintaan pembayaran (payment link)
// menggunakan Midtrans. Data ini nantinya akan dikirimkan ke service pembayaran untuk menghasilkan link pembayaran.
type PaymentRequest struct {
	PaymentLink    string          `json:"paymentLink"`    // Link pembayaran (jika ada)
	OrderID        string          `json:"orderID"`        // ID unik untuk pesanan
	ExpiredAt      time.Time       `json:"expiredAt"`      // Tanggal dan waktu kadaluarsa pembayaran
	Amount         float64         `json:"amount"`         // Jumlah total pembayaran
	Description    *string         `json:"description"`    // Deskripsi opsional tentang pembayaran
	CustomerDetail *CustomerDetail `json:"customerDetail"` // Informasi detail pelanggan
	ItemDetails    []ItemDetail    `json:"itemDetails"`    // Daftar item yang dibeli
}

// CustomerDetail menyimpan informasi pribadi pelanggan.
type CustomerDetail struct {
	Name  string `json:"name"`  // Nama pelanggan
	Email string `json:"email"` // Email pelanggan
	Phone string `json:"phone"` // Nomor telepon pelanggan
}

// ItemDetail menyimpan informasi tentang satu item dalam transaksi.
type ItemDetail struct {
	ID       string  `json:"id"`       // ID item
	Amount   float64 `json:"amount"`   // Harga item
	Name     string  `json:"name"`     // Nama item
	Quantity int     `json:"quantity"` // Jumlah item yang dibeli
}

// PaymentRequestParam digunakan untuk menerima parameter query ketika meminta daftar pembayaran,
// misalnya pada endpoint GET /payments?page=1&limit=10&sortColumn=createdAt&sortOrder=desc
type PaymentRequestParam struct {
	Page       int     `form:"page" validate:"required"`  // Halaman saat ini
	Limit      int     `form:"limit" validate:"required"` // Batas jumlah data per halaman
	SortColumn *string `form:"sortColumn"`                // Kolom untuk melakukan pengurutan
	SortOrder  *string `form:"sortOrder"`                 // Urutan pengurutan (ASC/DESC)
}

// UpdatePaymentRequest digunakan untuk memperbarui informasi status pembayaran,
// misalnya ketika menerima callback/webhook dari Midtrans.
type UpdatePaymentRequest struct {
	TransactionID *string                  `json:"transactionID"` // ID transaksi dari payment gateway
	Status        *constants.PaymentStatus `json:"status"`        // Status pembayaran
	PaidAt        *time.Time               `json:"paidAt"`        // Waktu pembayaran dilakukan
	VANumber      *string                  `json:"vaNumber"`      // Nomor Virtual Account (VA)
	Bank          *string                  `json:"bank"`          // Nama bank yang digunakan
	InvoiceLink   *string                  `json:"invoiceLink"`   // Link ke faktur (invoice)
	Acquirer      *string                  `json:"acquirer"`      // Informasi acquirer (penyedia layanan pembayaran)
}

// PaymentResponse adalah struktur data yang dikirimkan kembali ke client/merchant
// setelah permintaan pembayaran berhasil dibuat atau saat mengambil detail pembayaran.
type PaymentResponse struct {
	UUID          string                        `json:"uuid"`                    // UUID unik pembayaran
	OrderID       string                        `json:"orderID"`                 // ID pesanan
	Amount        float64                       `json:"amount"`                  // Total jumlah pembayaran
	Status        constants.PaymentStatusString `json:"status"`                  // Status pembayaran dalam bentuk string
	PaymentLink   string                        `json:"paymentLink"`             // Link untuk melakukan pembayaran
	InvoiceLink   *string                       `json:"invoiceLink,omitempty"`   // Link faktur jika tersedia
	TransactionID *string                       `json:"transactionID,omitempty"` // ID transaksi dari gateway
	VANumber      *string                       `json:"vaNumber,omitempty"`      // Nomor Virtual Account
	Bank          *string                       `json:"bank,omitempty"`          // Nama bank yang digunakan
	Acquirer      *string                       `json:"acquirer,omitempty"`      // Nama acquirer (penyedia gateway)
	Description   *string                       `json:"description"`             // Deskripsi pembayaran
	PaidAt        *time.Time                    `json:"paidAt,omitempty"`        // Waktu pembayaran
	ExpiredAt     *time.Time                    `json:"expiredAt,omitempty"`     // Waktu kadaluarsa pembayaran
	CreatedAt     *time.Time                    `json:"createdAt"`               // Waktu saat pembayaran dibuat
	UpdatedAt     *time.Time                    `json:"updatedAt"`               // Waktu saat pembayaran diperbarui
}

// Webhook merepresentasikan payload yang diterima dari Midtrans saat ada notifikasi status pembayaran.
// Struktur ini disesuaikan dengan format JSON yang dikirimkan oleh Midtrans (notification API).
type Webhook struct {
	VANumber          []VANumber                    `json:"va_numbers"`         // Daftar nomor VA
	TransactionTime   string                        `json:"transaction_time"`   // Waktu transaksi dilakukan
	TransactionStatus constants.PaymentStatusString `json:"transaction_status"` // Status transaksi
	TransactionID     string                        `json:"transaction_id"`     // ID transaksi dari gateway
	StatusMessage     string                        `json:"status_message"`     // Pesan status transaksi
	StatusCode        string                        `json:"status_code"`        // Kode status transaksi
	SignatureKey      string                        `json:"signature_key"`      // Kunci signature untuk verifikasi
	SettlementTime    string                        `json:"settlement_time"`    // Waktu settlement pembayaran
	PaymentType       string                        `json:"payment_type"`       // Jenis pembayaran (bank transfer, e-wallet, dll)
	PaymentAmount     []PaymentAmount               `json:"payment_amount"`     // Detail jumlah pembayaran
	OrderID           uuid.UUID                     `json:"order_id"`           // UUID pesanan
	MerchantID        string                        `json:"merchant_id"`        // ID merchant dari gateway
	GrossAmount       string                        `json:"gross_amount"`       // Jumlah kotor pembayaran
	FraudStatus       string                        `json:"fraud_status"`       // Status deteksi fraud
	Currency          string                        `json:"currency"`           // Mata uang pembayaran
	Acquirer          *string                       `json:"acquirer"`           // Acquirer opsional
}

// VANumber merepresentasikan informasi satu nomor Virtual Account dan bank terkait.
type VANumber struct {
	VaNumber string `json:"va_number"` // Nomor Virtual Account
	Bank     string `json:"bank"`      // Nama bank
}

// PaymentAmount menyimpan informasi tentang jumlah yang dibayarkan dan waktu pembayaran.
type PaymentAmount struct {
	PaidAt *string `json:"paid_at"` // Waktu pembayaran dilakukan (dalam string format)
	Amount *string `json:"amount"`  // Jumlah pembayaran (dalam string format)
}
