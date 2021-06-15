package repository

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
)

type SNSRepository interface {
	ListAllTopics() ([]*sns.Topic, error)
	ListAllSubscriptions() ([]*sns.Subscription, error)
}

type snsRepository struct {
	client snsiface.SNSAPI
	cache  cache.Cache
}

func NewSNSClient(session *session.Session, c cache.Cache) *snsRepository {
	return &snsRepository{
		sns.New(session),
		c,
	}
}

func (r *snsRepository) ListAllTopics() ([]*sns.Topic, error) {
	if v := r.cache.Get("snsListAllTopics"); v != nil {
		return v.([]*sns.Topic), nil
	}

	var topics []*sns.Topic
	input := &sns.ListTopicsInput{}
	err := r.client.ListTopicsPages(input, func(res *sns.ListTopicsOutput, lastPage bool) bool {
		topics = append(topics, res.Topics...)
		return !lastPage
	})
	if err != nil {
		return nil, err
	}

	r.cache.Put("snsListAllTopics", topics)
	return topics, nil
}

func (r *snsRepository) ListAllSubscriptions() ([]*sns.Subscription, error) {
	if v := r.cache.Get("snsListAllSubscriptions"); v != nil {
		return v.([]*sns.Subscription), nil
	}

	var subscriptions []*sns.Subscription
	input := &sns.ListSubscriptionsInput{}
	err := r.client.ListSubscriptionsPages(input, func(res *sns.ListSubscriptionsOutput, lastPage bool) bool {
		subscriptions = append(subscriptions, res.Subscriptions...)
		return !lastPage
	})
	if err != nil {
		return nil, err
	}

	r.cache.Put("snsListAllSubscriptions", subscriptions)
	return subscriptions, nil
}
