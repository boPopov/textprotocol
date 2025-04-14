package security

import (
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
	rateLimit.Channels = make(chan interface{}, maxSessions)
	rateLimit.CommandRateLimit = new(CommandRateLimit)
	rateLimit.CommandRateLimit.Setup(maxInputPerInterval, refillDurationInterval)
}

/**
 * Function Allocate(), is responsible for Allocating a space for a connection.
 * Each client from a given IP has maximum of 5 available connection.
 * Once this function is called, we lock the function for the other.
 */ 
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
	commandRateLimit.RefillDuration = time.Duration(refillDurationInterval * time.Second)
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
