package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	cluster "github.com/bsm/sarama-cluster"
)

func main6363() {

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// init consumer
	brokers := []string{"10.99.1.151:19092"}
	topics := []string{"test"}
	consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				m := make(map[string]interface{})
				err := json.Unmarshal([]byte(msg.Value), &m)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(m["device"])
					data := m["data"]
					if v, ok := data.([]interface{})[0].(map[string]interface{}); ok {
						fmt.Println(ok, v["humidity"], v["time"])
					}
				}
				consumer.MarkOffset(msg, "") // mark message as processed

			}
		case <-signals:
			return
		}
	}
}
