package message

import (
	"context"
	"encoding/json"

	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

func Subscribe(ctx context.Context, handler interface{}, log *logrus.Entry) error {
	subscribers := make(map[string]mq.Subscriber)

	defer func() {
		for k, s := range subscribers {
			if err := s.Unsubscribe(); err != nil {
				log.Errorf("failed to unsubscribe for topic:%s, err:%v", k, err)
			}
		}
	}()

	// register following
	s, err := registerFollowingHandler(handler)
	if err != nil {
		return err
	}
	if s != nil {
		subscribers[s.Topic()] = s
	}

	// register like
	s, err = registerLikeHandler(handler)
	if err != nil {
		return err
	}
	if s != nil {
		subscribers[s.Topic()] = s
	}

	// register end
	if len(subscribers) == 0 {
		return nil
	}

	<-ctx.Done()

	return nil
}

func registerFollowingHandler(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.FollowingHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.Following, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgFollowing{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		f := domain.Following{}
		if f.Owner, err = domain.NewAccount(body.Owner); err != nil {
			return
		}

		if f.Account, err = domain.NewAccount(body.Following); err != nil {
			return
		}

		switch body.Action {
		case actionAdd:
			return h.HandleEventAddFollowing(f)

		case actionRemove:
			return h.HandleEventRemoveFollowing(f)
		}

		return nil
	})
}

func registerLikeHandler(handler interface{}) (mq.Subscriber, error) {
	h, ok := handler.(message.LikeHandler)
	if !ok {
		return nil, nil
	}

	return kafka.Subscribe(topics.Like, func(e mq.Event) (err error) {
		msg := e.Message()
		if msg == nil {
			return
		}

		body := msgLike{}
		if err = json.Unmarshal(msg.Body, &body); err != nil {
			return
		}

		like := domain.Like{}
		if like.ResourceOwner, err = domain.NewAccount(body.Owner); err != nil {
			return
		}

		if like.ResourceType, err = domain.NewResourceType(body.Type); err != nil {
			return
		}

		like.ResourceId = body.Id

		switch body.Action {
		case actionAdd:
			return h.HandleEventAddLike(like)

		case actionRemove:
			return h.HandleEventRemoveLike(like)
		}

		return nil
	})
}
