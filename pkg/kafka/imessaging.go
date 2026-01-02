package kafka

import (
	"context"
)

//go:generate mockgen -source=imessaging.go -destination=messaging_mock.go -package=kafka

// CallBack ...
type CallBack func(context.Context, string, []byte) error

// IPublisher ...
type Publisher interface {
	Write(v interface{})
	WriteByTopic(interface{}, string) error
	WriteByTopicAndKey(v interface{}, key, topicName string) error

	Close()
}

// ISubscriber ...
type Subscriber interface {
	Read(callback CallBack, errRestart chan error)

	Close()
}
