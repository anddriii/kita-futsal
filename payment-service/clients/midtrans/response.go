package clients

// MidtransResponse merepresentasikan struktur respons standar dari layanan Midtrans,
// yang berisi status kode HTTP, status pesan, dan data detail pembayaran.
type MidtransResponse struct {
	Code   int          `json:"code"`   // Kode status HTTP (misal: 200, 400, dll)
	Status string       `json:"status"` // Pesan status (misal: "success", "error", dll)
	Data   MidtransData `json:"data"`   // Data utama yang berisi token dan URL redirect
}

// MidtransData menyimpan data penting dari hasil permintaan transaksi Midtrans,
// termasuk token pembayaran dan URL redirect untuk menyelesaikan pembayaran.
type MidtransData struct {
	Token       string `json:"token"`        // Token transaksi Midtrans
	RedirectURL string `json:"redirect_url"` // URL untuk redirect ke halaman pembayaran Midtrans
}
