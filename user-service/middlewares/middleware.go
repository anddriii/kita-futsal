package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/anddriii/kita-futsal/user-service/common/response"
	"github.com/anddriii/kita-futsal/user-service/config"
	"github.com/anddriii/kita-futsal/user-service/constants"
	errCons "github.com/anddriii/kita-futsal/user-service/constants/error"
	services "github.com/anddriii/kita-futsal/user-service/services/user"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func extraBearertoken(token string) string {
	splitToken := strings.Split(token, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
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
		return errCons.ErrUnauthorized
	}

	// Jika validasi sukses, kembalikan nil (tidak ada error)
	return nil
}

// validasi bearer token JWT
func validateBearerToken(ctx *gin.Context, token string) error {
	//Jika token tidak mengandung kata "Bearer", maka request akan langsung ditolak dengan Unauthorized error.
	if !strings.Contains(token, "Bearer") {
		return errCons.ErrUnauthorized
	}

	tokenString := extraBearertoken(token)
	if tokenString == "" {
		return errCons.ErrUnauthorized
	}

	claims := services.Claims{}
	tokenJwt, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errCons.ErrInvalidToken
		}

		jwtScret := []byte(config.Config.JwtSecretKey)
		return jwtScret, nil
	})

	if err != nil || !tokenJwt.Valid {
		return errCons.ErrUnauthorized
	}

	//Menyimpan informasi user dari token ke context request
	userLogin := ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), constants.UserLogin, claims.User))
	ctx.Request = userLogin
	ctx.Set(constants.Token, token) //Menyimpan token dalam context Gin
	return nil

}

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		token := ctx.GetHeader(constants.Authorization)
		//Jika tidak ada token, request langsung ditolak dengan Unauthorized error.
		if token == "" {
			responUnauthorized(ctx, errCons.ErrUnauthorized.Error())
			return
		}

		err = validateBearerToken(ctx, token)
		if err != nil {
			responUnauthorized(ctx, err.Error())
			return
		}

		err = validateApiKey(ctx)
		if err != nil {
			responUnauthorized(ctx, err.Error())
			return
		}
		ctx.Next()
	}
}
