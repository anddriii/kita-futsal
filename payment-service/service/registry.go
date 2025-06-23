package service

import (
	clients "github.com/anddriii/kita-futsal/payment-service/clients/midtrans"
	"github.com/anddriii/kita-futsal/payment-service/common/gcs"
	"github.com/anddriii/kita-futsal/payment-service/controllers/kafka"
	"github.com/anddriii/kita-futsal/payment-service/repositories"
	services "github.com/anddriii/kita-futsal/payment-service/service/payment"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
	gcs        gcs.IGCSClient
	kafka      kafka.IKafkaRegistry
	midtrans   clients.IMidtransClient
}

type IServiceRegistry interface {
	GetPayment() services.IPaymentService
}

func NewServiceRegistry(
	repository repositories.IRepositoryRegistry,
	gcs gcs.IGCSClient,
	kafka kafka.IKafkaRegistry,
	midtrans clients.IMidtransClient,
) IServiceRegistry {
	return &Registry{
		repository: repository,
		gcs:        gcs,
		kafka:      kafka,
		midtrans:   midtrans,
	}
}

func (r *Registry) GetPayment() services.IPaymentService {
	return services.NewPaymentService(r.repository, r.gcs, r.kafka, r.midtrans)
}
