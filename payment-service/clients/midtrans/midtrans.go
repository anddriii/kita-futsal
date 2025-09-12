package clients

import (
	"fmt"
	"time"

	errConstant "github.com/anddriii/kita-futsal/payment-service/constants/error/payment"
	"github.com/anddriii/kita-futsal/payment-service/domains/dto"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

// MidtransClient adalah struktur utama untuk melakukan komunikasi dengan Midtrans.
// Menyimpan server key dan flag apakah environment adalah production atau sandbox.
type MidtransClient struct {
	ServerKey    string // Kunci API dari Midtrans (didapat dari dashboard Midtrans)
	IsProduction bool   // True jika menggunakan mode produksi, false jika sandbox
}

// IMidtransClient adalah interface yang menyediakan kontrak fungsi untuk interaksi Midtrans,
// saat ini hanya menyediakan fungsi pembuatan payment link.
type IMidtransClient interface {
	CreatePaymentLink(request *dto.PaymentRequest) (*MidtransData, error)
}

// NewMidtransClient mengembalikan instance MidtransClient baru.
// Digunakan untuk inisialisasi Midtrans dengan konfigurasi tertentu.
func NewMidtransClient(serverKey string, isProduction bool) *MidtransClient {
	return &MidtransClient{
		ServerKey:    serverKey,
		IsProduction: isProduction,
	}
}

// CreatePaymentLink membuat transaksi pembayaran di Midtrans Snap dan mengembalikan token serta URL redirect.
// Parameter:
//   - request: pointer ke dto.PaymentRequest yang berisi detail transaksi seperti order ID, customer, item, dll.
//
// Return:
//   - *MidtransData: berisi token dan URL redirect dari Snap Midtrans jika sukses
//   - error: error jika terjadi kegagalan saat membuat transaksi
func (c *MidtransClient) CreatePaymentLink(request *dto.PaymentRequest) (*MidtransData, error) {
	var (
		snapClient   snap.Client
		isProduction = midtrans.Sandbox // Default ke sandbox
	)

	// Validasi waktu expired: harus lebih besar dari waktu sekarang
	expireDateTime := request.ExpiredAt
	currentTime := time.Now()
	duration := expireDateTime.Sub(currentTime)

	if duration <= 0 {
		logrus.Errorf("Expired at is invalid")
		return nil, errConstant.ErrExpireArInvalid
	}

	// Hitung durasi expired berdasarkan selisih waktu sekarang dan expiredAt
	expiryUnit := "minute"
	expiryDuration := int64(duration.Minutes())

	// Gunakan satuan hour jika lebih dari 1 jam
	if duration.Hours() >= 24 {
		expiryUnit = "day"
		expiryDuration = int64(duration.Hours() / 24)
	} else if duration.Hours() >= 1 {
		expiryUnit = "hour"
		expiryDuration = int64(duration.Hours())
	}

	// Ubah ke production jika flag IsProduction diset
	if c.IsProduction {
		isProduction = midtrans.Production
	}

	// Inisialisasi client Snap Midtrans dengan kunci dan environment
	snapClient.New(c.ServerKey, isProduction)

	// Siapkan list item untuk transaksi (saat ini hanya ambil item pertama)
	var items []midtrans.ItemDetails
	for _, item := range request.ItemDetails {
		items = append(items, midtrans.ItemDetails{
			ID:    item.ID,
			Name:  item.Name,
			Price: int64(item.Amount),
			Qty:   int32(item.Quantity),
		})
	}

	// Bangun request Snap Midtrans
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.OrderID,
			GrossAmt: int64(request.Amount),
		},

		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.CustomerDetail.Name,
			Email: request.CustomerDetail.Email,
			Phone: request.CustomerDetail.Phone,
		},

		Items: &items,

		Expiry: &snap.ExpiryDetails{
			Unit:     expiryUnit,
			Duration: expiryDuration,
		},
	}

	fmt.Printf("Midtrans Snap Request: %+v\n", req)

	// Kirim request transaksi ke Midtrans Snap API
	response, err := snapClient.CreateTransaction(req)
	if err != nil {
		logrus.Errorf("Failed to create transaction: %v", err)
		return nil, err
	}

	// Jika sukses, kembalikan token dan URL redirect
	return &MidtransData{
		RedirectURL: response.RedirectURL,
		Token:       response.Token,
	}, nil
}
