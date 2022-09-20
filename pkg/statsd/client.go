package statsd

import (
	"gopkg.in/alexcesaro/statsd.v2"
)

type Client struct {
	client *statsd.Client
}

func NewClient(addr, prefix string) (*Client, error) {
	cli, err := statsd.New(
		statsd.Address(addr),
		statsd.Prefix(prefix),
	)
	if err != nil {
		return nil, err
	}

	return &Client{client: cli}, nil
}

func (s *Client) Close() {
	s.client.Close()
}

func (s *Client) Increment(bucket string) {
	s.client.Increment(bucket)
}

func (s *Client) Timing(bucket string, value interface{}) {
	s.client.Timing(bucket, value)
}

func (s *Client) Count(bucket string, value interface{}) {
	s.client.Count(bucket, value)
}

func (s *Client) Gauge(bucket string, value interface{}) {
	s.client.Gauge(bucket, value)
}

func (s *Client) Unique(bucket string, value string) {
	s.client.Unique(bucket, value)
}
