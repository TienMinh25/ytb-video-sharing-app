package third_party

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"ytb-video-sharing-app-be/pkg"
	"ytb-video-sharing-app-be/pkg/worker"

	kafkaconfluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type queue struct {
	groupID     string                   // group consumer id
	producer    *kafkaconfluent.Producer // producer
	consumer    *kafkaconfluent.Consumer // consumer
	subscribers []*pkg.SubscriptionInfo  // information of subscribers
	mu          sync.RWMutex             // concurrent lock when subcribe
	workerPool  *worker.Pool             // Worker pool for processing messages
}

func NewQueue() (pkg.Queue, error) {
	brokersString := os.Getenv("KAFKA_BROKERS")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	q := &queue{
		groupID:     groupID,
		producer:    newKafkaProducer(brokersString),
		consumer:    newKafkaConsumer(brokersString, groupID),
		subscribers: make([]*pkg.SubscriptionInfo, 0),
	}

	workerPool := worker.NewWorkerPool(5, 100, q.processKafkaMessage)
	q.workerPool = workerPool

	subInfo := &pkg.SubscriptionInfo{
		Topic: os.Getenv("KAFKA_TOPIC"),
	}

	if err := q.Subscribe(subInfo); err != nil {
		return nil, err
	}

	q.workerPool.Start()

	return q, nil
}

func (q *queue) processKafkaMessage(message interface{}) error {
	kafkaMsg, ok := message.(*kafkaconfluent.Message)
	if !ok {
		return fmt.Errorf("invalid message type")
	}

	fmt.Printf("Processing Kafka message: %s\n", string(kafkaMsg.Value))

	// data, err := DeserializeVideoMessageEvent(kafkaMsg.Value)

	// if err != nil {
	// 	return err
	// }

	// send message through websocket
	// q.websock.SendMessage(data.AccountId, &websock.EventMessage{
	// 	Title:     data.Title,
	// 	SharedBy:  data.SharedBy,
	// 	Thumbnail: data.Thumbnail,
	// })

	return nil
}

func newKafkaProducer(brokers string) *kafkaconfluent.Producer {
	retries, err := strconv.Atoi(os.Getenv("KAFKA_RETRY_ATTEMPTS"))

	if err != nil {
		// fallback value
		retries = 5
	}

	producerMaxWait, err := strconv.Atoi(os.Getenv("KAFKA_PRODUCER_MAX_WAIT"))

	if err != nil {
		// fallback value
		producerMaxWait = 300
	}

	p, err := kafkaconfluent.NewProducer(&kafkaconfluent.ConfigMap{
		"bootstrap.servers":                     brokers,
		"client.id":                             "myProducer",
		"acks":                                  "all",
		"enable.idempotence":                    true,
		"max.in.flight.requests.per.connection": 5,
		"transactional.id":                      uuid.New().String(),
		"retries":                               retries,
		"linger.ms":                             producerMaxWait,
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	err = p.InitTransactions(context.Background())
	if err != nil {
		fmt.Printf("Failed to initialize transactions: %s\n", err)
		os.Exit(1)
	}

	return p
}

func newKafkaConsumer(brokers, groupID string) *kafkaconfluent.Consumer {
	fetchMinBytes, err := strconv.Atoi(os.Getenv("KAFKA_CONSUMER_FETCH_MIN_BYTES"))

	if err != nil {
		// fallback value
		fetchMinBytes = 1
	}

	fetchMaxBytes, err := strconv.Atoi(os.Getenv("KAFKA_CONSUMER_FETCH_MAX_BYTES"))

	if err != nil {
		// fallback value
		fetchMaxBytes = 1e6
	}

	timeMaxWait, err := strconv.Atoi(os.Getenv("KAFKA_CONSUMER_MAX_WAIT"))

	if err != nil {
		// fallback value
		timeMaxWait = 10000
	}

	c, err := kafkaconfluent.NewConsumer(&kafkaconfluent.ConfigMap{
		"bootstrap.servers":  brokers,
		"client.id":          "my-consumer",
		"group.id":           groupID,
		"enable.auto.commit": false,
		"auto.offset.reset":  "latest",
		// Consumer Tuning
		"max.poll.interval.ms":  60000, // 1p (kafka will kick consumer)
		"heartbeat.interval.ms": 5000,  // 5s
		"session.timeout.ms":    45000,
		"fetch.min.bytes":       fetchMinBytes,
		"fetch.max.bytes":       fetchMaxBytes,
		"fetch.wait.max.ms":     timeMaxWait,
	})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	return c
}

// Close implements pkg.Queue.
func (q *queue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.workerPool.GracefulStop() // Graceful shutdown worker pool
	q.consumer.Close()
	q.producer.Close()
	return nil
}

// Publish implements pkg.Queue.
func (q *queue) Produce(topic string, payload []byte) error {
	err := q.producer.BeginTransaction()

	if err != nil {
		fmt.Printf("Failed to begin transaction: %s\n", err)
		return err
	}

	deliveryChan := make(chan kafkaconfluent.Event)

	err = q.producer.Produce(&kafkaconfluent.Message{
		TopicPartition: kafkaconfluent.TopicPartition{Partition: kafkaconfluent.PartitionAny, Topic: &topic},
		Value:          payload,
	}, deliveryChan)

	if err != nil {
		fmt.Printf("Failed to produce message: %s\n", err)
		_ = q.producer.AbortTransaction(context.Background())
		return err
	}

	e := <-deliveryChan
	m := e.(*kafkaconfluent.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		_ = q.producer.AbortTransaction(context.Background())
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

		err = q.producer.CommitTransaction(context.Background())
		if err != nil {
			fmt.Printf("Failed to commit transaction: %s\n", err)
			_ = q.producer.AbortTransaction(context.Background())
			return err
		}
	}

	close(deliveryChan)

	return nil
}

// Subscribe implements pkg.Queue.
func (q *queue) Subscribe(payload *pkg.SubscriptionInfo) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.subscribers = append(q.subscribers, payload)
	go q.consume(payload)
	return nil
}

func (q *queue) consume(sub *pkg.SubscriptionInfo) {
	if err := q.consumer.Subscribe(sub.Topic, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topics: %s\n", err)
	}

	for {
		event := q.consumer.Poll(100)

		if event == nil {
			continue
		}

		switch msg := event.(type) {
		case *kafkaconfluent.Message:
			fmt.Printf("Received message: %s\n", string(msg.Value))
			q.workerPool.PushMessage(msg)

			_, err := q.consumer.CommitMessage(msg)

			if err != nil {
				fmt.Printf("Failed to commit message: %s\n", err)
			}
		case kafkaconfluent.Error:
			fmt.Printf("Consumer error: %s\n", msg)
		}
	}
}

// Serialize (marshal) to Protobuf []byte
func SerializeVideoMessageEvent(event *VideoMessageEvent) ([]byte, error) {
	return proto.Marshal(event)
}

// Deserialize (unmarshal) Protobuf []byte to struct
func DeserializeVideoMessageEvent(data []byte) (*VideoMessageEvent, error) {
	var video VideoMessageEvent
	err := proto.Unmarshal(data, &video)
	return &video, err
}
