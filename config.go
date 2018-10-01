package main

import "flag"

type Config struct {
	Topic     string
	Host      string
	Payload   string
	Consume   bool
	ShowStats bool
	Lookupd   string
	MaxLimit  int
	Empty     bool
	Delete    bool
}

func InitConfigs() Config {
	topic := flag.String("topic", "topic", "nsq topic")
	host := flag.String("host", "http://localhost:4151", "nsq host")
	payload := flag.String("payload", "test", "payload to post")
	consume := flag.Bool("consume", false, "consume")
	delete := flag.Bool("delete", false, "delete")
	showStats := flag.Bool("stats", false, "show stats")
	nsqLookupd := flag.String("lookupds", "localhost:4161", "lookupds comma separated")
	maxlimit := flag.Int("limit", 2, "max limit of messages consumed")
	empty := flag.Bool("empty", false, "empty topic")

	return Config{
		Topic:     *topic,
		Host:      *host,
		Payload:   *payload,
		Consume:   *consume,
		ShowStats: *showStats,
		Lookupd:   *nsqLookupd,
		MaxLimit:  *maxlimit,
		Empty:     *empty,
		Delete:    *delete,
	}
}
