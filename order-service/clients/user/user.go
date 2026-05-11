package clients

import (
	"context"
	"fmt"
	"net/http"

	"time"

	"github.com/anddriii/kita-futsal/order-service/clients/config"
	"github.com/anddriii/kita-futsal/order-service/common/util"
	config2 "github.com/anddriii/kita-futsal/order-service/config"
	"github.com/anddriii/kita-futsal/order-service/constants"
	"github.com/google/uuid"
	"github.com/parnurzeal/gorequest"
)

type UserClient struct {
	client config.IClientConfig
}

type IUserClient interface {
	GetUserByToken(context.Context) (*UserData, error)
	GetUserByUUID(context.Context, uuid.UUID) (*UserData, error)
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
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var response UserResponse
	request := gorequest.New().
		Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL())). // WAJIB TULIS GET DULUAN DI ATAS
		Set(constants.Authorization, bearerToken).                   // BARU SET HEADER DI BAWAHNYA
		Set(constants.XServiceName, config2.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime))

	resp, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user response: %s", response.Message)
	}

	return &response.Data, nil
}

func (u *UserClient) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		config2.Config.AppName,
		u.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var response UserResponse
	request := gorequest.New().
		Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL())). // WAJIB TULIS GET DULUAN DI ATAS
		Set(constants.Authorization, bearerToken).                   // BARU SET HEADER DI BAWAHNYA
		Set(constants.XServiceName, config2.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime))

	resp, _, errs := request.EndStruct(&response)
	fmt.Println("DEBUG: Requesting user data by UUID:", uuid)
	fmt.Println("DEBUG: Response user:", resp)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user response: %s", response.Message)
	}

	fmt.Println("DEBUG: Received user data:", response.Data)

	return &response.Data, nil
}
