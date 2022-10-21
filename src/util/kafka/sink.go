package kafka

import (
	"github.com/Shopify/sarama"
	"log"
	//"github.com/reugn/go-streams/flow"
	//"github.com/reugn/go-streams/util"
)

// GroupHandler represents a Sarama consumer group handler
type GroupHandler struct {
	ready chan struct{}
	out   chan interface{}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (handler *GroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(handler.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (handler *GroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (handler *GroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message != nil {
				log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s\n",
					string(message.Value), message.Timestamp, message.Topic)
				session.MarkMessage(message, "")
				handler.out <- message
			}
		case <-session.Context().Done():
			return session.Context().Err()
		}
	}
}

// KafkaSink connector
type KafkaSink struct {
	producer sarama.SyncProducer
	config   *sarama.Config
	topic    string
	in       chan interface{}
}

// NewKafkaSink returns a new KafkaSink instance
func NewKafkaSink(addrs []string, config *sarama.Config, topic string) *KafkaSink {
	producer, err := sarama.NewSyncProducer(addrs, config)
	//util.Check(err)
	if err != nil {
		panic(err)
	}
	sink := &KafkaSink{
		producer: producer,
		topic:    topic,
		in:       make(chan interface{}),
		config:   config,
	}
	go sink.init()
	return sink
}

// init starts the main loop
func (ks *KafkaSink) init() {
	for msg := range ks.in {
		switch m := msg.(type) {
		case *sarama.ProducerMessage:
			ks.producer.SendMessage(m)
		case *sarama.ConsumerMessage:
			sMsg := &sarama.ProducerMessage{
				Topic: ks.topic,
				Key:   sarama.StringEncoder(m.Key),
				Value: sarama.StringEncoder(m.Value),
			}
			ks.producer.SendMessage(sMsg)
		case string:
			sMsg := &sarama.ProducerMessage{
				Topic: ks.topic,
				Value: sarama.StringEncoder(m),
			}
			ks.producer.SendMessage(sMsg)
		default:
			log.Printf("Unsupported message type %v\n", m)
		}
	}
	log.Printf("Closing the Kafka producer\n")
	ks.producer.Close()
}

// In returns an input channel for receiving data
func (ks *KafkaSink) In() chan<- interface{} {
	return ks.in
}
