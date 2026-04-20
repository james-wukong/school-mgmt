package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/james-wukong/online-school-mgmt/internal/config"
	"github.com/segmentio/kafka-go"
)

func ProduceScheduleTask(schoolID, semesterID int64, reqVersion float64, excludeRooms bool) error {
	cfg := config.InitConfig()
	// 1. Create the topic if not exists
	if len(cfg.Kafka.Brokers) == 0 {
		return errors.New("no brokers configured")
	}
	if err := CreateTopicIfNotExists(cfg.Kafka.Topic, cfg.Kafka.Brokers[0]); err != nil {
		return err
	}

	// 2. Initialize writer
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Brokers...),
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	// 3. Write topic to brokers
	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Key: []byte(fmt.Sprintf("%d", semesterID)),
			Value: []byte(fmt.Sprintf(`{"school_id":%d, "semester_id":%d, "version": %.2f, "exclude_rooms": %t}`,
				schoolID, semesterID, reqVersion, excludeRooms)),
		},
	)
}

func CreateTopicIfNotExists(topic, broker string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Check if topic exists
	partitions, err := conn.ReadPartitions(topic)
	if err == nil && len(partitions) > 0 {
		return nil
	}

	// Create topic with specific config
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     3, // Allows 3 workers to work in parallel
		ReplicationFactor: 1, // Use 3 in a real production cluster
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return err
	}
	return nil
}
