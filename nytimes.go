package nytimes

import "github.com/bouwerp/log"

type Client struct {
	Key            string
	Secret         string
	LoggingEnabled bool
	Logger         log.Logger
}
