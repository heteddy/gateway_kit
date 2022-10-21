/*
@Copyright:
*/
/*
@Time : 2021/2/11 19:29
@Author : teddy
@File : source.go
*/

package kafka

import (
	"context"
	"github.com/Shopify/sarama"

	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// KafkaSource connector
type KafkaSource struct {
	consumer  sarama.ConsumerGroup
	handler   sarama.ConsumerGroupHandler
	topics    []string
	config    *sarama.Config
	out       chan interface{}
	ctx       context.Context
	cancelCtx context.CancelFunc
	wg        *sync.WaitGroup
}

// NewKafkaSource returns a new KafkaSource instance
func NewKafkaSource(ctx context.Context, addrs []string, groupID string,
	config *sarama.Config, topics ...string) *KafkaSource {
	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {

	}
	//util.Check(err)
	out := make(chan interface{})
	_ctx, cancel := context.WithCancel(ctx)
	// 创建消费者
	source := &KafkaSource{
		consumer:  consumerGroup,
		handler:   &GroupHandler{make(chan struct{}), out},
		topics:    topics,
		config:    config,
		out:       out,
		ctx:       _ctx,
		cancelCtx: cancel,
		wg:        &sync.WaitGroup{},
	}

	go source.init()
	return source
}

func (ks *KafkaSource) claimLoop() {
	ks.wg.Add(1)
	defer func() {
		ks.wg.Done()
		log.Printf("Exiting the Kafka claimLoop\n")
	}()
	for {
		handler := ks.handler.(*GroupHandler)
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := ks.consumer.Consume(ks.ctx, ks.topics, handler); err != nil {
			log.Printf("Kafka consumer.Consume failed with: %v\n", err)
		}

		select {
		case <-ks.ctx.Done():
			return
		default:
		}

		handler.ready = make(chan struct{})
	}
}

// init starts the main loop
func (ks *KafkaSource) init() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go ks.claimLoop()

	select {
	case <-sigChan:
		ks.cancelCtx()
	case <-ks.ctx.Done():
	}

	log.Printf("Closing the Kafka consumer\n")
	ks.wg.Wait()
	close(ks.out)
	ks.consumer.Close()
}

//// Via streams data through the given flow
//func (ks *KafkaSource) Via(_flow streams.Flow) streams.Flow {
//	flow.DoStream(ks, _flow)
//	return _flow
//}

// Out returns an output channel for sending data
func (ks *KafkaSource) Out() <-chan interface{} {
	return ks.out
}
