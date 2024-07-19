package msgbroker

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
)

type ProcessorState struct {
	Msgs chan string
	Done chan bool
	Wg   sync.WaitGroup
}

func (ps *ProcessorState) Close() error {
	for range viper.GetInt("NUM_PROCESSORS") {
		ps.Done <- true
	}
	ps.Wg.Wait()
	return nil
}

type TopicMan struct {
	bookKeepingWriteMutex sync.Mutex
	topicBookKeeping      map[string]chan bool
	topicsSubscribed      atomic.Int64
	wg                    sync.WaitGroup
	processorFunc         func(string)
}

func (tm *TopicMan) SubscribeTopic(topic string) {
	_, exists := tm.topicBookKeeping[topic]
	if exists {
		Log.Warn().Msg("Topic already subscribed")
		return
	}

	tm.topicsSubscribed.Add(1)

	// Map in golang cannot concurrent writes, but reads can be concurrent
	tm.bookKeepingWriteMutex.Lock()
	tm.topicBookKeeping[topic] = make(chan bool, 1)
	tm.bookKeepingWriteMutex.Unlock()

	tm.wg.Add(1) // Add to waitgroup
	go func() {
		defer tm.wg.Done()

		Log.Info().Msgf("Subscribing to topic %v", topic)

		c, err := kafka.NewConsumer(&kafka.ConfigMap{
			// User-specific properties that you must set
			"bootstrap.servers": viper.GetString("KAFKA_BROKER"),
			/*"sasl.username":     "<CLUSTER API KEY>",
			  "sasl.password":     "<CLUSTER API SECRET>",

			  // Fixed properties
			  "security.protocol": "SASL_SSL",*/
			"sasl.mechanisms":   "PLAIN",
			"group.id":          "kafka-go-getting-started",
			"auto.offset.reset": "earliest"})

		if err != nil {
			Log.Fatal().Msgf("Failed to create consumer: %s", err)
		}

		err = c.SubscribeTopics([]string{topic}, nil)

		if err != nil {
			Log.Fatal().Msgf("Failed to subscribe to topic %s", err)
		}

		// Read messages
		run := true
		for run {
			select {
			case <-tm.topicBookKeeping[topic]:
				tm.bookKeepingWriteMutex.Lock()
				delete(tm.topicBookKeeping, topic)
				tm.bookKeepingWriteMutex.Unlock()

				Log.Info().Msg("End received. Stopping reading the topic!")
				run = false
			default:
				ev, err := c.ReadMessage(100 * time.Millisecond)
				if err != nil {
					// Errors are informational and automatically handled by the consumer
					continue
				}
				tm.processorFunc(string(ev.Value)) //process the message
				Log.Info().Msgf("Consumed event from topic %s: key = %-10s value = %s\n",
					*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			}
		}
	}()
}

func (tm *TopicMan) UnsubscribeTopic(topic string) {
	tm.topicsSubscribed.Add(-1)

	Log.Info().Msgf("Unsubscribing from topic %v", topic)
	tm.topicBookKeeping[topic] <- true

	// Panic if any failure is encountered
	if false {
		Log.Fatal().Msg("Unsubscribing failed")
	}

}

func InitBroker(processorFunc func(string)) (*TopicMan, error) {
	var tm = TopicMan{processorFunc: processorFunc, topicBookKeeping: make(map[string]chan bool)}
	return &tm, nil
}

func (tm *TopicMan) Close() error {
	Log.Info().Msg("Calling Close")
	for topicK := range tm.topicBookKeeping {
		Log.Info().Msg("sendong close")
		tm.topicBookKeeping[topicK] <- true // Send Done to Topic Readers
	}

	tm.wg.Wait()
	return nil
}
