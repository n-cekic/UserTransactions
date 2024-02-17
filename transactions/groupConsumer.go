package transactions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	L "userTransactions/logging"

	"github.com/IBM/sarama"
)

func (srv *Service) kafkaSetup(broker, groupAAAA, topic string) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	groupConsumer, err := sarama.NewConsumerGroup(strings.Split(broker, ","), groupAAAA, config)
	if err != nil {
		fmt.Println("Error creating Kafka consumer group:", err)
		return
	}

	handler := &ConsumerGroupHandler{repo: &srv.repo}

	go func() {
		for {
			err := groupConsumer.Consume(context.TODO(), []string{topic}, handler)
			if err != nil {
				fmt.Println("Error in Kafka consumer group:", err)
			}
		}
	}()
}

type ConsumerGroupHandler struct {
	repo *Repo
}

func (h *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		L.Logger.Infof("Received message from partition %d at offset %d:\n", msg.Partition, msg.Offset)
		L.Logger.Infof("Key: %s", string(msg.Key))
		L.Logger.Infof("Value: %s", string(msg.Value))
		var val kafkaMessage
		err := json.Unmarshal(msg.Value, &val)
		if err != nil {
			L.Logger.Error(err)
		} else {
			err := h.repo.insertNewUser(val.Id, val.CreatedAt)
			if err != nil {
				L.Logger.Error(err)
			} else {
				L.Logger.Info("New user added to the database")
			}
		}
		session.MarkMessage(msg, "")
	}

	return nil
}
