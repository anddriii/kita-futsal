package config

import (
	"os"

	"github.com/anddriii/kita-futsal/payment-service/common/util"
	"github.com/sirupsen/logrus"
	_ "github.com/spf13/viper/remote"
)

var Config AppConfig

type AppConfig struct {
	Port                       int             `json:"port"`
	AppName                    string          `json:"appName"`
	AppEnv                     string          `json:"appEnv"`
	SignatureKey               string          `json:"signatureKey"`
	Database                   database        `json:"database"`
	RateLimiterMaxRequest      float64         `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond      int             `json:"rateLimiterTimeSecond"`
	InternalService            InternalService `json:"internalService"`
	GCSType                    string          `json:"gcsType"`
	GCSProjectID               string          `json:"gcsProjectID"`
	GCSPrivateKeyID            string          `json:"gcsPrivateKeyID"`
	GCSPrivateKey              string          `json:"gcsPrivateKey"`
	GCSClientEmail             string          `json:"gcsClientEmail"`
	GCSClientID                string          `json:"gcsClientID"`
	GCSAuthURI                 string          `json:"gcsAuthURI"`
	GCSTokenURI                string          `json:"gcsTokenURI"`
	GCSAuthProviderX509CertURL string          `json:"gcsAuthProviderX509CertURL"`
	GCSClientX509CertURL       string          `json:"gcsClientX509CertURL"`
	GCSUniverseDomain          string          `json:"gcsUniverseDomain"`
	GCSBucketName              string          `json:"gcsBucketName"`
	Kafka                      Kafka           `json:"kafka"`
	Midtrans                   Midtrans        `json:"midtrans"`
}

type database struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Name            string `json:"name"`
	UserName        string `json:"username"`
	Password        string `json:"password"`
	MaxOpenConn     int    `json:"maxOpenConnection"`
	MaxLifeTimeConn int    `json:"maxLifetimeConnection"`
	MaxIdleConn     int    `json:"maxIdleConnection"`
	MaxIdleTime     int    `json:"maxIdleTime"`
}

type InternalService struct {
	User User `json:"user"`
}

type User struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

type Kafka struct {
	Brokers  []string `json:"brokers"`
	Topic    string   `json:"topic"`
	TimeInMS int64    `json:"timeInMS"`
	MaxRetry int      `json:"maxRetry"`
}

type Midtrans struct {
	ServerKey    string `json:"serverKey"`
	ClienttKey   string `json:"clientKey"`
	IsProduction bool   `json:"isProduction"`
}

/*
jika config dari local maka akan mengambil dari file config.json.
Tetapi jika confignya berasal dari grpc maka akan menggunakan util "BindFromConsul"
*/
func Init() {
	err := util.BindFromJson(&Config, "config.json", ".")
	if err != nil {
		logrus.Infof("Failed to bind config: %v", err)
		err := util.BindFromConsul(Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_PATH"))
		if err != nil {
			panic(err)
		}
	}
}
