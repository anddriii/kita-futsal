package kafka

import (
	"github.com/IBM/sarama"
	configApp "github.com/anddriii/kita-futsal/payment-service/config"
	"github.com/sirupsen/logrus"
)

// Kafka adalah struct yang mewakili koneksi ke cluster Kafka.
// Menyimpan daftar alamat broker Kafka.
type Kafka struct {
	brokers []string // Daftar alamat broker Kafka (misal: ["localhost:9092"])
}

// IKafka adalah interface untuk Kafka client yang mendefinisikan
// fungsi utama yang harus diimplementasikan, yaitu pengiriman pesan.
type IKafka interface {
	ProduceMessage(topic string, data []byte) error
}

// NewKafkaProducer mengembalikan instance Kafka baru sebagai implementasi dari IKafka.
// Parameter:
//   - brokers: array string berisi alamat broker Kafka
func NewKafkaProducer(brokers []string) IKafka {
	return &Kafka{
		brokers: brokers,
	}
}

// ProduceMessage mengirim pesan ke Kafka ke topik tertentu.
// Parameter:
//   - topic: nama topik Kafka
//   - data: payload/message dalam bentuk byte array
//
// Return:
//   - error jika gagal mengirim pesan atau membuat producer
func (k *Kafka) ProduceMessage(topic string, data []byte) error {
	// Konfigurasi producer Kafka menggunakan sarama
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true                     // Aktifkan pelaporan keberhasilan
	config.Producer.RequiredAcks = sarama.WaitForAll            // Tunggu ack dari semua broker
	config.Producer.Retry.Max = configApp.Config.Kafka.MaxRetry // Ambil retry max dari konfigurasi aplikasi

	// Inisialisasi producer sinkron
	producer, err := sarama.NewSyncProducer(k.brokers, config)
	if err != nil {
		logrus.Errorf("Failed to create producer: %v", err)
		return err
	}

	// Pastikan producer ditutup setelah digunakan
	defer func(producer sarama.SyncProducer) {
		err = producer.Close()
		if err != nil {
			logrus.Errorf("Failed to close kafka: %v", err)
			return
		}
	}(producer)

	// Buat pesan Kafka
	message := &sarama.ProducerMessage{
		Topic:   topic,                    // Nama topik Kafka
		Headers: nil,                      // (Opsional) header bisa ditambahkan jika perlu
		Value:   sarama.ByteEncoder(data), // Payload sebagai byte encoder
	}

	// Kirim pesan ke Kafka
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		logrus.Errorf("Failed to produce to kafka: %v", err)
		return err
	}

	// Log info lokasi penyimpanan pesan di Kafka
	logrus.Infof("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}
