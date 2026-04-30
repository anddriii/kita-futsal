package kafka

import (
	kafka "github.com/anddriii/kita-futsal/order-service/controllers/kafka/payment"
	"github.com/anddriii/kita-futsal/order-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IKafkaRegistry interface {
	GetPayment() kafka.IPaymentKafka
}

func NewKafkaRegistry(service services.IServiceRegistry) IKafkaRegistry {
	return &Registry{service: service}
}

func (r *Registry) GetPayment() kafka.IPaymentKafka {
	return kafka.NewPaymentKafka(r.service)
}
