package sys

import (
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	BaseDelay = 2 * time.Second
	MaxDelay  = time.Minute
	Factor    = 2.0
	Jitter    = 0.5
)

type Backoff struct {
	baseDelay    time.Duration // initial delay
	currentDelay time.Duration // current delay
	maxDelay     time.Duration // max backoff delay
	//attempt      int           // retry attempt counter
	factor       float64 // backoff multiplier
	jitterFactor float64 // jitter percent (0.0 - 1.0)
	halfJitter   float64 // half of jitter percent
}

// New creates a Backoff object with sane defaults.
func NewBackoff(baseDelay, maxDelay time.Duration, factor, jitter float64) *Backoff {
	return &Backoff{
		baseDelay:    baseDelay,
		currentDelay: baseDelay,
		maxDelay:     maxDelay,
		factor:       factor,
		jitterFactor: jitter,
		halfJitter:   jitter / 2,
	}
}

func NewDefaultBackoff() *Backoff {
	return NewBackoff(BaseDelay, MaxDelay, Factor, Jitter)
}

// Sleep sleeps for the computed backoff time and increments the attempt count.
func (b *Backoff) Sleep() {
	backoff := b.next()
	log.Debug().Msgf("Backoff: sleeping for %v", backoff)
	time.Sleep(backoff)
	//b.attempt++
}

// Reset sets the attempt counter back to zero, e.g., after a successful call.
func (b *Backoff) Reset() {
	//b.attempt = 0
	b.currentDelay = b.baseDelay
}

// next computes the backoff delay with jitter.
func (b *Backoff) next() time.Duration {
	// // Exponential backoff: baseDelay * (factor ^ attempt)
	// backoff := float64(b.baseDelay) * math.Pow(b.factor, float64(b.attempt))
	// b.currentDelay *= 2

	// if backoff > float64(b.maxDelay) {
	// 	backoff = float64(b.maxDelay)
	// }

	// // Apply jitter
	// jitter := rand.Float64() * b.jitterFactor * backoff
	// backoffWithJitter := backoff - (b.jitterFactor/2)*backoff + jitter

	// return time.Duration(backoffWithJitter)

	// Apply jitter
	jitter := rand.Float64() * b.jitterFactor * float64(b.currentDelay)
	delay := float64(b.currentDelay) - b.halfJitter*float64(b.currentDelay) + jitter

	// Increase currentDelay for next call
	b.currentDelay = MinDuration(time.Duration(float64(b.currentDelay)*b.factor), b.maxDelay)

	return time.Duration(delay)
}
