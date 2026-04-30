package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/anddriii/kita-futsal/order-service/config"
	"github.com/sirupsen/logrus"
)

type (
	TopicName string
	Handler   func(ctx context.Context, message *sarama.ConsumerMessage) error
)

type ConsumerGroup struct {
	handler map[TopicName]Handler
}

func NewConsumerGroup() *ConsumerGroup {
	return &ConsumerGroup{
		handler: make(map[TopicName]Handler),
	}
}

func (c *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	logrus.Info("Consumer group setup completed")
	return nil
}

func (c *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	logrus.Info("Consumer group cleanup completed")
	return nil
}

func (c *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	messages := claim.Messages()
	for message := range messages {
		handler, ok := c.handler[TopicName(message.Topic)]
		if !ok {
			logrus.Warnf("No handler found for topic: %s", message.Topic)
			continue
		}
		var err error
		maxRetry := config.Config.Kafka.MaxRetry
		for attempt := 1; attempt <= maxRetry; attempt++ {
			err = handler(context.Background(), message)
			if err == nil {
				break
			}
			logrus.Errorf("error handling message pn %s, attempt %d: %v", message.Topic, attempt, err)
			if attempt == maxRetry {
				logrus.Errorf("max retry reached for message on topic %s: %v", message.Topic, err)
			}
		}

		if err != nil {
			logrus.Errorf("error handling message on %s: %v", message.Topic, err)
			session.MarkMessage(message, err.Error())
			break
		}
		session.MarkMessage(message, time.Now().UTC().String())
	}

	return nil
}

func (c *ConsumerGroup) RegisterHandler(topic TopicName, handler Handler) {
	c.handler[topic] = handler
	logrus.Infof("register handler for topic %s", topic)
}
