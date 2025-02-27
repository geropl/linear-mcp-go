package linear

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter manages API request rate limiting
type RateLimiter struct {
	requestsPerHour int
	minDelayMs      int64
	queue           chan func() error
	processing      bool
	lastRequestTime int64
	requestTimes    []int64
	requestDurations []int64
	mu              sync.Mutex
}

// RateLimiterMetrics contains metrics about the rate limiter
type RateLimiterMetrics struct {
	TotalRequests      int    `json:"totalRequests"`
	RequestsInLastHour int    `json:"requestsInLastHour"`
	AverageRequestTime int64  `json:"averageRequestTime"`
	QueueLength        int    `json:"queueLength"`
	LastRequestTime    int64  `json:"lastRequestTime"`
}

// NewRateLimiter creates a new rate limiter with the specified requests per hour limit
func NewRateLimiter(requestsPerHour int) *RateLimiter {
	minDelayMs := int64(3600000 / requestsPerHour)
	rl := &RateLimiter{
		requestsPerHour:  requestsPerHour,
		minDelayMs:       minDelayMs,
		queue:            make(chan func() error, 100), // Buffer size of 100
		processing:       false,
		lastRequestTime:  0,
		requestTimes:     []int64{},
		requestDurations: []int64{},
	}

	// Start the queue processor
	go rl.processQueue()

	return rl
}

// Enqueue adds a function to the rate limiter queue
func (rl *RateLimiter) Enqueue(fn func() error, operation string) error {
	startTime := time.Now().UnixMilli()
	queuePosition := len(rl.queue)

	fmt.Printf("[Linear API] Enqueueing request%s (Queue position: %d)\n", 
		formatOperation(operation), queuePosition)

	// Create a channel to receive the result
	resultCh := make(chan error, 1)

	// Wrap the function to capture its result
	wrappedFn := func() error {
		fmt.Printf("[Linear API] Starting request%s\n", formatOperation(operation))
		result := fn()
		endTime := time.Now().UnixMilli()
		duration := endTime - startTime

		fmt.Printf("[Linear API] Completed request%s (Duration: %dms)\n", 
			formatOperation(operation), duration)
		
		rl.trackRequest(startTime, endTime, operation)
		resultCh <- result
		return result
	}

	// Add to queue
	rl.queue <- wrappedFn

	// Wait for the result
	err := <-resultCh
	return err
}

// Batch processes a batch of items with the rate limiter
func (rl *RateLimiter) Batch(items []interface{}, batchSize int, fn func(item interface{}) error, operation string) []error {
	results := make([]error, len(items))
	var wg sync.WaitGroup

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		wg.Add(len(batch))

		for j, item := range batch {
			index := i + j
			go func(idx int, itm interface{}) {
				defer wg.Done()
				err := rl.Enqueue(func() error {
					return fn(itm)
				}, operation)
				results[idx] = err
			}(index, item)
		}

		wg.Wait()
	}

	return results
}

// processQueue processes the queue of functions
func (rl *RateLimiter) processQueue() {
	for {
		// If there are no items in the queue, wait for one
		fn := <-rl.queue

		// Process the item
		rl.mu.Lock()
		now := time.Now().UnixMilli()
		timeSinceLastRequest := now - rl.lastRequestTime

		// Check if we need to wait to respect rate limits
		requestsInLastHour := 0
		oneHourAgo := now - 3600000
		for _, t := range rl.requestTimes {
			if t > oneHourAgo {
				requestsInLastHour++
			}
		}

		if requestsInLastHour >= int(float64(rl.requestsPerHour)*0.9) && timeSinceLastRequest < rl.minDelayMs {
			waitTime := rl.minDelayMs - timeSinceLastRequest
			rl.mu.Unlock()
			time.Sleep(time.Duration(waitTime) * time.Millisecond)
		} else {
			rl.mu.Unlock()
		}

		// Execute the function
		rl.mu.Lock()
		rl.lastRequestTime = time.Now().UnixMilli()
		rl.mu.Unlock()

		_ = fn() // Execute the function
	}
}

// trackRequest tracks a request for metrics
func (rl *RateLimiter) trackRequest(startTime, endTime int64, operation string) {
	duration := endTime - startTime

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.requestTimes = append(rl.requestTimes, startTime)
	rl.requestDurations = append(rl.requestDurations, duration)

	// Keep only the last hour of requests
	oneHourAgo := time.Now().UnixMilli() - 3600000
	var newRequestTimes []int64
	var newRequestDurations []int64

	for i, t := range rl.requestTimes {
		if t > oneHourAgo {
			newRequestTimes = append(newRequestTimes, t)
			newRequestDurations = append(newRequestDurations, rl.requestDurations[i])
		}
	}

	rl.requestTimes = newRequestTimes
	rl.requestDurations = newRequestDurations
}

// GetMetrics returns metrics about the rate limiter
func (rl *RateLimiter) GetMetrics() RateLimiterMetrics {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now().UnixMilli()
	oneHourAgo := now - 3600000
	recentRequests := 0

	for _, t := range rl.requestTimes {
		if t > oneHourAgo {
			recentRequests++
		}
	}

	var avgRequestTime int64 = 0
	if len(rl.requestDurations) > 0 {
		var sum int64 = 0
		for _, d := range rl.requestDurations {
			sum += d
		}
		avgRequestTime = sum / int64(len(rl.requestDurations))
	}

	return RateLimiterMetrics{
		TotalRequests:      len(rl.requestTimes),
		RequestsInLastHour: recentRequests,
		AverageRequestTime: avgRequestTime,
		QueueLength:        len(rl.queue),
		LastRequestTime:    rl.lastRequestTime,
	}
}

// Helper function to format operation name for logging
func formatOperation(operation string) string {
	if operation != "" {
		return " for " + operation
	}
	return ""
}
