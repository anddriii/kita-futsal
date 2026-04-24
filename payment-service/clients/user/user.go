package clients

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/anddriii/kita-futsal/payment-service/clients/config"
	"github.com/anddriii/kita-futsal/payment-service/common/util"
	config2 "github.com/anddriii/kita-futsal/payment-service/config"
	"github.com/anddriii/kita-futsal/payment-service/constants"
)

// UserClient adalah struct yang digunakan untuk melakukan komunikasi dengan User Service.
type UserClient struct {
	client config.IClientConfig
}

type IUserClient interface {
	GetUserByToken(context.Context) (*UserData, error)
}

func NewUserClient(client config.IClientConfig) IUserClient {
	return &UserClient{client: client}
}

func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		config2.Config.AppName,
		u.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)

	// --- BAGIAN YANG DIPERBAIKI ---
	tokenVal := ctx.Value(constants.Token)
	if tokenVal == nil {
		// Jangan panic, return error biasa aja biar gampang di-trace
		return nil, fmt.Errorf("token tidak ditemukan di context")
	}

	token, ok := tokenVal.(string)
	if !ok {
		return nil, fmt.Errorf("format token di context bukan string")
	}
	// ------------------------------

	bearerToken := fmt.Sprintf("Bearer %s", token)

	var response UserResponse
	request := u.client.Client().Clone().
		Set(constants.Authorization, bearerToken).
		Set(constants.XServiceName, config2.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseUrl()))

	resp, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
		// Tambahin cetak error biar kelihatan di terminal
		fmt.Println("🔴 ERROR HTTP CLIENT:", errs[0])
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user response: %s", response.Message)
	}

	return &response.Data, nil
}
