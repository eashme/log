package coord_log

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
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	cli.Send(b)
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

	producer, err := sarama.NewAsyncProducer(cli.brokers, conf)
	if err != nil {
		logrus.Errorf("failed create producer ", err)
		return nil, err
	}
	cli.producer = producer
	return cli, nil
}

func (*KafkaCli) OnSend(msg *sarama.ProducerMessage) {
	logrus.Info("send msg %s", msg.Value)
}

func (cli *KafkaCli) Send(msg []byte) {
	go func() {
		cli.producer.Input() <- &sarama.ProducerMessage{
			Topic: cli.topic, Key: nil,
			Value:     sarama.ByteEncoder(msg),
			Timestamp: time.Now(),
		}
	}()
}

func (cli *KafkaCli) Close() {
	cli.producer.AsyncClose()
}
