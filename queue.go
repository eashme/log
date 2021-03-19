package log

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"time"
)

type KafkaCli struct {
	brokers  []string // broker 地址列表
	topic    string   // 上报的topic
	producer sarama.AsyncProducer
}

func (cli *KafkaCli) Levels() []logrus.Level {
	return []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel}
}

func (cli *KafkaCli) Fire(entry *logrus.Entry) error {
	cli.Send(entry.Message)
	return nil
}

func NewKafkaCli(topic string, broker ...string) (*KafkaCli, error) {
	cli := &KafkaCli{
		brokers: broker,
		topic:   topic,
	}
	conf := sarama.NewConfig()
	conf.Version = sarama.V0_11_0_0
	conf.Producer.Interceptors = []sarama.ProducerInterceptor{cli}
	conf.ChannelBufferSize = 10
	conf.Producer.Retry.Max = 3
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(cli.brokers, conf)
	if err != nil {
		logrus.Errorf("failed create producer ", err)
		return nil, err
	}
	cli.producer = producer

	// 错误日志
	go func() {
		for err := range producer.Errors() {
			logrus.Errorf("failed send msg(%s) err(%v):", Json(err.Msg), err.Err)
		}
	}()

	return cli, nil
}

func (*KafkaCli) OnSend(msg *sarama.ProducerMessage) {
	logrus.Info("send msg %s", msg.Value)
}

func (cli *KafkaCli) Send(msg string) {
	go func() {
		cli.producer.Input() <- &sarama.ProducerMessage{
			Topic: cli.topic, Key: nil,
			Value:     sarama.StringEncoder(msg),
			Timestamp: time.Now(),
		}
	}()
}

func (cli *KafkaCli) Close() {
	cli.producer.AsyncClose()
}

func Json(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
