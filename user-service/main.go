package main

import (
	"github.com/anddriii/kita-futsal/user-service/cmd"
)

func main() {
	cmd.Run()
}

/*
json untuk testing di pc lokal
{
    "port": 8001,
    "appName": "user-service",
    "appEnv": "local",
    "signatureKey": "c80a3afdd1600288da374e229b4a2a1f",
    "database": {
        "host": "localhost",
        "port": 5432,
        "name": "kita-futsal-test",
        "username": "postgres",
        "password": "280904",
        "maxOpenConnection": 10,
        "maxLifetimeConnection": 10,
        "maxIdleConnection": 10,
        "maxIdleTime": 10
    },
    "rateLimiterMaxRequest": 1000,
    "rateLimiterTimeSecond": 60,
    "jwtSecretKey": "336a6c766b8044aac272d43794d4d36924bb9bd9",
    "jwtExpirationTime": 1440
}


json untuk docker
{
    "port": 8001,
    "appName": "user-service",
    "appEnv": "local",
    "signatureKey": "c80a3afdd1600288da374e229b4a2a1f",
    "database": {
        "host": "localhost",
        "port": 5432,
        "name": "user_service",
        "username": "root",
        "password": "ikuyoKita",
        "maxOpenConnection": 10,
        "maxLifetimeConnection": 10,
        "maxIdleConnection": 10,
        "maxIdleTime": 10
    },
    "rateLimiterMaxRequest": 1000,
    "rateLimiterTimeSecond": 60,
    "jwtSecretKey": "336a6c766b8044aac272d43794d4d36924bb9bd9",
    "jwtExpirationTime": 1440
}
*/
