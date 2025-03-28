package config

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

type ClientConfig struct {
	client       *gorequest.SuperAgent
	baseUrl      string
	signatureKey string
}

type IClientConfig interface {
	Client() *gorequest.SuperAgent
	BaseUrl() string
	SignatureKey() string
}

type Option func(*ClientConfig)

func NewClientConfig(options ...Option) IClientConfig {
	clientConfig := &ClientConfig{
		client: gorequest.New().
			Set("Content-Type", "application/json").
			Set("Accept", "application/json"),
	}
	for _, option := range options {
		option(clientConfig)
	}

	return clientConfig
}

func (c *ClientConfig) Client() *gorequest.SuperAgent {
	return c.client
}

func (c *ClientConfig) BaseUrl() string {
	fmt.Println(c.baseUrl)
	return c.baseUrl
}

func (c *ClientConfig) SignatureKey() string {
	return c.signatureKey
}

func WithBaseURL(baseURL string) Option {
	return func(cc *ClientConfig) {
		cc.baseUrl = baseURL
	}
}

func WithSignatureKey(signatureKey string) Option {
	return func(cc *ClientConfig) {
		cc.signatureKey = signatureKey
	}
}
