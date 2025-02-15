package util

import (
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// digunakan untuk membaca file JSON dari server cloud atau dari Consul (penyimpanan konfigurasi berbasis key-value).

// BindFromJson membaca konfigurasi dari file JSON dan mengikatnya ke struct tujuan.
func BindFromJson(dest any, filename, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("Failed to unmarshal: %v", err)
		return err
	}
	return nil
}

// SetEnvFromConsulKV mengatur environment variables berdasarkan konfigurasi dari Consul KV.
func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)
	if err != nil {
		logrus.Errorf("Failed to unmarshal: %v", err)
		return err
	}

	for k, v := range env {
		var (
			valeOf = reflect.ValueOf(v)
			val    string
		)

		switch valeOf.Kind() {
		case reflect.String:
			val = valeOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valeOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valeOf.Uint()))
		case reflect.Float32:
			val = strconv.Itoa(int(valeOf.Float()))
		case reflect.Float64:
			val = strconv.Itoa(int(valeOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valeOf.Bool())
		default:
			panic("unsupported type")
		}

		err = os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("Failed to set env: %v", err)
			return err
		}
	}
	return nil
}

// BindFromConsul membaca konfigurasi dari Consul dan mengikatnya/menyimpan ke struct tujuan(path).
func BindFromConsul(dest any, endpoint, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endpoint, path)
	if err != nil {
		logrus.Errorf("Failed to add remote provider: %v", err)
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("Failed to read config: %v", err)
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("Failed to unmarshal: %v", err)
		return err
	}

	err = SetEnvFromConsulKV(v)
	if err != nil {
		logrus.Errorf("Failed to read env from consul kv %v", err)
		return err
	}

	return nil
}
