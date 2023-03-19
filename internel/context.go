package internel

import "github.com/fzft/my-actor/pkg"

// Suber is the interface that wraps the basic Subscribe methods.
// subscribe to a or more parent actor
type Suber interface {
	Sub(<-chan any)
}

// Puber is the interface that wraps the basic Publish methods.
// publish to a or more child actor
type Puber interface {
	Pub() <-chan any
}

type Storer[K comparable, V any] interface {
	Put(key K, value V)
	Get(key K) (V, bool)
}

type Context[K comparable, V any] struct {
	sub   Suber
	pub   Puber
	store Storer[K, V]

	pid Pid
}

func NewContext[K comparable, V any](pid Pid) *Context[K, V] {
	return &Context[K, V]{pid: pid, store: pkg.NewKeyValueStore[K, V]()}
}

// metrics is the metrics of the actor
