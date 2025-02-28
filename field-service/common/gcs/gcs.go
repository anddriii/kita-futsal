package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

// ServiceAccountKeyJson merepresentasikan struktur kunci akun layanan Google Cloud Storage
// yang biasanya diperoleh dalam format JSON dari Google Cloud Console.
type ServiceAccountKeyJson struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universal_domain"`
}

// GCSClient adalah struktur yang merepresentasikan client Google Cloud Storage.
type GCSClient struct {
	ServiceAccountKeyJson ServiceAccountKeyJson // Data akun layanan Google Cloud
	BucketName            string                // Nama bucket tempat penyimpanan file
}

// IGCSClient adalah interface yang mendefinisikan metode yang harus diimplementasikan oleh GCSClient.
type IGCSClient interface {
	UploadFile(context.Context, string, []byte) (string, error)
}

// NewGCSClient membuat instance baru dari GCSClient dan mengembalikan objek yang sesuai dengan interface IGCSClient.
func NewGCSClient(serviceAccountKeyJson ServiceAccountKeyJson, bucketName string) IGCSClient {
	return &GCSClient{
		ServiceAccountKeyJson: serviceAccountKeyJson,
		BucketName:            bucketName,
	}
}

// CreateClient membuat client penyimpanan Google Cloud menggunakan kredensial JSON.
func (g *GCSClient) CreateClient(ctx context.Context) (*storage.Client, error) {
	// Encode service account key ke dalam JSON
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(g.ServiceAccountKeyJson)
	if err != nil {
		logrus.Errorf("failed to encode service account key json: %v", err)
		return nil, err
	}

	jsonByte := reqBodyBytes.Bytes()
	// Buat client baru dengan kredensial JSON
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(jsonByte))
	if err != nil {
		logrus.Errorf("failed to create client: %v", err)
		return nil, err
	}

	return client, nil
}

// UploadFile mengunggah file ke Google Cloud Storage dan mengembalikan URL akses publiknya.
func (g *GCSClient) UploadFile(ctx context.Context, fileName string, data []byte) (string, error) {
	var (
		contentType      = "application/octet-stream"
		timeoutInSeconds = 60
	)

	// Membuat client penyimpanan
	client, err := g.CreateClient(ctx)
	if err != nil {
		logrus.Errorf("Failed to upload file: %v", err)
		return "", err
	}

	defer func(client *storage.Client) {
		if err := client.Close(); err != nil {
			logrus.Errorf("Failed to close client: %v", err)
		}
	}(client)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	// Referensi ke bucket dan objek file
	bucket := client.Bucket(g.BucketName)
	object := bucket.Object(fileName)
	buffer := bytes.NewBuffer(data)

	// Membuat writer untuk menulis file ke GCS
	writer := object.NewWriter(ctx)
	//ketika upload 1x request tidak akan dipecah peacah
	writer.ChunkSize = 0

	_, err = io.Copy(writer, buffer)
	if err != nil {
		logrus.Errorf("Failed to copy: %v", err)
		return "", err
	}

	// Menutup writer setelah upload selesai
	err = writer.Close()
	if err != nil {
		logrus.Errorf("Failed to close writer: %v", err)
		return "", err
	}

	// Memperbarui metadata objek dengan tipe konten yang sesuai
	_, err = object.Update(ctx, storage.ObjectAttrsToUpdate{ContentType: contentType})
	if err != nil {
		logrus.Errorf("Failed to update metadata: %v", err)
		return "", err
	}

	// Mengembalikan URL file yang diunggah
	url := fmt.Sprintf("https://www.storage.googleapis.com/%s/%s", g.BucketName, fileName)
	return url, nil
}
