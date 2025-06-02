package dto

import (
	"time"

	"github.com/google/uuid"
)

// KafkaEvent merepresentasikan metadata dasar dari event Kafka.
type KafkaEvent struct {
	Name string `json:"name"` // Nama event Kafka (contoh: "PaymentStatusChanged")
}

// KafkaMetaData menyimpan informasi pengirim dan waktu pengiriman pesan Kafka.
type KafkaMetaData struct {
	Sender    string `json:"sender"`    // Nama atau ID pengirim pesan Kafka
	SendingAt string `json:"sendingAt"` // Waktu pengiriman dalam format string ISO8601
}

// KafkaData berisi data utama yang dikirim melalui Kafka.
// Digunakan untuk menyampaikan status pembayaran atau data transaksi.
type KafkaData struct {
	OrderID   uuid.UUID  `json:"orderID"`   // UUID dari pesanan
	PaymentID uuid.UUID  `json:"paymentID"` // UUID dari pembayaran
	Status    string     `json:"status"`    // Status pembayaran (e.g., "PAID", "EXPIRED")
	ExpiredAt time.Time  `json:"expiredAt"` // Waktu kadaluarsa pembayaran
	PaidAt    *time.Time `json:"paidAt"`    // Waktu pembayaran dilakukan (bisa null)
}

// KafkaBody membungkus tipe data dan data payload yang dikirim dalam pesan Kafka.
type KafkaBody struct {
	Type string     `json:"type"` // Jenis payload data (contoh: "payment")
	Data *KafkaData `json:"data"` // Pointer ke data payload utama
}

// KafkaMessage merupakan struktur lengkap pesan Kafka yang terdiri dari event, metadata, dan payload.
type KafkaMessage struct {
	Event    KafkaEvent    `json:"event"`    // Informasi event Kafka
	Metadata KafkaMetaData `json:"metadata"` // Metadata tambahan seperti pengirim dan timestamp
	Body     KafkaBody     `json:"body"`     // Payload utama berisi data pembayaran/transaksi
}
