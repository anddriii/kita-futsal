package dto

import "github.com/anddriii/kita-futsal/payment-service/constants"

type PaymentHistoryRequest struct {
	PaymentID uint                          `json:"paymentID"` // ID unik untuk pembayaran yang ingin diambil riwayatnya
	Status    constants.PaymentStatusString `json:"status"`    // Status pembayaran yang ingin difilter (misalnya: "PAID", "EXPIRED")
}
