package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anddriii/kita-futsal/order-service/clients/config"
	"github.com/anddriii/kita-futsal/order-service/common/util"
	configApp "github.com/anddriii/kita-futsal/order-service/config"
	"github.com/anddriii/kita-futsal/order-service/constants"
	"github.com/anddriii/kita-futsal/order-service/domain/dto"
	"github.com/google/uuid"
	"github.com/parnurzeal/gorequest"
)

type PaymentClient struct {
	client config.IClientConfig
}

// CreatePaymentLink implements IPaymentClient.
func (p *PaymentClient) CreatePaymentLink(ctx context.Context, req *dto.PaymentRequest) (*PaymentData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		configApp.Config.AppName,
		p.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, bodyResp, errs := gorequest.New().
		Post(fmt.Sprintf("%s/api/v1/payment", p.client.BaseURL())).
		Set(constants.Authorization, bearerToken).
		Set(constants.XServiceName, configApp.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Send(string(body)).
		End()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	var response PaymentResponse
	if resp.StatusCode != http.StatusCreated {
		err = json.Unmarshal([]byte(bodyResp), &response)
		if err != nil {
			return nil, err
		}
		paymentError := fmt.Errorf("payment response: %s", response.Message)
		return nil, paymentError
	}

	err = json.Unmarshal([]byte(bodyResp), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// GetPaymentByUUID implements IPaymentClient.
func (p *PaymentClient) GetPaymentByUUID(ctx context.Context, uuid uuid.UUID) (*PaymentData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		configApp.Config.AppName,
		p.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var response PaymentResponse
	request := gorequest.New().
		Set(constants.Authorization, bearerToken).
		Set(constants.XServiceName, configApp.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/payment/%s", p.client.BaseURL(), uuid))

	resp, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment response: %s", response.Message)
	}

	return &response.Data, nil
}

type IPaymentClient interface {
	GetPaymentByUUID(context.Context, uuid.UUID) (*PaymentData, error)
	CreatePaymentLink(context.Context, *dto.PaymentRequest) (*PaymentData, error)
}

func NewPaymentClient(client config.IClientConfig) IPaymentClient {
	return &PaymentClient{client: client}
}
