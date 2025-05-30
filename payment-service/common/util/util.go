package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// PaginationParam merepresentasikan parameter untuk paginasi.
type PaginationParam struct {
	Count int64 `json:"count"` // Total jumlah data
	Page  int   `json:"page"`  // Halaman saat ini
	Limit int   `json:"limit"` // Jumlah item per halaman
	Data  any   `json:"data"`  // Data yang akan dipaginasi
}

// PaginationResult merepresentasikan hasil dari proses paginasi.
type PaginationResult struct {
	TotalPage    int   `json:"totalPage"`    // Total jumlah halaman
	TotalData    int64 `json:"totalData"`    // Total jumlah data
	NextPage     *int  `json:"nextPage"`     // Halaman berikutnya (jika ada)
	PreviousPage *int  `json:"previousPage"` // Halaman sebelumnya (jika ada)
	Page         int   `json:"page"`         // Halaman saat ini
	Limit        int   `json:"limit"`        // Jumlah item per halaman
	Data         any   `json:"data"`         // Data yang dipaginasi
}

// GeneratePagination menghasilkan hasil paginasi berdasarkan parameter yang diberikan.
func GeneratePagination(params PaginationParam) PaginationResult {
	totalPage := int(math.Ceil(float64(params.Count) / float64(params.Limit)))

	var (
		nextPage     int
		previousPage int
	)

	// Menentukan halaman berikutnya jika masih ada
	if params.Page < totalPage {
		nextPage = params.Page + 1
	}

	// Menentukan halaman sebelumnya jika lebih dari 1
	if params.Page > 1 {
		previousPage = params.Page - 1
	}

	return PaginationResult{
		TotalPage:    totalPage,
		TotalData:    params.Count,
		NextPage:     &nextPage,
		PreviousPage: &previousPage,
		Page:         params.Page,
		Limit:        params.Limit,
		Data:         params.Data,
	}
}

// GenerateSHA256 menghasilkan hash SHA-256 dari string input.
func GenerateSHA256(inputString string) string {
	hash := sha256.New()
	hash.Write([]byte(inputString))
	hashBytes := hash.Sum(nil)

	// Mengembalikan hasil hash dalam format string hexadecimal
	return hex.EncodeToString(hashBytes)
}

// RupiahFormat mengonversi angka menjadi format mata uang Rupiah.
func RupiahFormat(amount *float64) string {
	stringValue := "0"
	if amount != nil {
		// Menggunakan humanize untuk menambahkan separator ribuan
		humanizeValue := humanize.CommafWithDigits(*amount, 0)
		stringValue = strings.ReplaceAll(humanizeValue, ",", ".")
	}
	return fmt.Sprintf("Rp. %s", stringValue)
}

// digunakan untuk membaca file JSON dari server cloud atau dari Consul (penyimpanan konfigurasi berbasis key-value).

// BindFromJson membaca konfigurasi dari file JSON dan mengikatnya ke struct tujuan.
func BindFromJson(dest any, filename, path string) error {
	v := viper.New()

	v.SetConfigType("json")

	// Hapus ekstensi dari nama file agar sesuai dengan Viper
	configName := strings.TrimSuffix(filename, filepath.Ext(filename))
	v.SetConfigName(configName)

	// Cek apakah path kosong
	if path != "" {
		v.AddConfigPath(path)
	} else {
		v.AddConfigPath(".") // Pakai direktori saat ini jika path kosong
	}

	// Log untuk debug
	fmt.Printf("Loading config: %s.json from %s\n", configName, path)

	// Baca file konfigurasi
	if err := v.ReadInConfig(); err != nil {
		logrus.Errorf("Failed to read config file: %v", err)
		return err
	}

	// Unmarshal ke struct
	if err := v.Unmarshal(&dest); err != nil {
		logrus.Errorf("Failed to unmarshal config: %v", err)
		return err
	}

	fmt.Printf("Config loaded successfully: %+v\n", dest)
	return nil
}

// SetEnvFromConsulKV mengatur environment variables berdasarkan konfigurasi dari Consul KV.
func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)
	if err != nil {
		logrus.Errorf("Failed to unmarshal: %v", err)
		return err
	}

	for k, v := range env {
		var (
			valeOf = reflect.ValueOf(v)
			val    string
		)

		switch valeOf.Kind() {
		case reflect.String:
			val = valeOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valeOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valeOf.Uint()))
		case reflect.Float32:
			val = strconv.Itoa(int(valeOf.Float()))
		case reflect.Float64:
			val = strconv.Itoa(int(valeOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valeOf.Bool())
		}

		err = os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("Failed to set env: %v", err)
			return err
		}
	}
	return nil
}

// BindFromConsul membaca konfigurasi dari Consul dan mengikatnya/menyimpan ke struct tujuan(path).
func BindFromConsul(dest any, endpoint, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endpoint, path)
	if err != nil {
		logrus.Errorf("Failed to add remote provider: %v", err)
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("Failed to read config: %v", err)
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("Failed to unmarshal: %v", err)
		return err
	}

	err = SetEnvFromConsulKV(v)
	if err != nil {
		logrus.Errorf("Failed to read env from consul kv %v", err)
		return err
	}

	return nil
}

// add1 adalah fungsi pembantu yang menerima satu argumen integer dan mengembalikan nilainya setelah ditambah 1.
// Fungsi ini bisa digunakan di dalam template HTML untuk melakukan operasi penambahan.
func add1(a int) int {
	return a + 1
}

// GeneratePDFFromHTML menerima string template HTML dan data apapun,
// lalu mengembalikan file PDF dalam bentuk byte slice ([]byte) atau error jika gagal.
// Fungsi ini memanfaatkan library wkhtmltopdf untuk mengubah HTML ke PDF.
func GeneratePDFFromHTML(htmlTemplate string, data any) ([]byte, error) {
	// Membuat funcMap untuk digunakan dalam template. Di sini kita menambahkan fungsi add1 agar bisa dipakai dalam HTML template.
	funcMap := template.FuncMap{
		"add1": add1, // Menambahkan fungsi add1 ke dalam template
	}

	// Parse template HTML dan mengikat funcMap ke dalamnya
	template, err := template.New("htmlTemplate").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return nil, err // Mengembalikan error jika parsing gagal
	}

	// Menyimpan hasil eksekusi template ke dalam buffer
	var filledTemplate bytes.Buffer
	if err := template.Execute(&filledTemplate, data); err != nil {
		return nil, err // Mengembalikan error jika eksekusi template gagal
	}
	htmlContent := filledTemplate.String() // Mengambil hasil HTML yang sudah diisi data

	// Membuat objek PDF generator menggunakan library wkhtmltopdf
	pdfGenerator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		logrus.Errorf("Failed to create PDF generator: %v", err)
		return nil, err // Mengembalikan error jika gagal membuat PDF generator
	}

	// Mengatur konfigurasi PDF
	pdfGenerator.Dpi.Set(600)                                     // Resolusi DPI
	pdfGenerator.NoCollate.Set(false)                             // Aktifkan penggabungan dokumen
	pdfGenerator.Orientation.Set(wkhtmltopdf.OrientationPortrait) // Orientasi potret
	pdfGenerator.PageSize.Set(wkhtmltopdf.PageSizeA4)             // Ukuran halaman A4
	pdfGenerator.Grayscale.Set(false)                             // Nonaktifkan mode grayscale

	// Menambahkan halaman dari HTML yang sudah dirender
	pdfGenerator.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlContent)))

	// Membuat PDF
	err = pdfGenerator.Create()
	if err != nil {
		logrus.Errorf("Failed to create PDF: %v", err)
		return nil, err // Mengembalikan error jika proses generate PDF gagal
	}

	// Mengembalikan PDF dalam bentuk byte slice
	return pdfGenerator.Bytes(), nil
}
