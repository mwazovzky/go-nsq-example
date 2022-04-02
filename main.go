package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type Message struct {
	Title     string
	Content   string
	Timestamp string
}

const topic = "test_topic"
const channel = "test_channel"
const nsqAddr = "nsqd:4150"
const nsqlookupdAddr = "nsqlookupd:4161"

func main() {
	// pub()
	sub()
}

func pub() {
	config := nsq.NewConfig()

	producer, err := nsq.NewProducer(nsqAddr, config)
	if err != nil {
		log.Fatal("Could create producer", err)
	}

	msg := Message{
		Title:     "Hello",
		Content:   "World",
		Timestamp: time.Now().String(),
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	err = producer.Publish(topic, payload)
	if err != nil {
		log.Fatal("Could not connect", err)
	}

	producer.Stop()
}

type messageHandler struct{}

func sub() {
	config := nsq.NewConfig()
	config.MaxAttempts = 10
	config.MaxInFlight = 5
	config.MaxRequeueDelay = time.Second * 900
	config.DefaultRequeueDelay = time.Second * 0

	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(&messageHandler{})

	err = consumer.ConnectToNSQLookupd(nsqlookupdAddr)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()
}

// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
// Returning non-nil error will automatically send a REQ command to NSQ to re-queue the message.
func (mh *messageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Message with an empty body is simply ignored/discarded.
		return nil
	}

	return processMessage(m.Body)
}

func processMessage(body []byte) error {
	var message Message

	if err := json.Unmarshal(body, &message); err != nil {
		log.Println("Error when Unmarshaling the message body, Err : ", err)
		return err
	}

	log.Println("Message")
	log.Println("--------------------")
	log.Println("Title : ", message.Title)
	log.Println("Content : ", message.Content)
	log.Println("Timestamp : ", message.Timestamp)
	log.Println("--------------------")
	log.Println("")

	return nil
}
