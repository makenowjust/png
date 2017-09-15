package png

import (
	"context"
	"net"
	"net/url"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type AMQPPinger struct {
	url *url.URL
}

func (p *AMQPPinger) Ping(ctx context.Context) error {
	done := make(chan error)

	go func() {
		conn, err := amqp.DialConfig(p.url.String(), amqp.Config{
			Dial: func(network, addr string) (net.Conn, error) {
				dialer := &net.Dialer{}
				if conn, err := dialer.DialContext(ctx, network, addr); err != nil {
					return nil, err
				} else {
					if t, ok := ctx.Deadline(); ok {
						conn.SetDeadline(t)
						conn.SetReadDeadline(t)
						conn.SetWriteDeadline(t)
					}
					return conn, err
				}
			},
		})
		if err != nil {
			done <- errors.Wrap(err, "failed in connecting to AMQP server")
			return
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			done <- errors.Wrap(err, "failed in creating channel")
			return
		}
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			true,  // autoDelete
			true,  // exclusive
			false, // noWait
			nil,   // args
		)
		if err != nil {
			done <- errors.Wrap(err, "failed in creating queue")
			return
		}
		defer ch.QueueDelete(q.Name, false, false, false)

		if err := ch.Publish(
			"",     // exchange,
			q.Name, // key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				Body: []byte("png"),
			},
		); err != nil {
			done <- errors.Wrap(err, "failed in publush message")
			return
		}

		wait, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // autoAck
			true,   // exclusive
			false,  // noLocal
			false,  // noWait
			nil,    // args
		)
		if err != nil {
			done <- errors.Wrap(err, "failed in consume message")
			return
		}

		select {
		case <-ctx.Done():
			done <- errors.Wrap(err, "failed in consume message")
		case d := <-wait:
			if string(d.Body) == "png" {
				done <- nil
			} else {
				done <- errors.Errorf("invalid AMQP response: %#v", d.Body)
			}
		}

		return
	}()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "failed in AMQP connection")
	case err := <-done:
		return err
	}
}
