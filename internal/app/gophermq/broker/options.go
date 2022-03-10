package broker

type Options struct {
	MaxPending int
}

type Option func(*Options)

// MaxPending is the maximum valid number of inflight messages.
func MaxPending(m int) Option {
	return func(o *Options) {
		o.MaxPending = m
	}
}
