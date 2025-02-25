package pkg

type Queue interface {
	Subscriber
	Publisher

	// Closes the connection to the message broker, returning an error if the operation fails.
	Close() error
}

type Subscriber interface {
	// Subcribes to a topic or queue for receiving messages, returning an error if the subcription fails.
	Subscribe(payload *SubscriptionInfo) error
}

type Publisher interface {
	// Publishes a message to a specified topic, returning an error if the operation fails.
	Produce(topic string, payload []byte) error
}

type SubscriptionInfo struct {
	Topic string
}
