package mq

import (
	"fmt"
	"yimcom/comm"
	"yimcom/seelog"

	"github.com/streadway/amqp"
)

func TestJob(d amqp.Delivery, exta interface{}) error {

	seelog.Debug(fmt.Sprintf("got %dB delivery: [%v] %q %s", len(d.Body), d.DeliveryTag, d.Body, d.RoutingKey))

	d.Ack(false)

	return nil
}

type MqConsumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	tag        string
	closeChan  chan *amqp.Error
	deliveries <-chan amqp.Delivery
	job        JobHandle
	expires    int64
	ttl        int64
	exta       interface{}
	closeFlag  bool
	amqpURI    string
	exchange   string
	key        string
	queueName  string
}

func (c *MqConsumer) Stop() {
	if !c.closeFlag {
		c.channel.Close()
		c.conn.Close()
	}
}

func (c *MqConsumer) Run() {

	defer comm.PinacRecover()

	for {
		select {
		case d, ok := <-c.deliveries:

			if !ok {
				seelog.Error("MqConsumer <-c.deliveries no ok ##########")
				c.closeFlag = true
				return
			}

			err := c.job(d, c.exta)
			if err != nil {
				seelog.Error("MqConsumer run deliveries job do error:", err)
				continue
			}

		case ce := <-c.closeChan:
			seelog.Info("MqConsumer notify conn close error:", ce)
			c.closeFlag = true
			return
		}

	}
}

func (c *MqConsumer) Check() {

	if c.closeFlag == true {

		var err error

		c.conn, err = amqp.Dial(c.amqpURI)

		if err != nil {
			seelog.Error("MqConsumer Check dial error:", err)
			return
		}

		seelog.Debug("amqp dial :", c.amqpURI)

		c.closeChan = make(chan *amqp.Error)

		c.conn.NotifyClose(c.closeChan)

		c.channel, err = c.conn.Channel()
		if err != nil {
			seelog.Error("MqConsumer Check conn Channel error:", err)
			return
		}

		seelog.Debug("Check amqp channel :", c.amqpURI)

		err = c.channel.ExchangeDeclare(
			c.exchange, // name of the exchange
			"fanout",   // type
			true,       // durable
			false,      // delete when complete
			false,      // internal
			false,      // noWait
			nil,        // arguments
		)

		seelog.Debug("Check amqp ExchangeDeclare ")

		if err != nil {
			seelog.Error("MqConsumer Check ExchangeDeclare error:", err)
			return
		}

		args := make(amqp.Table)
		args["x-message-ttl"] = int64(300000)
		args["x-expires"] = int64(600000)

		queue, err := c.channel.QueueDeclare(
			c.queueName, // name of the queue
			true,        // durable
			false,       // delete when unused
			false,       // exclusive
			false,       // noWait
			args,        // arguments
		)

		seelog.Debug("amqp QueueDeclare ")

		if err != nil {
			seelog.Error("MqConsumer check QueueDeclare error:", err)
			return
		}

		err = c.channel.QueueBind(
			queue.Name, // name of the queue
			c.key,      // bindingKey
			c.exchange, // sourceExchange
			false,      // noWait
			nil,        // arguments
		)
		if err != nil {
			seelog.Error("MqConsumer check QueueBind error:", err)
			return
		}

		deliveries, err := c.channel.Consume(
			queue.Name, // name
			c.tag,      // consumerTag,
			false,      // noAck
			false,      // exclusive
			false,      // noLocal
			false,      // noWait
			nil,        // arguments
		)
		if err != nil {
			seelog.Error("MqConsumer check Consume error:", err)
			return
		}

		c.deliveries = deliveries

		go c.Run()

		c.closeFlag = false

	}
}

type JobHandle func(amqp.Delivery, interface{}) error

func NewMqConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string, job JobHandle, exta interface{}) (*MqConsumer, error) {

	c := &MqConsumer{
		conn:       nil,
		channel:    nil,
		tag:        ctag,
		closeChan:  make(chan *amqp.Error),
		deliveries: make(<-chan amqp.Delivery),
		job:        job,
		exta:       exta,
		amqpURI:    amqpURI,
		exchange:   exchange,
		queueName:  queueName,
	}

	var err error

	c.conn, err = amqp.Dial(amqpURI)

	if err != nil {
		seelog.Error("NewMqConsumer dial error:", err)
		return nil, err
	}

	seelog.Debug("amqp dial :", amqpURI)

	c.conn.NotifyClose(c.closeChan)

	c.channel, err = c.conn.Channel()
	if err != nil {
		seelog.Error("NewMqConsumer conn Channel error:", err)
		return nil, err
	}

	//c.channel.NotifyClose(c.closeChan)

	seelog.Debug("amqp channel :", amqpURI)

	err = c.channel.ExchangeDeclare(
		exchange, // name of the exchange
		"fanout", // type
		true,     // durable
		false,    // delete when complete
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)

	seelog.Debug("amqp ExchangeDeclare ")

	if err != nil {
		seelog.Error("NewMqConsumer ExchangeDeclare error:", err)
		return nil, err
	}

	args := make(amqp.Table)
	args["x-message-ttl"] = int64(300000)
	args["x-expires"] = int64(600000)

	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		args,      // arguments
	)

	seelog.Debug("amqp QueueDeclare ")

	if err != nil {
		seelog.Error("NewMqConsumer QueueDeclare error:", err)
		return nil, err
	}

	err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		seelog.Error("NewMqConsumer QueueBind error:", err)
		return nil, err
	}

	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		seelog.Error("NewMqConsumer Consume error:", err)
		return nil, err
	}

	c.deliveries = deliveries
	//go c.Run()

	return c, nil
}

/////////////////////////////
type MqProducet struct {
}

func (p *MqProducet) PublishFanoutNoConfirmation(uri string, exName string, body []byte) error {

	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(uri)
	if err != nil {
		seelog.Error("Publish amqp Dial error:", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		seelog.Error("Publish conn Channel error:", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exName,   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		seelog.Error("Publish ExchangeDeclare error:", err)
		return err
	}

	err = ch.Publish(
		exName, // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	if err != nil {
		seelog.Error("Publish mq error:", err)
		return err
	}

	seelog.Debug("Publish msg:", string(body))

	return nil
}

func (p *MqProducet) PublishNoConfirmation(uri string, exName string, topic string, body []byte) error {

	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(uri)
	if err != nil {
		seelog.Error("Publish amqp Dial error:", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		seelog.Error("Publish conn Channel error:", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exName,  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		seelog.Error("Publish ExchangeDeclare error:", err)
		return err
	}

	err = ch.Publish(
		exName, // exchange
		topic,  // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	if err != nil {
		seelog.Error("Publish mq error:", err)
		return err
	}

	seelog.Debug("Publish msg:", string(body))

	return nil
}
