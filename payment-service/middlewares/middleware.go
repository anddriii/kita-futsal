package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/anddriii/kita-futsal/payment-service/clients"
	"github.com/anddriii/kita-futsal/payment-service/common/response"
	"github.com/anddriii/kita-futsal/payment-service/config"
	"github.com/anddriii/kita-futsal/payment-service/constants"
	errCons "github.com/anddriii/kita-futsal/payment-service/constants/error"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HandlePanic menangani panic yang tidak terduga selama request lifecycle.
// Jika terjadi panic, akan mengembalikan response 500 dan mencatat error.
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

// RateLimiter membatasi jumlah request yang dikirim oleh client dalam periode waktu tertentu.
// Jika melebihi batas, mengembalikan response 429 Too Many Requests.
func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errCons.ErrToManyRequest.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// responUnauthorized mengembalikan response unauthorized (401) dan menghentikan request.
func responUnauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	ctx.Abort()
}

// validateApiKey memvalidasi API Key dari header dengan cara hashing signature key dan membandingkannya.
// Return error jika tidak valid, nil jika valid.
func validateApiKey(ctx *gin.Context) error {
	apiKey := ctx.GetHeader(constants.XApiKey)
	serviceName := ctx.GetHeader(constants.XServiceName)
	requestAt := ctx.GetHeader(constants.XRequestAt)

	log.Println("Apikey:", apiKey)
	log.Println("servicename:", serviceName)
	log.Println("request at:", requestAt)

	signatureKey := config.Config.SignatureKey
	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)

	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if apiKey != resultHash {
		log.Println("Apikey tidak sesuai dengan hash")
		log.Println("api key: ", apiKey)
		log.Println("resultHash: ", resultHash)
		return errCons.ErrUnauthorized
	}

	return nil
}

// contains mengecek apakah sebuah role terdapat dalam daftar role yang diperbolehkan.
func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// CheckRole adalah middleware untuk memverifikasi apakah user memiliki role yang diperbolehkan.
// Menggunakan client registry untuk mengambil data user berdasarkan token.
func CheckRole(roles []string, client clients.IClientRegistry) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := client.GetUser().GetUserByToken(ctx.Request.Context())
		// fmt.Println("user", user)
		if err != nil {
			fmt.Printf("error from get check role %s", err)
			return
		}

		if !contains(roles, user.Role) {
			fmt.Printf("unauthorized: user role %s not in %v\n", user.Role, roles)
			fmt.Println("role user:", user.Role)
			responUnauthorized(ctx, errCons.ErrUnauthorized.Error())
			return
		}
		ctx.Next()
	}
}

// extractBearerToken mengekstrak token dari header Authorization yang memiliki format "Bearer <token>".
func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}
	return ""
}

// Authenticate adalah middleware untuk memverifikasi token Authorization dan API key.
// Token disimpan dalam context untuk digunakan pada proses selanjutnya.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			log.Println("error: token kosong")
			responUnauthorized(c, errCons.ErrUnauthorized.Error())
			return
		}

		if err := validateApiKey(c); err != nil {
			log.Println("error vaidate API KEY")
			responUnauthorized(c, err.Error())
			return
		}

		tokenString := extractBearerToken(token)
		tokenUser := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.Token, tokenString))
		c.Request = tokenUser

		c.Next()
	}
}

// AuthenticateWithoutToken hanya memverifikasi API Key tanpa memerlukan Authorization token.
func AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateApiKey(c); err != nil {
			log.Println("error vaidate API KEY")
			responUnauthorized(c, err.Error())
			return
		}
		c.Next()
	}
}
