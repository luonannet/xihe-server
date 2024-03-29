package message

import (
	"encoding/json"

	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

var topics Topics

type Topics struct {
	Following string
	Like      string
}

func NewMessageSender() message.Sender {
	return sender{}
}

type sender struct{}

// Following
func (s sender) AddFollowing(msg domain.Following) error {
	return s.sendFollowing(msg, actionAdd)
}

func (s sender) RemoveFollowing(msg domain.Following) error {
	return s.sendFollowing(msg, actionRemove)
}

func (s sender) sendFollowing(msg domain.Following, action string) error {
	v := msgFollowing{
		Action:    action,
		Owner:     msg.Owner.Account(),
		Following: msg.Account.Account(),
	}

	return s.send(topics.Following, &v)
}

// Like
func (s sender) AddLike(msg domain.Like) error {
	return s.sendLike(msg, actionAdd)
}

func (s sender) RemoveLike(msg domain.Like) error {
	return s.sendLike(msg, actionRemove)
}

func (s sender) sendLike(msg domain.Like, action string) error {
	v := msgLike{
		Action: action,
		Owner:  msg.ResourceOwner.Account(),
		Type:   msg.ResourceType.ResourceType(),
		Id:     msg.ResourceId,
	}

	return s.send(topics.Like, &v)
}

func (s sender) send(topic string, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return kafka.Publish(topic, &mq.Message{
		Body: body,
	})
}
