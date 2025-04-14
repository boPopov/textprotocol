package security

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter interface {
	CreateRateLimiter(maxSessions int)
	Allocate() bool
	Release()
}

type CommandRateLimiter interface {
	Setup(maxInputPerInterval int, refillDurationInterval int)
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

func (rateLimit *RateLimit) CreateRateLimiter(maxSessions int, maxInputPerInterval int, refillDurationInterval int) {
	fmt.Println("Max Sessions is: ", maxSessions)
	fmt.Println("Max Input Per Interval is:", int64(maxInputPerInterval))
	fmt.Println("Refill Duration Interval is:", int64(refillDurationInterval))
	rateLimit.Channels = make(chan interface{}, maxSessions)
	rateLimit.CommandRateLimit = new(CommandRateLimit)
	rateLimit.CommandRateLimit.Setup(maxInputPerInterval, refillDurationInterval)
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

func (commandRateLimit *CommandRateLimit) Setup(maxInputPerInterval int, refillDurationInterval int) {
	commandRateLimit.MaxToken = maxInputPerInterval
	commandRateLimit.AwailableTokens = commandRateLimit.MaxToken
	commandRateLimit.RefillDuration = time.Duration(int64(refillDurationInterval) * int64(time.Second))
	commandRateLimit.LastTimeRefilled = time.Now()
}

func (commandRateLimit *CommandRateLimit) Allow() bool {
	currentTime := time.Now()
	timeDifferenceFromLastPing := currentTime.Sub(commandRateLimit.LastTimeRefilled)

	if timeDifferenceFromLastPing >= commandRateLimit.RefillDuration {
		commandRateLimit.AwailableTokens = commandRateLimit.MaxToken
		commandRateLimit.LastTimeRefilled = time.Now()
	}

	if commandRateLimit.AwailableTokens > 0 {
		commandRateLimit.AwailableTokens--
		return true
	}

	return false
}
