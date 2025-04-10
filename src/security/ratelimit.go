package security

import "sync"

type RateLimiter interface {
	CreateRateLimiter()
	Allocate()
	Release()
}

type RateLimit struct {
	Lock     sync.Mutex
	Channels chan interface{}
}

func (rateLimit *RateLimit) CreateRateLimiter() {
	rateLimit.Channels = make(chan interface{}, 5)
}

func (rateLimit *RateLimit) Allocate() bool {
	rateLimit.Lock.Lock()

	select {
	case rateLimit.Channels <- 1:
		rateLimit.Lock.Unlock()
		return true
	default:
		rateLimit.Lock.Unlock()
		return false
	}
}

func (RateLimit *RateLimit) Release() {
	<-RateLimit.Channels
}
