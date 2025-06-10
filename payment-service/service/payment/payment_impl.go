package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	clients "github.com/anddriii/kita-futsal/payment-service/clients/midtrans"
	"github.com/anddriii/kita-futsal/payment-service/common/gcs"
	"github.com/anddriii/kita-futsal/payment-service/common/util"
	config2 "github.com/anddriii/kita-futsal/payment-service/config"
	"github.com/anddriii/kita-futsal/payment-service/constants"
	errPayment "github.com/anddriii/kita-futsal/payment-service/constants/error/payment"
	"github.com/anddriii/kita-futsal/payment-service/controllers/kafka"
	"github.com/anddriii/kita-futsal/payment-service/domains/dto"
	"github.com/anddriii/kita-futsal/payment-service/domains/models"
	"github.com/anddriii/kita-futsal/payment-service/repositories"

	"gorm.io/gorm"
)

type PaymentService struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
	kafka      kafka.IKafkaRegistry
	midtrans   clients.IMidtransClient
}

// Create implements IPaymentService.
// Fungsi ini membuat entitas pembayaran baru berdasarkan data permintaan dari client.
// Langkah-langkah utama:
// 1. Validasi apakah tanggal kedaluwarsa pembayaran valid.
// 2. Membuat payment link menggunakan Midtrans.
// 3. Menyimpan data pembayaran dan histori pembayaran ke database dalam satu transaksi.
// Jika berhasil, akan mengembalikan objek PaymentResponse; jika gagal, akan mengembalikan error.
func (p *PaymentService) Create(ctx context.Context, req *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	var (
		txErr, err error
		payment    *models.Payment
		response   *dto.PaymentResponse
		midtrans   *clients.MidtransData
	)

	// Eksekusi dalam konteks transaksi DB
	err = p.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		// Validasi bahwa waktu kedaluwarsa lebih dari waktu sekarang
		if !req.ExpiredAt.After(time.Now()) {
			return errPayment.ErrExpireArInvalid
		}

		// Membuat payment link dari Midtrans
		midtrans, txErr = p.midtrans.CreatePaymentLink(req)
		if txErr != nil {
			return txErr
		}

		// Persiapan data request untuk disimpan ke database
		paymentRequest := &dto.PaymentRequest{
			OrderID:     req.OrderID,
			Amount:      req.Amount,
			ExpiredAt:   req.ExpiredAt,
			Description: req.Description,
			PaymentLink: midtrans.RedirectURL,
		}

		// Menyimpan data pembayaran ke database
		payment, txErr = p.repository.GetPayment().Create(ctx, tx, paymentRequest)
		if txErr != nil {
			return txErr
		}

		// Menyimpan histori pembayaran pertama kali
		txErr = p.repository.GetPaymentHistory().Create(ctx, tx, &dto.PaymentHistoryRequest{
			PaymentID: payment.ID,
			Status:    payment.Status.GetStatusString(),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Menyusun response untuk dikembalikan ke client
	response = &dto.PaymentResponse{
		UUID:        payment.UUID,
		OrderID:     payment.OrderID,
		Amount:      payment.Amount,
		Status:      payment.Status.GetStatusString(),
		PaymentLink: payment.PaymentLink,
		Description: payment.Description,
	}

	return response, nil
}

// GetAllWithPagination implements IPaymentService.
// Fungsi ini mengambil semua data pembayaran dari database dengan fitur pagination.
// Data yang diambil kemudian dikonversi ke bentuk DTO PaymentResponse.
// Hasilnya dibungkus dalam objek PaginationResult agar mendukung pagination di sisi client.
func (p *PaymentService) GetAllWithPagination(ctx context.Context, req *dto.PaymentRequestParam) (*util.PaginationResult, error) {
	// Ambil data pembayaran dan total count dari repository
	payments, total, err := p.repository.GetPayment().FindAllWithPagination(ctx, req)
	if err != nil {
		return nil, err
	}

	// Ubah data model menjadi slice DTO PaymentResponse
	paymentResults := make([]dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		paymentResults = append(paymentResults, dto.PaymentResponse{
			UUID:          payment.UUID,
			TransactionID: payment.TransactionID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			Status:        payment.Status.GetStatusString(),
			PaymentLink:   payment.PaymentLink,
			InvoiceLink:   payment.InvoiceLink,
			VANumber:      payment.VANumber,
			Bank:          payment.Bank,
			Description:   payment.Description,
			ExpiredAt:     payment.ExpiredAt,
			CreatedAt:     payment.CreatedAt,
			UpdatedAt:     payment.UpdatedAt,
		})
	}

	// Buat parameter pagination untuk hasil akhir
	paginationParam := util.PaginationParam{
		Page:  req.Page,
		Limit: req.Limit,
		Count: total,
		Data:  paymentResults,
	}

	// Generate response pagination
	response := util.GeneratePagination(paginationParam)
	return &response, nil
}

// GetByUUID implements IPaymentService.
func (p *PaymentService) GetByUUID(ctx context.Context, uuid string) (*dto.PaymentResponse, error) {
	payment, err := p.repository.GetPayment().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		UUID:          payment.UUID,
		TransactionID: payment.TransactionID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Status:        payment.Status.GetStatusString(),
		PaymentLink:   payment.PaymentLink,
		InvoiceLink:   payment.InvoiceLink,
		VANumber:      payment.VANumber,
		Bank:          payment.Bank,
		Description:   payment.Description,
		ExpiredAt:     payment.ExpiredAt,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}, nil
}

// Fungsi untuk mengkonversi nama bulan dari bahasa Inggris ke Indonesia
func (p *PaymentService) convertToIndonesianMonth(englishMonth string) string {
	// Peta (map) yang berisi mapping nama bulan Inggris-Indonesia
	months := map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"July":      "Juli",
		"August":    "Agustus",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"December":  "Desember",
	}

	// Mencari nama bulan Indonesia berdasarkan input bulan Inggris
	indonesianMonth, ok := months[englishMonth]
	if !ok {
		// Jika bulan tidak ditemukan, kembalikan error (dikonversi ke string)
		return errors.New("month not found").Error()
	}
	return indonesianMonth // Mengembalikan nama bulan dalam bahasa Indonesia
}

// Fungsi untuk menghasilkan PDF dari template HTML
func (p *PaymentService) generatePDF(req *dto.InvoiceRequest) ([]byte, error) {
	// Path ke template HTML invoice
	htmlTemplatePath := "templates/invoice.html"

	// Membaca file template HTML
	htmlTemplate, err := os.ReadFile(htmlTemplatePath)
	if err != nil {
		return nil, err
	}

	// Mempersiapkan data untuk template
	var data map[string]interface{}
	jsonData, _ := json.Marshal(req)      // Convert request ke JSON
	err = json.Unmarshal(jsonData, &data) // Convert JSON ke map
	if err != nil {
		return nil, err
	}

	// Memanggil fungsi utility untuk generate PDF dari HTML dan data
	pdf, err := util.GeneratePDFFromHTML(string(htmlTemplate), data)
	if err != nil {
		return nil, err
	}

	return pdf, nil
}

// Fungsi untuk mengupload file PDF ke Google Cloud Storage
func (p *PaymentService) uploadToGCS(ctx context.Context, invoiceNumber string, pdf []byte) (string, error) {
	// Membuat nama file yang aman untuk GCS
	// Mengganti karakter '/' dengan '-' dan mengubah ke lowercase
	invoiceNumberReplace := strings.ToLower(strings.ReplaceAll(invoiceNumber, "/", "-"))
	filename := fmt.Sprintf("%s.pdf", invoiceNumberReplace)

	// Mengupload file ke GCS dan mendapatkan URL
	url, err := p.gcs.UploadFile(ctx, filename, pdf)
	if err != nil {
		return "", err
	}

	return url, nil
}

// Fungsi untuk menghasilkan nomor acak 6 digit
func (p *PaymentService) randomNumber() int {
	// Membuat generator random dengan seed berdasarkan waktu sekarang
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Menghasilkan angka random antara 100000-999999
	number := random.Intn(900000) + 100000
	return number
}

// Fungsi untuk memetakan status transaksi ke event Kafka
func (p *PaymentService) mapTransactionStatusToEvent(status constants.PaymentStatusString) string {
	var paymentStatus string

	// Switch case untuk menentukan status pembayaran
	switch status {
	case constants.PendingString:
		paymentStatus = strings.ToUpper(constants.PendingString.String())
	case constants.SettlementString:
		paymentStatus = strings.ToUpper(constants.SettlementString.String())
	case constants.ExpireString:
		paymentStatus = strings.ToUpper(constants.ExpireString.String())
	}
	return paymentStatus
}

// Fungsi untuk memproduksi message ke Kafka
func (p *PaymentService) produceToKafka(req *dto.Webhook, payment *models.Payment, paidAt *time.Time) error {
	// Membuat struktur event Kafka
	event := dto.KafkaEvent{
		Name: p.mapTransactionStatusToEvent(req.TransactionStatus),
	}

	// Metadata Kafka message
	metadata := dto.KafkaMetaData{
		Sender:    "payment-service",               // Service pengirim
		SendingAt: time.Now().Format(time.RFC3339), // Waktu pengiriman
	}

	// Body Kafka message
	body := dto.KafkaBody{
		Type: "JSON",
		Data: &dto.KafkaData{
			OrderID:   req.OrderID,                    // ID order
			PaymentID: payment.UUID,                   // ID pembayaran
			Status:    req.TransactionStatus.String(), // Status transaksi
			PaidAt:    paidAt,                         // Waktu pembayaran
			ExpiredAt: *payment.ExpiredAt,             // Waktu kadaluarsa
		},
	}

	// Membuat struktur message Kafka lengkap
	kafkaMessage := dto.KafkaMessage{
		Event:    event,
		Metadata: metadata,
		Body:     body,
	}

	// Mendapatkan topic Kafka dari config
	topic := config2.Config.Kafka.Topic
	// Mengconvert message ke JSON
	kafkaMessageJSON, _ := json.Marshal(kafkaMessage)

	// Memproduksi message ke Kafka
	err := p.kafka.GetKafkaProducer().ProduceMessage(topic, kafkaMessageJSON)
	if err != nil {
		return err
	}

	return nil
}

// Implementasi WebHook untuk menangani callback dari payment gateway
func (p *PaymentService) WebHook(ctx context.Context, req *dto.Webhook) error {
	var (
		txErr, err         error
		paymentAfterUpdate *models.Payment
		paidAt             *time.Time
		invoiceLink        string
		pdf                []byte
	)

	// Memulai transaksi database
	err = p.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		// Mencari data pembayaran berdasarkan order ID
		_, txErr = p.repository.GetPayment().FindByOrderID(ctx, req.OrderID.String())
		if err != nil {
			return txErr
		}

		// Jika status settlement, set waktu pembayaran ke waktu sekarang
		if req.TransactionStatus == constants.SettlementString {
			now := time.Now()
			paidAt = &now
		}

		// Konversi status string ke integer
		status := req.TransactionStatus.GetStatusInt()
		vaNumber := req.VANumbers[0].VaNumber // Nomor VA
		bank := req.VANumbers[0].Bank         // Bank

		// Update data pembayaran
		_, txErr = p.repository.GetPayment().Update(ctx, tx, req.OrderID.String(), &dto.UpdatePaymentRequest{
			TransactionID: &req.TransactionID,
			Status:        &status,
			PaidAt:        paidAt,
			VANumber:      &vaNumber,
			Bank:          &bank,
			Acquirer:      req.Acquirer,
		})

		if txErr != nil {
			return txErr
		}

		// Mendapatkan data pembayaran setelah diupdate
		paymentAfterUpdate, txErr = p.repository.GetPayment().FindByOrderID(ctx, req.OrderID.String())
		if txErr != nil {
			return txErr
		}

		// Membuat history pembayaran
		txErr = p.repository.GetPaymentHistory().Create(ctx, tx, &dto.PaymentHistoryRequest{
			PaymentID: paymentAfterUpdate.ID,
			Status:    paymentAfterUpdate.Status.GetStatusString(),
		})

		// Jika status settlement, generate invoice
		if req.TransactionStatus == constants.SettlementString {
			// Format tanggal pembayaran
			paidDay := paidAt.Format("02")
			paidMonth := p.convertToIndonesianMonth(paidAt.Format("January"))
			paidYear := paidAt.Format("2006")

			// Generate nomor invoice acak
			invoiceNumber := fmt.Sprintf("INV/%s/ORD/%d", time.Now().Format(time.DateOnly), p.randomNumber())

			// Format jumlah pembayaran ke format Rupiah
			total := util.RupiahFormat(&paymentAfterUpdate.Amount)

			// Membuat request invoice
			invoiceRequest := &dto.InvoiceRequest{
				InvoiceNumber: invoiceNumber,
				Data: dto.InvoiceData{
					PaymentDetail: dto.InvoicePaymentDetail{
						PaymentMethod: req.PaymentType,
						BankName:      strings.ToUpper(*paymentAfterUpdate.Bank),
						VANumber:      *paymentAfterUpdate.VANumber,
						Date:          fmt.Sprintf("%s %s %s", paidDay, paidMonth, paidYear),
						IsPaid:        true,
					},
					Items: []dto.InvoiceItem{
						{
							Description: *paymentAfterUpdate.Description,
							Price:       total,
						},
					},
					Total: total,
				},
			}

			// Generate PDF invoice
			pdf, txErr = p.generatePDF(invoiceRequest)
			if txErr != nil {
				return txErr
			}

			// Upload invoice ke GCS
			invoiceLink, txErr = p.uploadToGCS(ctx, invoiceNumber, pdf)
			if txErr != nil {
				return txErr
			}

			// Update link invoice di database
			_, txErr = p.repository.GetPayment().Update(ctx, tx, req.OrderID.String(), &dto.UpdatePaymentRequest{
				InvoiceLink: &invoiceLink,
			})
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Memproduksi message ke Kafka
	err = p.produceToKafka(req, paymentAfterUpdate, paidAt)
	if err != nil {
		return fmt.Errorf("failed to produce message to kafka: %w", err)
	}

	return nil
}

// Constructor untuk PaymentService
func NewPaymentService(repository repositories.IRepositoryRegistry, gcs gcs.IGCSClient, kafka kafka.IKafkaRegistry, midtrans clients.IMidtransClient) IPaymentService {
	return &PaymentService{
		repository: repository, // Dependency repository
		gcs:        gcs,        // Dependency Google Cloud Storage client
		kafka:      kafka,      // Dependency Kafka
		midtrans:   midtrans,   // Dependency Midtrans client
	}
}
