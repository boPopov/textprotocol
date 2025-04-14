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
 * Once this function is called from an IP, a temporary lock is called to block multiple connection at the same time.
 * If there are available spaces 1 is sent to the Channel and the lock is unlocked.
 * If there are no available spaces the lock is unlocked and a false flag is returned.
 */
func (rateLimit *RateLimit) Allocate() bool {
	rateLimit.Lock.Lock()

	select {
	case rateLimit.Channels <- 1: //There is still a free channel for a connection.
		rateLimit.Lock.Unlock()
		return true
	default:
		rateLimit.Lock.Unlock() //There are no more free channels for a connection.
		return false
	}
}

/**
 * Function Release, is called to release a channel an enable to user from a dedicated IP to make a new connection.
 */
func (rateLimit *RateLimit) Release() {
	<-rateLimit.Channels
}

func (commandRateLimit *CommandRateLimit) Setup(maxInputPerInterval int, refillDurationInterval int) {
	commandRateLimit.MaxToken = maxInputPerInterval
	commandRateLimit.AwailableTokens = commandRateLimit.MaxToken
	commandRateLimit.RefillDuration = time.Duration(int64(refillDurationInterval) * int64(time.Second))
	commandRateLimit.LastTimeRefilled = time.Now()
}

/**
 * Function Allow, handles how many commands are entered in a time interval.
 * With this function we don't allow the user to spam the application with commands more than the value of AwailableTokens variable.
 * After a given interval the AwailableTokens variable is reset to the Maximum value.
 */
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
