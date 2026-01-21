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

type FieldClient struct {
	client config.IClientConfig
}

// GetFieldByUUID implements IFieldClient.
func (f *FieldClient) GetFieldByUUID(ctx context.Context, uuid uuid.UUID) (*FieldData, error) {
	unixTime := time.Now().Unix()
	generateApiKey := fmt.Sprintf("%s:%s:%d", configApp.Config.AppName, f.client.SignatureKey(), unixTime)
	apiKey := util.GenerateSHA256(generateApiKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)
	url := fmt.Sprintf("%s/api/v1/field/schedule/%s", f.client.BaseURL(), uuid)

	var response FieldResponse
	request := gorequest.New().
		Set(constants.Authorization, bearerToken).
		Set(constants.XServiceName, configApp.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(url)

	resp, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user response: %s", response.Message)
	}

	return &response.Data, nil
}

// UpdateStatus implements IFieldClient.
func (f *FieldClient) UpdateStatus(request *dto.UpdateFieldScheduleStatusRequest) error {
	unixTime := time.Now().Unix()
	generateApiKey := fmt.Sprintf("%s:%s:%d", configApp.Config.AppName, f.client.SignatureKey(), unixTime)
	apiKey := util.GenerateSHA256(generateApiKey)
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, bodyResp, errs := gorequest.New().
		Patch(fmt.Sprintf("%s/api/v1/schedule/status", f.client.BaseURL())).
		Set(constants.XServiceName, configApp.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Send(string(body)).
		End()

	if len(errs) > 0 {
		return errs[0]
	}

	var response FieldResponse
	if resp.StatusCode != http.StatusOK {
		err = json.Unmarshal([]byte(bodyResp), &response)
		if err != nil {
			return err
		}
		fieldError := fmt.Errorf("field response: %s", response.Message)
		return fieldError
	}

	err = json.Unmarshal([]byte(bodyResp), &response)
	if err != nil {
		return err
	}

	return nil
}

type IFieldClient interface {
	GetFieldByUUID(context.Context, uuid.UUID) (*FieldData, error)
	UpdateStatus(request *dto.UpdateFieldScheduleStatusRequest) error
}

func NewFieldClient(client config.IClientConfig) IFieldClient {
	return &FieldClient{client: client}
}
