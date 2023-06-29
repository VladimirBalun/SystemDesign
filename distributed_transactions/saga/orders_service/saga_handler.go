package main

import (
	"encoding/json"

	"github.com/Shopify/sarama"
)

type SagaHandler struct {
	kafkaProducer sarama.AsyncProducer
	kafkaConsumer sarama.Consumer

	orderCreatedTopic  string
	goodsCreatedTopic  string
	goodsRejectedTopic string
}

func NewSagaHandler(brokers []string, oct, gct, grt string) (SagaHandler, error) {
	producer, err := sarama.NewAsyncProducer(brokers, sarama.NewConfig())
	if err != nil {
		return SagaHandler{}, err
	}

	consumer, err := sarama.NewConsumer(brokers, sarama.NewConfig())
	if err != nil {
		return SagaHandler{}, err
	}

	return SagaHandler{
		kafkaProducer:      producer,
		kafkaConsumer:      consumer,
		orderCreatedTopic:  oct,
		goodsCreatedTopic:  gct,
		goodsRejectedTopic: grt,
	}, nil
}

type SagaOrderCreatedEvent struct {
	OrderID int      `json:"order_id"`
	Goods   []string `json:"goods"`
}

func (s *SagaHandler) NotifyOrderCreated(orderID int, goods []string) error {
	var event SagaOrderCreatedEvent
	msgData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: s.orderCreatedTopic,
		Value: sarama.StringEncoder(msgData),
	}

	_, _, err := s.kafkaProducer.SendMessage(producerMsg)
	return err
}
