package util

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// digunakan untuk membaca file JSON dari server cloud atau dari Consul (penyimpanan konfigurasi berbasis key-value).

// BindFromJson membaca konfigurasi dari file JSON dan mengikatnya ke struct tujuan.
func BindFromJson(dest any, filename, path string) error {
	v := viper.New()

	v.SetConfigType("json")

	// Hapus ekstensi dari nama file agar sesuai dengan Viper
	configName := strings.TrimSuffix(filename, filepath.Ext(filename))
	v.SetConfigName(configName)

	// Cek apakah path kosong
	if path != "" {
		v.AddConfigPath(path)
	} else {
		v.AddConfigPath(".") // Pakai direktori saat ini jika path kosong
	}

	// Log untuk debug
	fmt.Printf("Loading config: %s.json from %s\n", configName, path)

	// Baca file konfigurasi
	if err := v.ReadInConfig(); err != nil {
		logrus.Errorf("Failed to read config file: %v", err)
		return err
	}

	// Unmarshal ke struct
	if err := v.Unmarshal(&dest); err != nil {
		logrus.Errorf("Failed to unmarshal config: %v", err)
		return err
	}

	fmt.Printf("Config loaded successfully: %+v\n", dest)
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
