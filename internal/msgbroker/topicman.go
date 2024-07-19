package msgbroker

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
)
func SubscribeTopic(topic string) {
	Log.Info().Msgf("subscribing to topic %v", topic)


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
    // Set up a channel for handling Ctrl-C, etc
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    // Process messages
    run := true
    for run {
        select {
        case sig := <-sigchan:
            fmt.Printf("Caught signal %v: terminating\n", sig)
            run = false
        default:
            ev, err := c.ReadMessage(100 * time.Millisecond)
            if err != nil {
                // Errors are informational and automatically handled by the consumer
                continue
            }
            fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
                *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
        }
    }
}

func UnsubscribeTopic(topic string) {
	Log.Info().Msgf("Unsubscribing frim topic %v", topic)

	// Panic if any failure is encountered
	if false{
		Log.Fatal().Msg("Unsubscribing failed")
	}
	
}
