package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anddriii/kita-futsal/payment-service/clients"
	midtransClient "github.com/anddriii/kita-futsal/payment-service/clients/midtrans"
	"github.com/anddriii/kita-futsal/payment-service/common/gcs"
	"github.com/anddriii/kita-futsal/payment-service/common/response"
	"github.com/anddriii/kita-futsal/payment-service/config"
	"github.com/anddriii/kita-futsal/payment-service/constants"
	controllers "github.com/anddriii/kita-futsal/payment-service/controllers/http"
	kafkaClient "github.com/anddriii/kita-futsal/payment-service/controllers/kafka"
	"github.com/anddriii/kita-futsal/payment-service/domains/models"
	"github.com/anddriii/kita-futsal/payment-service/middlewares"
	"github.com/anddriii/kita-futsal/payment-service/repositories"
	"github.com/anddriii/kita-futsal/payment-service/routes"
	"github.com/anddriii/kita-futsal/payment-service/service"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// command adalah objek Cobra yang digunakan untuk menjalankan perintah "serve"
// Perintah ini akan menjalankan server payment-service.
var command = cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		// Memuat file .env ke dalam environment
		_ = godotenv.Load()
		config.Init()
		fmt.Printf("Config setelah Init: %+v\n", config.Config)

		// Inisialisasi koneksi database menggunakan konfigurasi
		db, err := config.InitDB()
		if err != nil {
			log.Fatalf("error in config init db %s", err)
			panic(err)
		}

		// Set zona waktu server ke Asia/Jakarta
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Fatalf("load location %s", err)
			panic(err)
		}
		time.Local = loc

		// Migrasi tabel-tabel yang dibutuhkan
		err = db.AutoMigrate(
			&models.Payment{},
			&models.PaymentHistory{},
		)
		if err != nil {
			log.Fatalf("error in migrate %s", err)
			panic(err)
		}

		// Inisialisasi klien Google Cloud Storage
		gcsClient := initGCS()

		// Inisialisasi Kafka dan Midtrans client
		kafka := kafkaClient.NewKafkaRegistry(config.Config.Kafka.Brokers)
		midtrans := midtransClient.NewMidtransClient(
			config.Config.Midtrans.ServerKey,
			config.Config.Midtrans.IsProduction,
		)

		// Inisialisasi client internal antar layanan
		client := clients.NewClientRegistry()

		// Inisialisasi layer repository, service, dan controller
		repository := repositories.NewRepositoryRegistry(db)
		service := service.NewServiceRegistry(repository, gcsClient, kafka, midtrans)
		controller := controllers.NewControllerRegistry(service)

		// Buat router Gin dan pasang middleware
		router := gin.Default()

		// Middleware untuk menangani panic secara global
		router.Use(middlewares.HandlePanic())

		// Middleware untuk mengembalikan response 404 jika path tidak ditemukan
		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})

		// Endpoint dasar untuk pengecekan apakah service berjalan
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, response.Response{
				Status:  constants.Succes,
				Message: "Welcome to payment service",
			})
		})

		// Middleware untuk mengatur CORS agar bisa menerima request dari domain lain
		router.Use(func(ctx *gin.Context) {
			ctx.Writer.Header().Set("Acces-Control-Allow-Origin", "*")
			ctx.Writer.Header().Set("Acces-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Writer.Header().Set("Acces-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-api-key, x-request-at")
			if ctx.Request.Method == "OPTIONS" {
				ctx.AbortWithStatus(204)
				return
			}
			ctx.Next()
		})

		// Middleware untuk membatasi jumlah request (rate limiting)
		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			},
		)
		router.Use(middlewares.RateLimiter(lmt))

		// Inisialisasi route group versi 1 (v1)
		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controller, group, client)
		route.Serve() // Daftarkan seluruh endpoint

		// Menjalankan server pada port yang dikonfigurasi di file .env
		port := fmt.Sprintf(":%d", config.Config.Port)
		router.Run(port)
	},
}

// Run digunakan untuk mengeksekusi command "serve" saat aplikasi dijalankan.
func Run() {
	err := command.Execute()
	if err != nil {
		log.Fatalf("error run %s", err)
		panic(err)
	}
	log.Println("Server running on port 8001")
}

// initGCS menginisialisasi Google Cloud Storage Client dengan private key dari konfigurasi
func initGCS() gcs.IGCSClient {
	decode, err := base64.StdEncoding.DecodeString(config.Config.GCSPrivateKey)
	if err != nil {
		log.Fatalf("error in initGCS %s", err)
		panic(err)
	}

	stringPrivateKey := string(decode)
	gcsServiceAccount := gcs.ServiceAccountKeyJson{
		Type:                    config.Config.GCSType,
		ProjectId:               config.Config.GCSProjectID,
		PrivateKeyId:            config.Config.GCSPrivateKeyID,
		PrivateKey:              stringPrivateKey,
		ClientEmail:             config.Config.GCSClientEmail,
		ClientId:                config.Config.GCSClientID,
		AuthURI:                 config.Config.GCSAuthURI,
		TokenURI:                config.Config.GCSTokenURI,
		AuthProviderX509CertUrl: config.Config.GCSAuthProviderX509CertURL,
		ClientX509CertUrl:       config.Config.GCSClientX509CertURL,
		UniverseDomain:          config.Config.GCSUniverseDomain,
	}

	gcsClient := gcs.NewGCSClient(gcsServiceAccount, config.Config.GCSBucketName)
	return gcsClient
}

/*
Dokumentasi Tambahan:
https://chatgpt.com/canvas/shared/67b81457fdd88191900c9806170cc048
*/
