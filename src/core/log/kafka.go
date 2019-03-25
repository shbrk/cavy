package log

import "github.com/Shopify/sarama" //使用的v1.9，以后的版本引入了cgo

// kafka case for zap logger
type KafkaLogger struct {
	Producer sarama.SyncProducer
	Topic    string
}

func (kl *KafkaLogger) Write(p []byte) (n int, err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = kl.Topic
	msg.Value = sarama.ByteEncoder(p)
	_, _, err = kl.Producer.SendMessage(msg)
	if err != nil {
		return
	}
	return
}

func newKafkaLogger(addr []string, topic string) *KafkaLogger {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		panic("kafka producer create failed:" + err.Error())
	}
	return &KafkaLogger{p, topic}
}
