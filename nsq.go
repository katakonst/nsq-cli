package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

type MessageHandler struct {
	countChannel chan int
}

func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	h.countChannel <- 1
	Logger.Infof("message body: %s", m.Body)
	return nil
}

func (h *MessageHandler) CountMessages(limit int, stopChan chan int) {
	for i := 0; i < limit; i++ {
		<-h.countChannel
	}
	stopChan <- 1
}

func Publish(conf Config) error {
	url := conf.Host + "/pub?topic=" + conf.Topic
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(conf.Payload)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("publishing to nsq: %s", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	Logger.Infof("response Body:", string(body))

	return nil
}

func Consume(conf Config) error {

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(conf.Topic, "metrics", config)
	if err != nil {
		return fmt.Errorf("creating consummer: %s", err)
	}
	consumer.ChangeMaxInFlight(200)
	msgHandler := MessageHandler{countChannel: make(chan int)}

	terminate := make(chan int)
	go msgHandler.CountMessages(conf.MaxLimit, terminate)

	consumer.AddConcurrentHandlers(
		&msgHandler,
		3,
	)

	nsqlds := []string{conf.Lookupd}
	if err := consumer.ConnectToNSQLookupds(nsqlds); err != nil {
		return fmt.Errorf("connecting to lookup: %s", err)
	}
	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT)

	for {
		select {
		case <-consumer.StopChan:
			return nil
		case <-shutdown:
			consumer.Stop()
			os.Exit(1)
		case <-terminate:
			consumer.Stop()
		}
	}
}

func Empty(conf Config) error {
	url := conf.Host + "/topic/empty?topic=" + conf.Topic
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("empty topic: %s", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	Logger.Infof("empty response body:", string(body))

	return nil
}

func Delete(conf Config) error {
	url := conf.Host + "/topic/delete?topic=" + conf.Topic
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("delete topic: %s", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	Logger.Infof("delete response body:", string(body))

	return nil
}
