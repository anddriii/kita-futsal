package config

import (
	"os"

	"github.com/anddriii/kita-futsal/user-service/common/util"
	"github.com/sirupsen/logrus"
	_ "github.com/spf13/viper/remote"
)

var Config AppConfig

type AppConfig struct {
	Port                  int      `json:"port"`
	AppName               string   `json:"appName"`
	AppEnv                string   `json:"appEnv"`
	SignatureKey          string   `json:"signatureKey"`
	Database              database `json:"database"`
	RateLimitMaxRequest   float64  `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond int      `json:"rateLimiterTimeSecond"`
	JwtSecretKey          string   `json:"jwtSecretKey"`
	JwtExpirationTime     int      `json:"jwtExpirationTime"`
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
