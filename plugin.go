package sakura

import "context"

type Plugin interface{}

type PluginInitializer interface {
	Initialize(ctx context.Context, sakura *Sakura) error
}

type PluginBeforePublish interface {
	BeforePublish(ctx context.Context, sakura *Sakura, topic string, data []byte) error
}

type PluginAfterPublish interface {
	AfterPublish(ctx context.Context, sakura *Sakura, topic string, data []byte)
}

type PluginBeforeSubscribe interface {
	BeforeSubscribe(ctx context.Context, sakura *Sakura, user, topic string) error
}

type PluginAfterSubscribe interface {
	AfterSubscribe(ctx context.Context, sakura *Sakura, user, topic string)
}

type PluginBeforeUnsubscribe interface {
	BeforeUnsubscribe(ctx context.Context, sakura *Sakura, user, topic string) error
}

type PluginAfterUnsubscribe interface {
	AfterUnsubscribe(ctx context.Context, sakura *Sakura, user, topic string)
}

type PluginBeforeUserDrop interface {
	BeforeUserDrop(ctx context.Context, sakura *Sakura, user string) error
}

type PluginAfterUserDrop interface {
	AfterUserDrop(ctx context.Context, sakura *Sakura, user string)
}

type PluginBeforeTopicDrop interface {
	BeforeTopicDrop(ctx context.Context, sakura *Sakura, topic string) error
}

type PluginAfterTopicDrop interface {
	AfterTopicDrop(ctx context.Context, sakura *Sakura, topic string)
}
