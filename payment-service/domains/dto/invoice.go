package dto

type InvoiceRequest struct {
	InvoiceNumber string      `json:"invoiceNumber"` // Nomor faktur unik
	Data          InvoiceData `json:"data"`          // Data terkait faktur
}

type InvoiceData struct {
	PaymentDetail InvoicePaymentDetail `json:"paymentDetail"` // Detail pembayaran yang terkait dengan faktur
	Items         []InvoiceItem        `json:"items"`         // Daftar item yang termasuk dalam faktur
	Total         string               `json:"total"`         // Total jumlah yang harus dibayar
}

type InvoicePaymentDetail struct {
	BankName      string `json:"bankName"`      // Nama bank yang digunakan untuk pembayaran
	PaymentMethod string `json:"paymentMethod"` // Metode pembayaran yang digunakan (misalnya: "bank transfer", "credit card")
	VANumber      string `json:"vaNumber"`      // Nomor Virtual Account (jika menggunakan bank transfer)
	Date          string `json:"date"`          // Tanggal pembayaran yang dilakukan
	IsPaid        bool   `json:"isPaid"`        // Status apakah pembayaran sudah dilakukan atau belum
}

type InvoiceItem struct {
	Description string `json:"description"` // Deskripsi item yang termasuk dalam faktur
	Price       string `json:"price"`       // Harga per item
}
