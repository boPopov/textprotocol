package security

import (
	"sync"
	"time"
)

type RateLimiter interface {
	CreateRateLimiter()
	Allocate() bool
	Release()
}

type CommandRateLimiter interface {
	Setup()
	Allow() bool
}

type RateLimit struct {
	Lock             sync.Mutex
	Channels         chan interface{}
	CommandRateLimit *CommandRateLimit
	RateLimiter
}

type CommandRateLimit struct {
	AwailableTokens  int
	MaxToken         int
	RefillDuration   time.Duration
	LastTimeRefilled time.Time
	CommandRateLimiter
}

func (rateLimit *RateLimit) CreateRateLimiter() {
	rateLimit.Channels = make(chan interface{}, 5)
	rateLimit.CommandRateLimit = new(CommandRateLimit)
	rateLimit.CommandRateLimit.Setup()
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

func (rateLimit *RateLimit) Release() {
	<-rateLimit.Channels
}

func (commandRateLimit *CommandRateLimit) Setup() {
	commandRateLimit.MaxToken = 5
	commandRateLimit.AwailableTokens = commandRateLimit.MaxToken
	commandRateLimit.RefillDuration = time.Duration(15 * time.Second)
	commandRateLimit.LastTimeRefilled = time.Now()
}

func (commandRateLimit *CommandRateLimit) Allow() bool {
	currentTime := time.Now()
	timeDifferenceFromLastPing := currentTime.Sub(commandRateLimit.LastTimeRefilled)

	if timeDifferenceFromLastPing >= commandRateLimit.RefillDuration {
		commandRateLimit.AwailableTokens = commandRateLimit.MaxToken
	}

	if commandRateLimit.AwailableTokens > 0 {
		commandRateLimit.AwailableTokens--
		return true
	}

	return false
}
