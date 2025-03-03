package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/anddriii/kita-futsal/field-service/clients"
	"github.com/anddriii/kita-futsal/field-service/common/response"
	"github.com/anddriii/kita-futsal/field-service/config"
	"github.com/anddriii/kita-futsal/field-service/constants"
	errCons "github.com/anddriii/kita-futsal/field-service/constants/error"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				ctx.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errCons.ErrInternalServerError.Error(),
				})
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

// membatasi request akses
func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errCons.ErrToManyRequest.Error(),
			})
			ctx.Abort()
		}
		ctx.Next()
	}
}

func responUnauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	ctx.Abort()
}

func validateApiKey(ctx *gin.Context) error {
	// Ambil nilai API Key dan metadata dari request header
	apiKey := ctx.GetHeader(constants.XApiKey)           // API Key yang dikirim oleh client
	serviceName := ctx.GetHeader(constants.XServiceName) // Nama layanan yang melakukan request
	requestAt := ctx.GetHeader(constants.XRequestAt)     // Timestamp saat request dikirim
	log.Println("Apikey:", apiKey)
	log.Println("servicename:", serviceName)
	log.Println("request at:", requestAt)

	// Ambil Signature Key dari konfigurasi server
	signatureKey := config.Config.SignatureKey

	// Buat string validasi dengan format "serviceName:signatureKey:requestAt"
	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)

	// Buat hash SHA-256 dari validateKey
	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil)) // Konversi hash ke format string heksadesimal

	// Bandingkan API Key dari request dengan hasil hash yang dihasilkan
	if apiKey != resultHash {
		// Jika tidak cocok, kembalikan error Unauthorized
		log.Println("Apikey tidak sesuai dengan hash")
		log.Println("api key: ", apiKey)
		log.Println("resultHash: ", resultHash)
		return errCons.ErrUnauthorized
	}

	// Jika validasi sukses, kembalikan nil (tidak ada error)
	return nil
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func CheckRole(roles []string, client clients.IClientRegistry) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := client.GetUser().GetUserByToken(ctx.Request.Context())
		if err != nil {
			responUnauthorized(ctx, errCons.ErrUnauthorized.Error())
			return
		}

		if !contains(roles, user.Role) {
			responUnauthorized(ctx, errCons.ErrUnauthorized.Error())
			return
		}
		ctx.Next()
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			log.Println("error: token kosong")
			responUnauthorized(c, errCons.ErrUnauthorized.Error())
			return
		}

		err = validateApiKey(c)
		if err != nil {
			log.Println("error vaidate API KEY")
			responUnauthorized(c, err.Error())
			return
		}

		c.Next()
	}
}

func AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := validateApiKey(c)
		if err != nil {
			log.Println("error vaidate API KEY")
			responUnauthorized(c, err.Error())
			return
		}

		c.Next()
	}
}
