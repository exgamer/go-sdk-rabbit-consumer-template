package publishers

import (
	"github.com/exgamer/gosdk-rabbit-core/pkg/config"
)

func GetPublishers() []config.PublisherDefinition {
	return []config.PublisherDefinition{
		{
			Name:   "TestPublisher",
			Config: config.NewPublisherTopicDurableConfig("test"),
		},
	}
}
