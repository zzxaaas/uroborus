package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
)

type Client struct {
	Producer sarama.SyncProducer
	msgChan  chan *sarama.ProducerMessage
	admin    sarama.ClusterAdmin
}

type Message struct {
	Data  string
	Topic string
}

func (c *Client) SendLog(msg *sarama.ProducerMessage) (err error) {
	select {
	case c.msgChan <- msg:
	default:
		err = fmt.Errorf("msgChan is full")
	}
	return
}

func NewKafkaClient() *Client {
	config := sarama.NewConfig()
	cli := &Client{}
	config.Producer.RequiredAcks = sarama.WaitForAll
	//config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	//config.Producer.Return.Errors = false //设置成false，防止失败过多，没有消费导致阻塞，而且不是必现，极难排查
	//config.Producer.Return.Successes = false

	var err error
	fmt.Println(viper.GetStringSlice("kafka.addrs"))
	cli.admin, err = sarama.NewClusterAdmin(viper.GetStringSlice("kafka.addrs"), config)
	if err != nil {
		panic(err.Error())
	}
	cli.Producer, err = sarama.NewSyncProducer(viper.GetStringSlice("kafka.addrs"), config)
	if err != nil {
		panic(err.Error())
	}
	cli.msgChan = make(chan *sarama.ProducerMessage, viper.GetInt("kafka.chanSize"))
	go cli.sendMsg()
	return cli

}

func (c *Client) CreateTopic(topic string) error {
	return c.admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
}

func PackMsg(topic, value string, key int32) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Key = sarama.StringEncoder(strconv.Itoa(int(key)))
	msg.Value = sarama.StringEncoder(value)
	return msg
}

func (c *Client) sendMsg() {
	for {
		select {
		case msg := <-c.msgChan:
			pid, offset, err := c.Producer.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg failed,err:", err.Error())
				return
			}
			logrus.Infof("send msg successd,pid:%v,offset:%v", pid, offset)
		}
	}
}
