package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rank1zen/yujin/internal"
)

type Summoner struct {
	ch *amqp.Channel
}

func NewSummoner(ch *amqp.Channel) *Summoner {
	return &Summoner{
		ch: ch,
	}
}

func (s *Summoner) PassOn(ctx context.Context,) error {
	err := s.ch.PublishWithContext(
		ctx,
		"tasks",
		"HELP",
		false,
		false,
		amqp.Publishing{},
	)

	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "ch.Publish")
	}
	return nil
}
