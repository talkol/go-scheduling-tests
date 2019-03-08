package ratelimiter

type RateLimiter struct {
	slots chan struct{}
}

func NewRateLimiter(maxConcurrent int) *RateLimiter {
	return &RateLimiter{
		slots: make(chan struct{}, maxConcurrent),
	}
}

func (r *RateLimiter) RequestSlot() {
	r.slots <- struct{}{}
}

func (r *RateLimiter) ReleaseSlot() {
	<- r.slots
}