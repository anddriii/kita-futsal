package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/anddriii/kita-futsal/field-service/clients/config"
	"github.com/anddriii/kita-futsal/field-service/common/util"
	config2 "github.com/anddriii/kita-futsal/field-service/config"
	"github.com/anddriii/kita-futsal/field-service/constants"
)

// UserClient adalah struct yang digunakan untuk melakukan komunikasi dengan User Service.
type UserClient struct {
	Client config.IClientConfig // Client konfigurasi yang digunakan untuk HTTP request.
}

// IUserClient adalah interface yang mendefinisikan metode untuk mendapatkan data user berdasarkan token.
type IUserClient interface {

	// GetUserByToken mengambil data user berdasarkan token dalam context.
	GetUserByToken(context.Context) (*UserData, error)
}

// NewUserClient mengembalikan instance baru dari UserClient.
func NewUserClient(client config.IClientConfig) IUserClient {
	return &UserClient{Client: client}
}

// GetUserByToken mengambil data user berdasarkan token yang disertakan dalam context.
func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
	// Mengambil timestamp saat ini
	unixTime := time.Now().Unix()

	// Membuat API key menggunakan SHA-256
	generateAPIKey := fmt.Sprintf("%s:%s:%d", config2.Config.AppName, u.Client.SignatureKey(), unixTime)
	apiKey := util.GenerateSHA256(generateAPIKey)

	// Mengambil token dari context
	tokenValue := ctx.Value(constants.Token)
	if tokenValue == nil {
		return nil, errors.New("token is missing from context")
	}

	token, ok := tokenValue.(string)
	if !ok {
		return nil, errors.New("token is not a valid string")
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)

	// Menyiapkan request ke user service
	var response UserRespone
	request := u.Client.Client().Clone().
		Set(constants.Authorization, bearerToken).
		Set(constants.XServiceName, config2.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/auth/user", u.Client.BaseUrl()))

	// Mengirim request dan menangani response
	resp, _, err := request.EndStruct(&response)
	if len(err) > 0 {
		return nil, err[0]
	}

	// Jika response tidak OK, kembalikan error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user response: %s", response.Message)
	}

	// Mengembalikan data user jika berhasil
	return &response.Data, nil
}
