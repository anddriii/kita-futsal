package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anddriii/kita-futsal/field-service/clients"
	"github.com/anddriii/kita-futsal/field-service/common/gcs"
	"github.com/anddriii/kita-futsal/field-service/common/response"
	"github.com/anddriii/kita-futsal/field-service/config"
	"github.com/anddriii/kita-futsal/field-service/constants"
	"github.com/anddriii/kita-futsal/field-service/controllers"
	"github.com/anddriii/kita-futsal/field-service/domains/models"
	"github.com/anddriii/kita-futsal/field-service/middlewares"
	"github.com/anddriii/kita-futsal/field-service/repositories"
	"github.com/anddriii/kita-futsal/field-service/routes"
	"github.com/anddriii/kita-futsal/field-service/services"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// command adalah objek Cobra untuk menjalankan perintah "serve"
var command = cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		// Memuat variabel lingkungan dari file .env
		_ = godotenv.Load()
		config.Init()
		fmt.Printf("Config setelah Init: %+v\n", config.Config)

		// Inisialisasi koneksi database
		db, err := config.InitDB()
		if err != nil {
			log.Fatalf("error in config init db %s", err)
			panic(err)
		}

		// Mengatur zona waktu lokal ke "Asia/Jakarta"
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Fatalf("load location %s", err)
			panic(err)
		}
		time.Local = loc

		// Migrasi database untuk model Role dan User
		err = db.AutoMigrate(
			&models.Field{},
			&models.FieldSchedule{},
			&models.Time{},
		)
		if err != nil {
			log.Fatalf("error in migrate %s", err)
			panic(err)
		}

		// Inisialisasi klien GCS
		gcsClient := initGCS()
		client := clients.NewClientRegistry()

		// Inisialisasi repository, service, dan controller
		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository, gcsClient)
		controller := controllers.NewControllerRegistry(service)

		// Membuat instance router Gin
		router := gin.Default()

		// Middleware untuk menangani panic dan mengembalikan response yang sesuai
		// router.Use(middlewares.HandlePanic())

		// Handler untuk route yang tidak ditemukan
		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})

		// Endpoint utama
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, response.Response{
				Status:  constants.Succes,
				Message: "Welcome to user service",
			})
		})

		// Middleware untuk menangani CORS (Cross-Origin Resource Sharing)
		router.Use(func(ctx *gin.Context) {
			ctx.Writer.Header().Set("Acces-Control-Allow-Origin", "*")
			ctx.Writer.Header().Set("Acces-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Writer.Header().Set("Acces-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-api-key, x-request-at")
			ctx.Next()
		})

		// Middleware untuk membatasi jumlah permintaan (Rate Limiting)
		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			},
		)
		router.Use(middlewares.RateLimiter(lmt))

		// Inisialisasi route untuk API versi 1
		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controller, group, client)
		route.Serve()

		// Menjalankan server pada port yang telah dikonfigurasi
		port := fmt.Sprintf(":%d", config.Config.Port)
		router.Run(port)
	},
}

// Run menjalankan command "serve" untuk memulai server
func Run() {
	err := command.Execute()
	if err != nil {
		log.Fatalf("error run %s", err)
		panic(err)
	}
	log.Println("Server running on port 8001")
}

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
Documen for this code : https://chatgpt.com/canvas/shared/67b81457fdd88191900c9806170cc048
*/
